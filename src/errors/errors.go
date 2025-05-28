package app_errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	return fmt.Sprintf("%s not found", e.Entity)
}

func (e *ErrNotFound) JsonError(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
}

type ErrInternalServer struct {
	Reason  error
	Message string
}

func (e *ErrInternalServer) Error() string {
	return "Internal Server Error"
}

func (e *ErrInternalServer) JsonError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
}

type ErrBadRequest struct {
	Reason  error
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
	Reason  error
	Message string
}

func (e *ErrUnauthorized) Error() string {
	return "unauthorized"
}

func (e *ErrUnauthorized) JsonError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
}

type ErrForbidden struct {
	Reason  error
	Message string
}

func (e *ErrForbidden) Error() string {
	return "forbidden"
}

func (e *ErrForbidden) JsonError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
}

// Keycloak errors

type ErrRsaPublicKey struct {
	Reason  error
	Message string
}

func (e *ErrRsaPublicKey) Error() string {
	return e.Message

}
func (e *ErrRsaPublicKey) JsonError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
}

type ErrNotJwkFound struct {
	Reason  error
}

func (e *ErrNotJwkFound) Error() string {
	return "Not JWK found"

}
func (e *ErrNotJwkFound) JsonError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
}

type ErrVerifyToken struct {
	Reason  error
    Message string
}

func (e *ErrVerifyToken) Error() string {
	return e.Message

}
func (e *ErrVerifyToken) JsonError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
}


