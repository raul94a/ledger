package handlers

import (
	"net/http"
	dto "src/api/dto"
	api_keycloak "src/api/keycloak"
	appRedis "src/db/redis"
	utils "src/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AuthorizationHandler interface {
	Authorization(c *gin.Context)
}

type IAuthorizationHandler struct {
	KeycloakClient api_keycloak.KeycloakClient
	Logger         *zap.Logger
	RedisClient    *redis.Client
}

func (h *IAuthorizationHandler) Authorization(c *gin.Context) {
	// 0. Bind the request body
	var req dto.AuthorizationRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. User sends data with session-id request header - else, login
	if sessionIdReq := c.Request.Header.Get("session-id"); sessionIdReq != "" {
		tokenResponse, responseType := appRedis.GetToken(c.Request.Context(), h.RedisClient, sessionIdReq)
		switch responseType.(type) {
		case appRedis.LoginAgainResponse:
			h.login(c, req, "")
		case appRedis.ErrorResponse:
			c.AbortWithStatus(500)
		case appRedis.OkResponse:
			err := h.KeycloakClient.VerifyIssuer(tokenResponse.AccessToken, req.User)
			if err != nil {
				h.login(c, req, "")
				return
			}
			if h.KeycloakClient.HasTokenExpired(tokenResponse.AccessToken) {
				h.login(c, req, sessionIdReq)
				return
			}
			h.redisUpsertNewSessionId(c, sessionIdReq, tokenResponse)
			c.JSON(200, tokenResponse)
		}
	} else {
		h.login(c, req, sessionIdReq)
	}

}

func (h *IAuthorizationHandler) login(c *gin.Context, req dto.AuthorizationRequest, sessionIdReq string) {
	token, err := h.KeycloakClient.AuthUser(req.User, req.Password)
	if err != nil {
		err.JsonError(c)
		return
	}
	h.redisUpsertNewSessionId(c, sessionIdReq, &token)
	c.JSON(http.StatusOK, token)
}

func (h *IAuthorizationHandler) redisUpsertNewSessionId(c *gin.Context, sessionIdReq string, token *api_keycloak.TokenResponse) {
	// session-id
	sessionId, hashError := utils.CreateSessionId()
	if hashError != nil {
		h.Logger.Error(hashError.Error())
		c.AbortWithStatus(500)
		return
	}
	ctx := c.Request.Context()
	// delete previous session id if applies
	if sessionIdReq != "" {
		appRedis.DeleteToken(ctx, h.RedisClient, sessionIdReq)
	}
	// register new session ID in redis
	appRedis.SetToken(ctx, h.RedisClient, sessionId, token)
	c.Header("session-id", sessionId)

}
