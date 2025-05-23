package utils

import (
	"fmt"
	app_utils "src/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	const spainCode string = "ES"
	// BBBB
	const myBankDigits string = "0182"
	// GGGG
	const myBranchBankDigits string = "0600"

	const iterations int = 10
	for i := 0; i < iterations; i++ {
		handler := app_utils.IbanHandler{}
		accNr := handler.GenerateAccountNumber(10)
		cc := handler.DomesticCheckDigits(myBankDigits, myBranchBankDigits, accNr)
		bban := app_utils.Bban {
			BankCode: myBankDigits,
			BranchCode: myBranchBankDigits,
			DomesticCheckDigits: cc,
			AccountNumber: accNr,
		}
		iban, _ := handler.ComputeIban(bban,spainCode)
		isValidIban := handler.Verify(iban)
		assert.True(t,isValidIban)
		fmt.Printf("IBAN: %s %s %s %s %s ; Is Valid? %v\n", iban[0:4], iban[4:8], iban[8:12], iban[12:14], iban[14:], isValidIban)
	}

}

