package cliententity

import (
    "database/sql"
    "time"
)

// Client represents the clients table in the database.
type ClientEntity struct {
    ID            int            `db:"id" json:"id"`
    Name          string         `db:"name" json:"name"`
    Surname1      string         `db:"surname1" json:"surname1"`
    Surname2      sql.NullString `db:"surname2" json:"surname2,omitempty"`
    Email         string         `db:"email" json:"email"`
    Identification string         `db:"identification" json:"identification"`
    Nationality   string         `db:"nationality" json:"nationality"`
    DateOfBirth   time.Time      `db:"date_of_birth" json:"date_of_birth"`
    Sex           string         `db:"sex" json:"sex"`
    Address       string         `db:"address" json:"address"`
    City          string         `db:"city" json:"city"`
    Province      string         `db:"province" json:"province"`
    State         sql.NullString `db:"state" json:"state,omitempty"`
    ZipCode       string         `db:"zip_code" json:"zip_code"`
    Telephone     string         `db:"telephone" json:"telephone"`
    CreatedAt     time.Time      `db:"created_at" json:"created_at"`
    UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}