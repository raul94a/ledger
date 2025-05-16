package main

import (
	"log"
	clientdto "src/api/dto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
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
    r := gin.Default()
	
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

    account:=clientdto.AccountDto{
        ID: 1004564563,
        Balance: 65665.33,
        UserID: "77",
        AccountNumber: "447474788564",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

	r.GET("/")
    r.GET("/ping/1", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "ping",
			"underage": isUnderAge,
        })
    })
	r.GET("/ping/2", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
			"underage": isUnderAge2,
        })
    })
    r.GET("/accounts/1", func(c *gin.Context){
        c.JSON(200, account)
    })
    r.Run() // Listen on :8080 by default
}