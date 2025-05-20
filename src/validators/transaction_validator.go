package validators

import (
	"fmt"
	"log/slog"
	transaction_entity "src/domain/transaction"
	errors "src/errors"
	"strings"
)

func ValidateTransactionBalance(transaction transaction_entity.TransactionEntity, balance float64, logger *slog.Logger) error {
	transactionType := strings.ToUpper(transaction.Type)

	if transactionType == "ADD" {
		return nil
	}
	if balance < transaction.Amount {
		errStr := fmt.Sprintf(
			"Not enough funds. Account %d has %v money. Tried to %s %v units",
			transaction.AccountID,
			balance,
			transactionType,
			transaction.Amount,
		)
		logger.Error(errStr)
		return &errors.ErrNotEnoughFunds{Message: errStr}
	}

	return nil

}
