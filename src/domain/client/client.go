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
    KcUserId      sql.NullInt64  `db:"kc_user_id"`
}

func ScanClientEntity(r *sql.Rows, client *ClientEntity) error {
    return r.Scan(
			&client.ID,
			&client.Name,
			&client.Surname1,
			&client.Surname2,
			&client.Email,
			&client.Identification,
			&client.Nationality,
			&client.DateOfBirth,
			&client.Sex,
			&client.Address,
			&client.City,
			&client.Province,
			&client.State,
			&client.ZipCode,
			&client.Telephone,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
}

func FetchClientEntities(r *sql.Rows) ([]ClientEntity,error) {
    var clients []ClientEntity
    defer r.Close()
    for r.Next(){
        var client ClientEntity
        er := ScanClientEntity(r,&client)
		if er != nil {
			return nil,er
		}
		clients = append(clients, client)
    }

    return clients,nil
}