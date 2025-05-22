package mappers

import
(	
	"time"
	"fmt"
	"database/sql"
	clientdto "src/api/dto"
	cliententity "src/domain/client"
)

func ToClientEntity(client clientdto.CreateClientRequest) (cliententity.ClientEntity, error) {
	var entity cliententity.ClientEntity

	// Parsear DateOfBirth
	dob, err := time.Parse("2006-01-02", client.DateOfBirth)
	if err != nil {
		return cliententity.ClientEntity{}, fmt.Errorf("error parsing DateOfBirth: %w", err)
	}

	// Parsear CreatedDate
	// Nota: El formato "2006-01-02 15:04:05" es para "YYYY-MM-DD HH:mm:ss"
	createdDate, err := time.Parse("2006-01-02 15:04:05", client.CreatedDate)
	if err != nil {
		return cliententity.ClientEntity{}, fmt.Errorf("error parsing CreatedDate: %w", err)
	}

	// Parsear UpdatedDate
	updatedDate, err := time.Parse("2006-01-02 15:04:05", client.UpdatedDate)
	if err != nil {
		return cliententity.ClientEntity{}, fmt.Errorf("error parsing UpdatedDate: %w", err)
	}

	// Asignar campos directos
	entity.Address = client.Address
	entity.City = client.City
	entity.Email = client.Email
	entity.Identification = client.Identification
	entity.Name = client.Name
	entity.Nationality = client.Nationality
	entity.Province = client.Province
	entity.Sex = client.Sex
	entity.Surname1 = client.Surname1
	entity.Telephone = client.Telephone
	entity.ZipCode = client.ZipCode

	// Asignar fechas parseadas
	entity.DateOfBirth = dob
	entity.CreatedAt = createdDate
	entity.UpdatedAt = updatedDate

	// Manejar campos opcionales usando sql.NullString
	// Si el campo no está vacío en el DTO, asigna el valor.
	// Si está vacío, sql.NullString automáticamente se marcará como NULL en la base de datos.
	if client.Surname2 != "" {
		entity.Surname2 = sql.NullString{String: client.Surname2, Valid: true}
	} else {
		entity.Surname2 = sql.NullString{Valid: false}
	}

	if client.State != "" {
		entity.State = sql.NullString{String: client.State, Valid: true}
	} else {
		entity.State = sql.NullString{Valid: false}
	}

	// Asignar ID (en un caso real, esto probablemente se generaría en la base de datos
	// después de la inserción, o sería un UUID generado aquí si tu DB lo soporta).
	// Para este ejemplo, lo inicializamos en 0.
	entity.ID = 0 // O genera un ID temporal si es necesario antes de la DB

	return entity, nil
}

func ToClientDTO(entity cliententity.ClientEntity) (clientdto.ClientResponse, error) {
	var dto clientdto.ClientResponse

	// Formatear DateOfBirth
	// El formato "2006-01-02" es para "YYYY-MM-DD"
	dto.DateOfBirth = entity.DateOfBirth.Format("2006-01-02")

	// Formatear CreatedAt
	// El formato "2006-01-02 15:04:05" es para "YYYY-MM-DD HH:mm:ss"
	dto.CreatedDate = entity.CreatedAt.Format("2006-01-02 15:04:05")

	// Formatear UpdatedAt
	dto.UpdatedDate = entity.UpdatedAt.Format("2006-01-02 15:04:05")

	// Asignar campos directos
	dto.ID = entity.ID
	dto.Address = entity.Address
	dto.City = entity.City
	dto.Email = entity.Email
	dto.Identification = entity.Identification
	dto.Name = entity.Name
	dto.Nationality = entity.Nationality
	dto.Province = entity.Province
	dto.Sex = entity.Sex
	dto.Surname1 = entity.Surname1
	dto.Telephone = entity.Telephone
	dto.ZipCode = entity.ZipCode

	// Manejar campos opcionales que usan sql.NullString
	// Si el campo es válido (no nulo en la base de datos), asigna su valor.
	if entity.Surname2.Valid {
		dto.Surname2 = entity.Surname2.String
	} else {
		dto.Surname2 = "" // O nil, dependiendo de cómo quieras representar un campo nulo en tu DTO
	}

	if entity.State.Valid {
		dto.State = entity.State.String
	} else {
		dto.State = "" // O nil
	}

	return dto, nil
}