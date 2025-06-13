package app_router

import (
	"src/api/handlers"
	api_keycloak "src/api/keycloak"
	services "src/api/service"
	"src/api/middleware"
	"src/repositories"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	RepositoryWrapper *repositories.RepositoryWrapper
	KeycloakClient    *api_keycloak.KeycloakClient
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
	}
	/**
	 * ROUTES
	 */
	router.GET("/")
	authorization := router.Group("/authorization")
	{
		authorization.POST("/login", authHandler.Authorization)
	}
	accounts := router.Group("/accounts")
	{
		accounts.POST("", authHandlerMiddleware(), accountHandler.CreateAccount)
		accounts.POST("/completeNewUserRegistration", accountHandler.CompleteNewUserRegistration)
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
		clients.POST("", clientHandler.CreateClient)

	}
	transactions := router.Group("/transactions", authHandlerMiddleware())
	{
		// verificar que la cuenta corresponda al cliente
		transactions.GET(
			"/:account_id",
			middleware.AuthenticateByAccountIdHandler(),
			transactionHandler.GetTransactions,
		)
		// verificar que la cuenta corresponda al cliente
		transactions.POST(
			"",
		    middleware.AuthenticatePerformTransactionHandler(),
			transactionHandler.PerformTransaction,
		)
	}
}
