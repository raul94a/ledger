package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	handlers "src/api/handlers"
	api_keycloak "src/api/keycloak"
	"src/api/middleware"
	services "src/api/service"
	"strconv"

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
		middleware.RepositoryWrapperMiddleware(ctx, repositoryWrapper)
	})

	authHandlerMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			middleware.AuthorizationMiddleware(c)
		}
	}

	/**
	 * HANDLERS
	 */
	clientHandler := handlers.IClientHandler{
		ClientRepository:             repositoryWrapper.ClientRepository,
		RegistryAccountOtpRepository: repositoryWrapper.RegistryAccountOtpRepository,
		ClientService:                services.NewClientService(repositoryWrapper.ClientRepository, repositoryWrapper.RegistryAccountOtpRepository),
	}
	accountHandler := handlers.IAccountHandler{
		KeycloakClient:               keycloakClient,
		AccountService:               services.NewAccountService(*repositoryWrapper),
		ClientRepository:             repositoryWrapper.ClientRepository,
		AccountRepository:            repositoryWrapper.AccountRepository,
		TransactionRepository:        repositoryWrapper.TransactionRepository,
		RegistryAccountOtpRepository: repositoryWrapper.RegistryAccountOtpRepository,
	}

	transactionHandler := handlers.ITransactionHandler{
		AccountRepository:     repositoryWrapper.AccountRepository,
		TransactionRepository: repositoryWrapper.TransactionRepository,
	}

	authHandler := handlers.IAuthorizationHandler{
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
		accounts.POST("", authHandlerMiddleware(), accountHandler.CreateAccount)
		accounts.POST("/completeNewUserRegistration", accountHandler.CompleteNewUserRegistration)
		// Verificar el client ID en el middleware de auth
		accounts.GET("/:client_id", authHandlerMiddleware(), func(c *gin.Context) {
			clientIDStr := c.Param("client_id")

			clientID, err := strconv.Atoi(clientIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid identifier"})
				return
			}
			middleware.AuthenticateUserByClientIdMiddleware(c, clientID)
		}, accountHandler.FetchAccounts)
	}
	clients := router.Group("/clients")
	{
		// Verificar que identificacion corresponde al clientID
		clients.GET("/:identification", authHandlerMiddleware(), func(c *gin.Context) {
			identification := c.Param("identification")
			middleware.AuthenticateUserByIdentificationMiddleware(c, identification)
		}, clientHandler.GetClientByIdentification)
		// Este endpoint debe recibir algún token especial para la autorización
		clients.POST("", clientHandler.CreateClient)

	}
	transactions := router.Group("/transactions", authHandlerMiddleware())
	{
		// verificar que la cuenta corresponda al cliente
		transactions.GET("/:account_id", func(c *gin.Context) {
			accountIdStr := c.Param("account_id")

			accountID, err := strconv.Atoi(accountIdStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid identifier"})
				return
			}
			middleware.AuthenticateUserByAccountIdMiddleware(c, accountID)
		}, transactionHandler.GetTransactions)
		// verificar que la cuenta corresponda al cliente
		transactions.POST("", transactionHandler.PerformTransaction)
	}
	router.Run() // Listen on :8080 by default
	defer db.Close()
}
