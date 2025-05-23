package app_errors

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)


// type ErrNotEnoughFunds struct {
//     Message string
// }

// type ErrEntityNotFound struct {
//     Identifier any
// }

// type ErrTransactionInsertFailed struct {
//     Message string
//     Reason  error // Wraps original DB error
// }

// type ErrLedgerEntryInsertFailed struct {
//     Message string
//     Reason  error // Wraps original DB error
// }

// type ErrNoRowsAffected struct {
//     Message string
//     Reason  error // Wraps original DB error
// }

// func (e *ErrEntityNotFound) Error() string {
//     return fmt.Sprintf("Entity not found: %v", e.Identifier)
// }

// func (e *ErrTransactionInsertFailed) Error() string {
//     return fmt.Sprintf("Failed transaction insertion: %s", e.Message)
// }

// func (e *ErrTransactionInsertFailed) Unwrap() error {
//     return e.Reason
// }

// func (e *ErrLedgerEntryInsertFailed) Error() string {
//     return fmt.Sprintf("Failed ledger entry insertion: %s", e.Message)
// }

// func (e *ErrLedgerEntryInsertFailed) Unwrap() error {
//     return e.Reason
// }

// func (e *ErrNoRowsAffected) Error() string {
//     return fmt.Sprintf("No rows affected: %s", e.Message)
// }

// func (e *ErrNoRowsAffected) Unwrap() error {
//     return e.Reason
// }

// func (e *ErrNotEnoughFunds) Error() string {
//     return fmt.Sprintf(e.Message)
// }

/**
* Error types
*
*/

type AppError interface {
    Error() string
    JsonError(c *gin.Context)
}

type ErrNotFound struct {
    Reason error
    Entity string
}
func (e *ErrNotFound) Error() string {
    return fmt.Sprintf("%s not found",e.Entity)
}

func (e *ErrNotFound) JsonError(c *gin.Context) {
    c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
}


type ErrInternalServer struct {
    Reason error
    Message string
}

func (e *ErrInternalServer) Error() string {
    return "Internal Server Error"
}

func (e *ErrInternalServer) JsonError(c *gin.Context) {
    c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
}

type ErrBadRequest struct {
    Reason error
    Message string
}

func (e *ErrBadRequest) Error() string {
    return "bad request"
}

func (e *ErrBadRequest) JsonError(c *gin.Context) {
    c.JSON(http.StatusBadRequest, gin.H{"error": e.Error(), "message": e.Message})
}

type ErrNotEnoughFunds struct {
    Message string
}
func (e *ErrNotEnoughFunds) Error() string {
    return "not enough funds"
}

func (e *ErrNotEnoughFunds) JsonError(c *gin.Context) {
    c.JSON(http.StatusConflict, gin.H{"error": e.Error(), "message": e.Message})
}
type ErrUnauthorized struct {
    Reason error
    Message string
}

func (e *ErrUnauthorized) Error() string {
    return "unauthorized"
}

func (e *ErrUnauthorized) JsonError(c *gin.Context) {
    c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
}

type ErrForbidden struct {
    Reason error
    Message string
}
func (e *ErrForbidden) Error() string {
    return "forbidden"
}

func (e *ErrForbidden) JsonError(c *gin.Context) {
    c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
}