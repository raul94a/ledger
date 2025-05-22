package main

import (
    "os"
	"log"
	"log/slog"
	clientdto "src/api/dto"
	handlers "src/api/handlers"
	middleware "src/api/middleware"
	"src/repositories"
	"time"

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

	req := clientdto.CreateClientRequest{
		Name:           "Maria",
		Surname1:       "Garcia",
		Surname2:       "Lopez",
		Email:          "maria.garcia@example.com",
		Identification: "X12345678Z",
		TaxID:          "123456789",
		Nationality:    "ES",
		DateOfBirth:    "1995-03-10",
		Sex:            "F",
		Address:        "123 Banking Street",
		City:           "Madrid",
		Province:       "Madrid",
		State:          "", // Optional, not used in Spain
		ZipCode:        "28001",
		Telephone:      "+34 612 345 678",
	}
	req2 := clientdto.CreateClientRequest{
		Name:           "Maria",
		Surname1:       "Garcia",
		Surname2:       "Lopez",
		Email:          "maria.garcia@example.com",
		Identification: "X12345678Z",
		TaxID:          "123456789",
		Nationality:    "ES",
		DateOfBirth:    "2015-03-10",
		Sex:            "F",
		Address:        "123 Banking Street",
		City:           "Madrid",
		Province:       "Madrid",
		State:          "", // Optional, not used in Spain
		ZipCode:        "28001",
		Telephone:      "+34 612 345 678",
	}
	isUnderAge, _ := req.IsUnderage()
	isUnderAge2, _ := req2.IsUnderage()

	account := clientdto.AccountDto{
		ID:            1004564563,
		Balance:       65665.33,
		ClientID:      "77",
		AccountNumber: "447474788564",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

    /**
    * HANDLERS
    */
   
	clientHandler := handlers.IClientHandler{
		ClientRepository: repositories.NewClientRepository(db.DB, logger),
	}
    /**
    * ROUTES
    */
	router.GET("/")
	router.GET("/ping/1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "ping",
			"underage": isUnderAge,
		})
	})
	router.GET("/ping/2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "pong",
			"underage": isUnderAge2,
		})
	})
	router.GET("/accounts/1", func(c *gin.Context) {
		c.JSON(200, account)
	})
    router.POST("/clients",clientHandler.CreateClient)
	router.Run() // Listen on :8080 by default
}
