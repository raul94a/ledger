package main

import (
	"log"
	"log/slog"
	"os"
	handlers "src/api/handlers"
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
	repositoryWrapper = &repositories.RepositoryWrapper{
		ClientRepository:      clientRepository,
		AccountRepository:     accountRepository,
		TransactionRepository: transactionRepository,
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

	defer db.Close()

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
	router := gin.Default()
	/**
	 * HANDLERS
	 */
	clientHandler := handlers.IClientHandler{
		ClientRepository: repositoryWrapper.ClientRepository,
	}
	accountHandler := handlers.IAccountHandler{
		AccountRepository:     repositoryWrapper.AccountRepository,
		TransactionRepository: repositoryWrapper.TransactionRepository,
	}

	transactionHandler := handlers.ITransactionHandler{
		AccountRepository:     repositoryWrapper.AccountRepository,
		TransactionRepository: repositoryWrapper.TransactionRepository,
	}
	/**
	 * ROUTES
	 */
	router.GET("/")
	accounts := router.Group("/accounts")
	{
		accounts.POST("", accountHandler.CreateAccount)
		accounts.GET("/:client_id", accountHandler.FetchAccounts)
	}
	clients := router.Group("/clients")
	{
		clients.POST("", clientHandler.CreateClient)

	}
	transactions := router.Group("/transactions")
	{
		transactions.POST("", transactionHandler.PerformTransaction)
	}
	router.Run() // Listen on :8080 by default
}
