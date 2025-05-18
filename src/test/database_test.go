package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	cliententity "src/domain/client"
	"src/test/utils"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestMain sets up the test environment and runs tests.
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start the PostgreSQL container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:17-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container after tests
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Run tests
	m.Run()
}

func InsertClient(ctx context.Context, db *sql.DB, client *cliententity.ClientEntity) error {
	query := `
        INSERT INTO clients (
            name, surname1, surname2, email, identification, nationality, 
            date_of_birth, sex, address, city, province, state, 
            zip_code, telephone
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := db.QueryRowContext(ctx, query,
		client.Name,
		client.Surname1,
		client.Surname2,
		client.Email,
		client.Identification,
		client.Nationality,
		client.DateOfBirth,
		client.Sex,
		client.Address,
		client.City,
		client.Province,
		client.State,
		client.ZipCode,
		client.Telephone,
	).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert client: %w", err)
	}
	return nil
}

// TestDatabaseConnection verifies the database connection and performs a simple query.
func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	// Start the container (in a real suite, this would be in TestMain or setup)
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:17-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("../migrations/00001_tables.up.sql"),
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
	err = InsertClient(ctx, db, &exampleClient)

	assert.NoError(t, err, "failed to insert client")
    err = InsertClient(ctx,db, &exampleClient2)
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

}
