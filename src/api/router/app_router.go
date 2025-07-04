package app_router

import (
	"src/api/handlers"
	api_keycloak "src/api/keycloak"
	"src/api/middleware"
	services "src/api/service"
	"src/repositories"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

)

type AppRouter struct {
	RepositoryWrapper *repositories.RepositoryWrapper
	KeycloakClient    *api_keycloak.KeycloakClient
	RedisClient		  *redis.Client
	ZapLogger		  *zap.Logger
}

func (appRouter AppRouter) BuildRoutes(router *gin.Engine) {
	/**
	* Middlewares
	 */

	router.Use(func(ctx *gin.Context) {
		middleware.KeycloakClientMiddleware(ctx, *appRouter.KeycloakClient)

	})

	router.Use(func(ctx *gin.Context) {
		middleware.RepositoryWrapperMiddleware(ctx, appRouter.RepositoryWrapper)
	})

	authHandlerMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			middleware.AuthorizationMiddleware(c)
		}
	}

	logger := middleware.RequestLogger()

	clientHandler := handlers.IClientHandler{
		ClientRepository:             appRouter.RepositoryWrapper.ClientRepository,
		RegistryAccountOtpRepository: appRouter.RepositoryWrapper.RegistryAccountOtpRepository,
		ClientService:                services.NewClientService(appRouter.RepositoryWrapper.ClientRepository, appRouter.RepositoryWrapper.RegistryAccountOtpRepository),
	}
	accountHandler := handlers.IAccountHandler{
		KeycloakClient:               *appRouter.KeycloakClient,
		AccountService:               services.NewAccountService(*appRouter.RepositoryWrapper),
		ClientRepository:             appRouter.RepositoryWrapper.ClientRepository,
		AccountRepository:            appRouter.RepositoryWrapper.AccountRepository,
		TransactionRepository:        appRouter.RepositoryWrapper.TransactionRepository,
		RegistryAccountOtpRepository: appRouter.RepositoryWrapper.RegistryAccountOtpRepository,
	}

	transactionHandler := handlers.ITransactionHandler{
		AccountRepository:     appRouter.RepositoryWrapper.AccountRepository,
		TransactionRepository: appRouter.RepositoryWrapper.TransactionRepository,
	}

	authHandler := handlers.IAuthorizationHandler{
		KeycloakClient: *appRouter.KeycloakClient,
		Logger: appRouter.ZapLogger,
		RedisClient: appRouter.RedisClient,
	}
	/**
	 * ROUTES
	 */
	
	router.GET("/")
	authorization := router.Group("/authorization")
	{
		authorization.POST("/login",logger,authHandler.Authorization)
	}

	accounts := router.Group("/accounts")
	{
		accounts.POST("", logger,authHandlerMiddleware(), accountHandler.CreateAccount)
		accounts.POST("/completeNewUserRegistration", logger,accountHandler.CompleteNewUserRegistration)
		// Verificar el client ID en el middleware de auth
		accounts.GET(
			"/:client_id",
			authHandlerMiddleware(),
			middleware.AuthenticationByClientIdHandler(),
			accountHandler.FetchAccounts,
		)
	}
	clients := router.Group("/clients")
	{
		// Verificar que identificacion corresponde al clientID
		clients.GET(
			"/:identification",
			authHandlerMiddleware(),
			middleware.AuthenticateUserByIdentificationHandler(),
			clientHandler.GetClientByIdentification,
		)
		// Este endpoint debe recibir algún token especial para la autorización
		clients.POST("",logger, clientHandler.CreateClient)

	}
	transactions := router.Group("/transactions",logger, authHandlerMiddleware())
	{
		// verificar que la cuenta corresponda al cliente
		transactions.GET(
			"/:account_id",
			middleware.AuthenticateByAccountIdHandler(),
			transactionHandler.GetTransactions,
		)
		// verificar que la cuenta corresponda al cliente
		transactions.POST(
			"",logger,
		    middleware.AuthenticatePerformTransactionHandler(),
			transactionHandler.PerformTransaction,
		)
	}
}
