package repository_Test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"src/repositories"
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

func TestRepository(t *testing.T) {
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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Debug level for detailed test output
	}))

	clientRepository := repositories.NewClientRepository(db, logger)
	accountRepository := repositories.NewAccountRepository(db, logger)
	transactionRepository := repositories.NewTransactionRepository(db, logger)
	// Test a simple query
	var exampleClient = utils.CreateClientTest(1, "Jhon", "jhon@test.es")
	var exampleClient2 = utils.CreateClientTest(2, "Joe", "joe@insertion.com")
	// Insert the client
	err = clientRepository.InsertClient(ctx, &exampleClient)
	assert.NoError(t, err, "failed to insert client")
	err = clientRepository.InsertClient(ctx, &exampleClient2)
	assert.NoError(t, err, "failed to insert client 2")

	// Verify some returned fields
	assert.NotZero(t, exampleClient.ID, "expected non-zero ID")
	assert.False(t, exampleClient.CreatedAt.IsZero(), "expected non-zero created_at")
	assert.False(t, exampleClient.UpdatedAt.IsZero(), "expected non-zero updated_at")

	accountJhon := insert_entities.CreateAccount(1)
	accountJoe := insert_entities.CreateAccount(2)

	accountRepository.InsertAccount(ctx, &accountJhon)
	accountRepository.InsertAccount(ctx, &accountJoe)

	toAccountId := sql.NullInt32{}
	joeAddsMoneyTransaction := insert_entities.CreateTransaction(accountJoe.ID, toAccountId, 2000.65, "ADD")
	txErr := transactionRepository.InsertTransactionLedgerTx(ctx, &joeAddsMoneyTransaction)
	assert.NoError(t, txErr)

	toAccountIdJhon := sql.NullInt32{Int32: int32(accountJhon.ID), Valid: true}

	joeTransfersMoneyTransactionToJhon := insert_entities.CreateTransaction(accountJoe.ID, toAccountIdJhon, 1000, "TRANSFERENCIA")
	txErr = transactionRepository.InsertTransactionLedgerTx(ctx, &joeTransfersMoneyTransactionToJhon)
	assert.NoError(t, txErr)

	balance, err := transactionRepository.FetchAccountBalance(ctx,nil, accountJhon.ID)
	t.Logf("Jhon account balance %v",*balance)
	balance, err = transactionRepository.FetchAccountBalance(ctx,nil, accountJoe.ID)
	t.Logf("Joe account balance %v",*balance)
	


}
