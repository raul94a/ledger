


using migrate:

migrate -database "postgres://postgres:root@localhost:5432/ledger?sslmode=disable" -path ./migrations up
