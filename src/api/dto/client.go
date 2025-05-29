package clientdto


import
(
	    "time"
)

type CreateClientRequest struct {
 	Address        string `json:"address" binding:"required"`
    City           string `json:"city" binding:"required"`
    DateOfBirth    string `json:"date_of_birth" binding:"required,datetime=2006-01-02"` // ISO 8601 date (YYYY-MM-DD)
    Email          string `json:"email" binding:"required,email"`
    Identification string `json:"identification" binding:"required"` // e.g., passport, national ID
    Name           string `json:"name" binding:"required"`
    Nationality    string `json:"nationality" binding:"required"` // e.g., ISO 3166-1 alpha-2 code
    Province       string `json:"province" binding:"required"`
    Sex            string `json:"sex" binding:"required,oneof=M F"` // M(ale), F(emale)
    State          string `json:"state"` // Optional, depending on country
    Surname1       string `json:"surname1" binding:"required"`
    Surname2       string `json:"surname2"` // Optional, as not all have second surname
    TaxID          string `json:"tax_id"` // e.g., SSN, TIN
    Telephone      string `json:"telephone" binding:"required"` // Basic phone validation can be added
    ZipCode        string `json:"zip_code" binding:"required"`
}

type ClientResponse struct {
    ID             int    `json:"id" binding:"required"`
    OTP            string `json:"otp"`
 	Address        string `json:"address" binding:"required"`
    City           string `json:"city" binding:"required"`
    DateOfBirth    string `json:"date_of_birth" binding:"required,datetime=2006-01-02"` // ISO 8601 date (YYYY-MM-DD)
    Email          string `json:"email" binding:"required,email"`
    Identification string `json:"identification" binding:"required"` // e.g., passport, national ID
    Name           string `json:"name" binding:"required"`
    Nationality    string `json:"nationality" binding:"required"` // e.g., ISO 3166-1 alpha-2 code
    Province       string `json:"province" binding:"required"`
    Sex            string `json:"sex" binding:"required,oneof=M F"` // M(ale), F(emale)
    State          string `json:"state"` // Optional, depending on country
    Surname1       string `json:"surname1" binding:"required"`
    Surname2       string `json:"surname2"` // Optional, as not all have second surname
    TaxID          string `json:"tax_id" binding:"required"` // e.g., SSN, TIN
    Telephone      string `json:"telephone" binding:"required"` // Basic phone validation can be added
    ZipCode        string `json:"zip_code" binding:"required"`
    CreatedDate    string `json:"created_date" binding:"required,datetime=2006-01-02 15:04:05"` // ISO 8601 date (YYYY-MM-DD HH:mm:ss)
    UpdatedDate    string `json:"updated_date" binding:"required,datetime=2006-01-02 15:04:05"` // ISO 8601 date (YYYY-MM-DD HH:mm:ss)
}



// IsUnderage returns true if the client is under 18 years old
func (r CreateClientRequest) IsUnderage() (bool, error) {
    dob, err := time.Parse("2006-01-02", r.DateOfBirth)
    if err != nil {
        return false, err
    }
    now := time.Now()
    age := now.Year() - dob.Year()
    if now.YearDay() < dob.YearDay() {
        age--
    }
    return age < 18, nil
}