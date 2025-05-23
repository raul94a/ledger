package utils

import(
	"fmt"
	"strconv"
	"math/rand"
)



type Bban struct {
    BankCode string
    BranchCode string
    DomesticCheckDigits string
    AccountNumber string
}

func (b Bban) String() string {
    return fmt.Sprintf("%s%s%s%s",b.BankCode, b.BranchCode, b.DomesticCheckDigits, b.AccountNumber)
}

type Iban struct {
    CountryCode string
    CheckDigits string    
    BankCode string
    BranchCode string
    DomesticCheckDigits string
    AccountNumber string
}

func (b Iban) String() string {
    return fmt.Sprintf("%s%s%s%s%s%s",b.CountryCode,b.CheckDigits,b.BankCode, b.BranchCode, b.DomesticCheckDigits, b.AccountNumber)
}


type IbanActions interface {
    GenerateAccountNumber(length int) string
    ComputeDomesticCheckDigits(bankDigits, branchDigits, accountNr string) (string, error)
    ComputeIban(bban Bban, countryCode string) string
    VerifyIban(iban string) string
}

type IbanHandler struct {}

func (i *IbanHandler) GenerateAccountNumber(length int) string {
	alphabet := "0123456789"
	seed := 55600746159796123
	rand.NewSource(int64(seed))
	strLen := len(alphabet) - 1
	accNumber := ""
	for i := 0; i < length; i++{
		random := rand.Intn(strLen)
		char := string(alphabet[random])
		accNumber = accNumber + char
	}

	return accNumber

} 

// Domestic Check Digits (CC):
// Calculated using weighted sums:
// For BBBBGGGG: Weights [1, 2, 4, 8, 5, 10, 9, 7], modulo 11.

// For AAAAAAAAAA: Weights [1, 2, 4, 8, 5, 10, 9, 7, 3, 6], modulo 11.

// Adjust results: 10 → 1, 11 → 0.

func (i *IbanHandler) DomesticCheckDigits(bankDigits, branchDigits, accountNr string) (string, error) {

	// BBBBGGGG Construction
	bankDigitsWithBranchDigits := fmt.Sprintf("%s%s", bankDigits, branchDigits)
	// fmt.Println(tag + " BBBBGGGG: ", bankDigitsWithBranchDigits)

	weights := []int{1, 2, 4, 8, 5, 10, 9, 7}

	sum := 0
	for i, weight := range weights {
		n, _ := strconv.Atoi(string(bankDigitsWithBranchDigits[i]))
		sum += n * weight
	}

	firstDigit := 11 - (sum % 11)
	if firstDigit == 11 {
		firstDigit = 0
	} else if firstDigit == 10 {
		firstDigit = 1
	}

	weights = []int{1, 2, 4, 8, 5, 10, 9, 7, 3, 6}
	sum = 0
	for i, weight := range weights {
		digit, _ := strconv.Atoi(string(accountNr[i]))
		sum += digit * weight
	}
	secondDigit := 11 - (sum % 11)
	if secondDigit == 11 {
		secondDigit = 0
	} else if secondDigit == 10 {
		secondDigit = 1
	}

	return fmt.Sprintf("%v%v", firstDigit, secondDigit), nil

}

func (i *IbanHandler) numericBban(bbanCountry string) string {
	numericIBAN := ""
	for _, char := range bbanCountry {
		if char >= 'A' && char <= 'Z' {
			numericIBAN += fmt.Sprintf("%d", char-'A'+10)
		} else {
			numericIBAN += string(char)
		}
	}
	return numericIBAN
}

