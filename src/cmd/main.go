package main

import (
	"log"
	"os"
	logger "src/logger"
	api_keycloak "src/api/keycloak"
	app_router "src/api/router"
	"src/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var repositoryWrapper *repositories.RepositoryWrapper

func LoadRepositoryWrapper() {
	transactionRepository := repositories.NewTransactionRepository(db.DB, logger)
	accountRepository := repositories.NewAccountRepository(db.DB, logger)
	clientRepository := repositories.NewClientRepository(db.DB, logger)
	registryAccountOtpRepository := repositories.NewRegistryAccountOtpRepository(db.DB, logger)
	repositoryWrapper = &repositories.RepositoryWrapper{
		ClientRepository:             clientRepository,
		AccountRepository:            accountRepository,
		TransactionRepository:        transactionRepository,
		RegistryAccountOtpRepository: registryAccountOtpRepository,
	}
}
func initializer() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Warning: Could not load .env file: %v. Falling back to system environment variables.", err)
		panic("environment variables could not be loaded!")
	}
	// logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug, // Debug level for detailed test output
	// }))

	logger := logger.GetLogger()

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
	keycloakClient := api_keycloak.BuildKeycloakClientFromEnv()
	router := gin.Default()

	appRouter := app_router.AppRouter{
		KeycloakClient: &keycloakClient,
		RepositoryWrapper: repositoryWrapper,
	}

	appRouter.BuildRoutes(router)
	router.Run() // Listen on :8080 by default
	defer db.Close()
}
