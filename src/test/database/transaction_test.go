package transaction_test

import (
	"context"
	"database/sql"
	"fmt"
	cliententity "src/domain/client"
	"src/test/utils"
	insert_entities "src/test/utils"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDatabaseConnection verifies the database connection and performs a simple query.
func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	// Start the container (in a real suite, this would be in TestMain or setup)
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:17-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("../../migrations/00001_tables.up.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get the connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Test the connection with a ping
	err = db.PingContext(ctx)
	assert.NoError(t, err, "failed to ping database")

	// Test a simple query
	var exampleClient = utils.CreateClientTest(1, "Jhon", "jhon@test.es")
    var exampleClient2 = utils.CreateClientTest(2,"Joe", "joe@insertion.com")
	// Insert the client
	err = insert_entities.InsertClient(ctx, db, &exampleClient)

	assert.NoError(t, err, "failed to insert client")
    err = insert_entities.InsertClient(ctx,db, &exampleClient2)
    assert.NoError(t, err, "failed to insert client 2")

	// Verify the insertion
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM clients").Scan(&count)
	assert.NoError(t, err)
    t.Logf("Number of rows %d",count)
	assert.Equal(t, 2, count, "expected two clients records")

	// Verify some returned fields
	assert.NotZero(t, exampleClient.ID, "expected non-zero ID")
	assert.False(t, exampleClient.CreatedAt.IsZero(), "expected non-zero created_at")
	assert.False(t, exampleClient.UpdatedAt.IsZero(), "expected non-zero updated_at")

	query := "SELECT * from clients"

	r, err := db.QueryContext(ctx, query)
	var clients []cliententity.ClientEntity 
	if err != nil {
        t.Fatalf("failed to query clients: % s", err)
	}
    clients, er := cliententity.FetchClientEntities(r)
	defer r.Close()
	
	
    if er != nil {
        t.Fatalf("failed to scan client: %s", er)
	}

	t.Log(clients)
	assert.Equal(t, len(clients), 2)
	accountJhon := insert_entities.CreateAccount(1)
	accountJoe := insert_entities.CreateAccount(2)
	fmt.Println("Account Jhon: " + string(accountJhon.AccountNumber))
	fmt.Println("Account Joe: " + string(accountJoe.AccountNumber))
	insert_entities.InsertAccount(ctx,db,&accountJhon, accountJhon.ClientID)
	insert_entities.InsertAccount(ctx,db,&accountJoe, accountJoe.ClientID)
	fmt.Println("Account ID Jhon: " + fmt.Sprintf("%d", accountJhon.ID))
	fmt.Println("Account ID Joe: " + fmt.Sprintf("%d", accountJoe.ID))
	toAccountId := sql.NullInt32{}
	joeAddsMoneyTransaction := insert_entities.CreateTransaction(accountJoe.ID,toAccountId,2000.65, "ingreso")
	txErr := insert_entities.AccountTransactionTx(ctx,db, &joeAddsMoneyTransaction, accountJoe.ClientID, "CREDITO", "")
	assert.NoError(t, txErr)
	
	toAccountIdJhon := sql.NullInt32{Int32: int32(accountJhon.ID),Valid: true}
	joeTransfersMoneyTransactionToJhon := insert_entities.CreateTransaction(accountJoe.ID, toAccountIdJhon, 1000, "TRANSFERENCIA")
	txErr = insert_entities.AccountTransactionTx(ctx,db, &joeTransfersMoneyTransactionToJhon, accountJoe.ClientID, "DEBITO", "CREDITO")
	assert.NoError(t, txErr)




}