// **MOD 97-10 is the algorithm used for IBAN generation**
//
// [iban] parameter corresponds to the following iban structure:
//
// BBBBGGGGCCAAAAAAAAAAXX00
//  * BBBB is the Bank Code
//  * GGGG is the Branch Code for the bank
//  * AAAAAAAAAA is the Account number
//  * XX is the country code (Two letters code)
//  * 00 is the placeholder
// However, in this case the country code is not the two letters code,
// but the position in the alphabet of each letter plus ten. Example: 
//
//      The letter E corresponds to 14 (position 4 + 10 = 14)
//      The letter S corresponds to 28 (position 18 + 10 = 28)
//      So, the two letter code (ES) is translated to 1428, as the 
//      values are concatenated. Along with the placeholder (00), FOR
//      this case the value is BBBBGGGGCCAAAAAAAAAA142800
func (i *IbanHandler) mod_97_10(iban string) (string, error) {
	// Calculate MOD-97-10
	const chunkSize = 9
	remainder := 0
	for i := 0; i < len(iban); i += chunkSize {
		end := i + chunkSize
		if end > len(iban) {
			end = len(iban)
		}
		chunk := strconv.Itoa(remainder) + iban[i:end]
		num, err := strconv.Atoi(chunk)
		if err != nil {
			return "", fmt.Errorf("error converting chunk to integer: %v", err)
		}
		remainder = num % 97
	}

	// Calculate IBAN check digits
	checkDigitsValue := 98 - remainder
	if checkDigitsValue < 0 {
		checkDigitsValue += 97
	}
	checkDigitsStr := fmt.Sprintf("%02d", checkDigitsValue)
	return checkDigitsStr, nil
}

// Spanish IBAN Structure: ESkkBBBBGGGGCCAAAAAAAAAA (24 characters)
// MOD-97-10 Algorithm for kk:
// The check digits (kk) are calculated by:
// Constructing the BBAN: BBBBGGGGCCAAAAAAAAAA.

// Forming the string: BBBBGGGGCCAAAAAAAAAAES00.

// Converting letters to numbers
// (E=14, S=28, digits unchanged).

// Computing the remainder
// of the numeric string modulo 97.

// Check digits = 98 - remainder,
// formatted as two digits.
func (i *IbanHandler) ComputeIban(bban Bban, countryCode string) (string, error) {
	
    const placeholder string = "00"
	bbanCountry := fmt.Sprintf("%s%s%s", bban.String(), countryCode, placeholder)
	// Convert letters to numbers (A=10, B=11, ..., Z=35)
	numericIBAN := i.numericBban(bbanCountry)
	checkDigitsStr, _ := i.mod_97_10(numericIBAN)
	// Construct final IBAN
	ibanStr := Iban {
        CountryCode: countryCode,
        CheckDigits: checkDigitsStr,
        BankCode: bban.BankCode,
        BranchCode: bban.BranchCode,
        DomesticCheckDigits: bban.DomesticCheckDigits,
        AccountNumber: bban.AccountNumber,
    }.String()
	if len(ibanStr) != 24 {
		return "", fmt.Errorf("generated IBAN has invalid length: %d", len(ibanStr))
	}
	return ibanStr, nil

}

//
// ESkk BBBB GGGG CC AAAAAAAAAA
func (i *IbanHandler) Verify(iban string) bool {
	if len(iban) < 24 {
		return false
	}
	bankDigits := iban[4:8]
	branchDigits := iban[8:12]
	accountNr := iban[14:]
	ibanControlDigits := iban[12:14]

	controlDigits, err := i.DomesticCheckDigits(bankDigits, branchDigits, accountNr)
	if err != nil || controlDigits != ibanControlDigits {
		fmt.Printf("\nBank digits:%s\tBranch Digits:%s\tAccountNr:%s\n", bankDigits, branchDigits, accountNr)
		fmt.Printf("Computed ControlDigits %s : TestControlDigits %s\n", controlDigits, ibanControlDigits)
		return false
	}
	country := iban[0:2]
	bban := iban[4:]
	placeholder := "00"
	// BBBBGGGGccAAAAAAAAAAES00
	bbanCountry := fmt.Sprintf("%s%s%s", bban, country, placeholder)
	bbanNumeric := i.numericBban(bbanCountry)
	checkDigits, err := i.mod_97_10(bbanNumeric)
	if err != nil {
        
		return false
	}
	testControlDigits := iban[2:4]
	return checkDigits == testControlDigits

}
