package utils

import (
	"database/sql"
	cliententity "src/domain/client"
	"time"
)

func CreateClientTest(ID int, name string, email string) cliententity.ClientEntity {
	return cliententity.ClientEntity{
		ID:             ID,
		Name:           name,
		Surname1:       "Doe",
		Surname2:       sql.NullString{String: "Smith", Valid: true}, // Optional, can be {Valid: false} for NULL
		Email:          email,
		Identification: "ABC123456" + string(ID),
		Nationality:    "US",
		DateOfBirth:    time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Sex:            "M",
		Address:        "123 Main St",
		City:           "Springfield",
		Province:       "Illinois",
		State:          sql.NullString{String: "IL", Valid: true}, // Optional, can be {Valid: false} for NULL
		ZipCode:        "62701",
		Telephone:      "+1-555-123-4567",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
}
