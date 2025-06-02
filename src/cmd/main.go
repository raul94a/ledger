package main

import (
	"log"
	"log/slog"
	"os"
	handlers "src/api/handlers"
	api_keycloak "src/api/keycloak"
	"src/api/middleware"
	services "src/api/service"

	// middleware "src/api/middleware"
	"src/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var logger *slog.Logger
var db *sqlx.DB
var repositoryWrapper *repositories.RepositoryWrapper

func LoadRepositoryWrapper() {
	transactionRepository := repositories.NewTransactionRepository(db.DB, logger)
	accountRepository := repositories.NewAccountRepository(db.DB, logger)
	clientRepository := repositories.NewClientRepository(db.DB, logger)
	registryAccountOtpRepository := repositories.NewRegistryAccountOtpRepository(db.DB,logger)
	repositoryWrapper = &repositories.RepositoryWrapper{
		ClientRepository:      clientRepository,
		AccountRepository:     accountRepository,
		TransactionRepository: transactionRepository,
		RegistryAccountOtpRepository: registryAccountOtpRepository,
	}
}
func initializer() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Warning: Could not load .env file: %v. Falling back to system environment variables.", err)
		panic("environment variables could not be loaded!")
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Debug level for detailed test output
	}))

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
	/**
	* Middlewares
	*/
	
	router.Use(func(ctx *gin.Context) {
		middleware.KeycloakClientMiddleware(ctx, keycloakClient)
	})
	/**
	 * HANDLERS
	 */
	clientHandler := handlers.IClientHandler{
		ClientRepository: repositoryWrapper.ClientRepository,
		RegistryAccountOtpRepository: repositoryWrapper.RegistryAccountOtpRepository,
		ClientService: services.NewClientService(repositoryWrapper.ClientRepository, repositoryWrapper.RegistryAccountOtpRepository),
	}
	accountHandler := handlers.IAccountHandler{
		KeycloakClient: keycloakClient,
		AccountService: services.NewAccountService(*repositoryWrapper),
		ClientRepository: repositoryWrapper.ClientRepository,
		AccountRepository:     repositoryWrapper.AccountRepository,
		TransactionRepository: repositoryWrapper.TransactionRepository,
		RegistryAccountOtpRepository: repositoryWrapper.RegistryAccountOtpRepository,
	}

	transactionHandler := handlers.ITransactionHandler{
		AccountRepository:     repositoryWrapper.AccountRepository,
		TransactionRepository: repositoryWrapper.TransactionRepository,
	}

	authHandler := handlers.IAuthorizationHandler {
		KeycloakClient: keycloakClient,
	}
	
	/**
	 * ROUTES
	 */
	router.GET("/")
	authorization := router.Group("/authorization")
	{
		authorization.POST("/login", authHandler.Authorization)
	}
	accounts := router.Group("/accounts")
	{
		accounts.POST("", accountHandler.CreateAccount)
		accounts.POST("/completeNewUserRegistration", accountHandler.CompleteNewUserRegistration)
		accounts.GET("/:client_id", accountHandler.FetchAccounts)
	}
	clients := router.Group("/clients")
	{
		clients.POST("", clientHandler.CreateClient)

	}
	transactions := router.Group("/transactions")
	{
		transactions.GET("/:account_id", transactionHandler.GetTransactions)	
		transactions.POST("", transactionHandler.PerformTransaction)
	}
	router.Run() // Listen on :8080 by default
	defer db.Close()
}
