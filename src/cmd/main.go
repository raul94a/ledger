package main

import (
	"log"
	"log/slog"
	"os"
	handlers "src/api/handlers"
	middleware "src/api/middleware"
	"src/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title API Bank Clients
// @version 1.0
// @description Clients management of a bank system.
// @host localhost:8080
// @BasePath /
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Debug level for detailed test output
	}))
	urlString := "user=postgres dbname=ledger sslmode=disable password=root host=localhost"
	db, err := sqlx.Connect("postgres", urlString)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	router := gin.Default()
	// Register Middlewares!
	for _, middlewareFn := range middleware.AppMiddlewares() {
		router.Use(middlewareFn())
	}

	
	


	/**
	 * HANDLERS
	 */

	clientHandler := handlers.IClientHandler{
		ClientRepository: repositories.NewClientRepository(db.DB, logger),
	}
	accountHandler := handlers.IAccountHandler {
		AccountRepository: repositories.NewAccountRepository(db.DB,logger),
		TransactionRepository: repositories.NewTransactionRepository(db.DB,logger),
	}
	
	/**
	 * ROUTES
	 */
	router.GET("/")
	router.GET("/accounts/:client_id", accountHandler.FetchAccounts)
	router.POST("/clients", clientHandler.CreateClient)
	router.Run() // Listen on :8080 by default
}
