package main

import (
	"context"
	"fmt"
	"log"
	"os"
	api_keycloak "src/api/keycloak"
	app_router "src/api/router"
	appRedis "src/db/redis"
	logger "src/logger"
	"src/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)
var zlogger *zap.Logger
var db *sqlx.DB
var repositoryWrapper *repositories.RepositoryWrapper

func LoadRepositoryWrapper() {
	transactionRepository := repositories.NewTransactionRepository(db.DB, zlogger)
	accountRepository := repositories.NewAccountRepository(db.DB, zlogger)
	clientRepository := repositories.NewClientRepository(db.DB, zlogger)
	registryAccountOtpRepository := repositories.NewRegistryAccountOtpRepository(db.DB, zlogger)
	repositoryWrapper = &repositories.RepositoryWrapper{
		ClientRepository:             clientRepository,
		AccountRepository:            accountRepository,
		TransactionRepository:        transactionRepository,
		RegistryAccountOtpRepository: registryAccountOtpRepository,
	}
}
func initializer() {
	zlogger = logger.GetLogger()
	err := godotenv.Load()

	if err != nil {
		zlogger.Sugar().Warn("Warning: Could not load .env file: %v. Falling back to system environment variables. " +  err.Error())
		panic("environment variables could not be loaded!")
	}
	// logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug, // Debug level for detailed test output
	// }))


	connectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
		panic("error " + err.Error())
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		panic("error " + err.Error())
	} else {
		log.Println("Successfully Connected")
		LoadRepositoryWrapper()
	}

}

// @title API Bank Clients
// @version 1.0
// @description Clients management of a bank system.
// @host localhost:8080
// @BasePath /
func main() {
	initializer()
	redisClient := appRedis.Get()
	appRedis.CreateAllIndexes(context.Background(),redisClient,zlogger)
	

	keycloakClient := api_keycloak.BuildKeycloakClientFromEnv()
	router := gin.Default()

	appRouter := app_router.AppRouter{
		KeycloakClient: &keycloakClient,
		RepositoryWrapper: repositoryWrapper,
		RedisClient: redisClient,
		ZapLogger: zlogger,
	}

	
	appRouter.BuildRoutes(router)
	fmt.Println("Ledger is running")
	zlogger.Info("Server has been started")
	router.Run() // Listen on :8080 by default
	defer db.Close()
	zlogger.Fatal("Server has been shut down")

}
