package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	cliententity "src/domain/client"
	"src/test/utils"
	insert_entities "src/test/utils"
)

// TestDatabaseConnection verifies the database connection and performs a simple query.
func main() {
	ctx := context.Background()

	// Start the container (in a real suite, this would be in TestMain or setup)
	connStr := "user=postgres dbname=ledger sslmode=disable password=root host=localhost"
	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("failed to open database: %s\n", err)
	}
	defer db.Close()

	// Test the connection with a ping
	err = db.PingContext(ctx)

	// Test a simple query
	var exampleClient = utils.CreateClientTest(1, "Jhon", "jhon@test.es")
	var exampleClient2 = utils.CreateClientTest(2, "Joe", "joe@insertion.com")
	// Insert the client
	err = insert_entities.InsertClient(ctx, db, &exampleClient)

	err = insert_entities.InsertClient(ctx, db, &exampleClient2)

	// Verify the insertion
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM clients").Scan(&count)

	query := "SELECT * from clients"

	r, err := db.QueryContext(ctx, query)
	if err != nil {
		fmt.Printf("failed to query clients: % s", err)
	}
	clients, er := cliententity.FetchClientEntities(r)
	fmt.Println(clients)

	if er != nil {
	}

	accountJhon := insert_entities.CreateAccount(1)
	accountJoe := insert_entities.CreateAccount(2)
	fmt.Println("Account Jhon: " + string(accountJhon.AccountNumber))
	fmt.Println("Account Joe: " + string(accountJoe.AccountNumber))
	errJhonAcc := insert_entities.InsertAccount(ctx, db, &accountJhon, accountJhon.ClientID)
	if errJhonAcc != nil {
		fmt.Println("Jhon acc: " + errJhonAcc.Error())
	}
	insert_entities.InsertAccount(ctx, db, &accountJoe, accountJoe.ClientID)
	fmt.Println("Account ID Jhon: " + fmt.Sprintf("%d", accountJhon.ID))
	fmt.Println("Account ID Joe: " + fmt.Sprintf("%d", accountJoe.ID))

	joeAddsMoneyTransaction := insert_entities.CreateTransaction(accountJoe.ClientID, sql.NullInt32{
		
	}, 2000.65, "ingreso")
	joeAddsMoneyTransaction.AccountID = accountJoe.ID
	txErr := insert_entities.AccountTransactionTx(ctx, db, &joeAddsMoneyTransaction, accountJoe.ClientID, "CREDITO", "")
	if txErr != nil {
		fmt.Println("Error " + txErr.Error())
	}

	joeTransfersMoneyTransactionToJhon := insert_entities.CreateTransaction(
		accountJoe.ID, 
		sql.NullInt32{
			Int32: int32(accountJhon.ID),
			Valid: true,
			}, 
		1000, 
		"TRANSFERENCIA",
	)
	txErr = insert_entities.AccountTransactionTx(ctx, db, &joeTransfersMoneyTransactionToJhon, accountJoe.ClientID, "DEBITO", "CREDITO")
	if txErr != nil {
		fmt.Println("Error " + txErr.Error())
	}

}
