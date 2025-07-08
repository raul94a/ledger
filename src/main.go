package main

import (
	"context"
	"fmt"
	"log"
	"os"
	api_keycloak "src/api/keycloak"
	app_router "src/api/router"
	appRedis "src/db/redis"
	docs "src/docs"
	logger "src/logger"
	"src/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var zlogger *zap.Logger
var db *gorm.DB
var dbSqlite *sqlx.DB
var repositoryWrapper *repositories.RepositoryWrapper

func LoadRepositoryWrapper() {
	transactionRepository := repositories.NewTransactionRepository(dbSqlite.DB, zlogger)
	accountRepository := repositories.NewAccountRepository(db, zlogger)
	clientRepository := repositories.NewClientRepository(db, zlogger)
	registryAccountOtpRepository := repositories.NewRegistryAccountOtpRepository(db, zlogger)
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

	connectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	dbSqlite, _ = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
		panic("error " + err.Error())
	}

	// // Test the connection to the database
	
	if err := dbSqlite.Ping(); err != nil {
		log.Fatal(err)
		panic("error " + err.Error())
	} else {
		log.Println("Successfully Connected")
		LoadRepositoryWrapper()
	}

}

func initSwagger(router *gin.Engine){
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/apidoc/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


}

// @title API Bank Clients
// @version 1.0
// @description Clients management of a bank system.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	initializer()
	redisClient := appRedis.Get()
	appRedis.CreateAllIndexes(context.Background(),redisClient,zlogger)
	

	keycloakClient := api_keycloak.BuildKeycloakClientFromEnv()
	router := gin.Default()
	
	initSwagger(router)
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
	// defer db.Close()
	zlogger.Fatal("Server has been shut down")

}
