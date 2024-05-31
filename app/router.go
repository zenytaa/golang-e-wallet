package app

import (
	"assignment-go-rest-api/handler"
	"assignment-go-rest-api/middleware"
	"assignment-go-rest-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouterOpt struct {
	AuthHandler   *handler.AuthHandler
	UserHandler   *handler.UserHandlerImpl
	Transaction   *handler.TransactionHandler
	PasswordReset *handler.PasswordResetHandler
}

func NewRouter(routerOpt *RouterOpt, config utils.Config) *gin.Engine {
	router := gin.Default()
	router.ContextWithFallback = true

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestId)
	router.Use(middleware.Logger(log))
	router.Use(middleware.ErrorHandler())

	router.POST("/register", routerOpt.AuthHandler.Register)
	router.POST("/login", routerOpt.AuthHandler.Login)

	protected := router.Group("/users")
	{
		protected.Use(middleware.JWTAuthMiddleware(config))
		protected.GET("/profiles", routerOpt.UserHandler.GetProfile)
		protected.POST("/transfer", routerOpt.Transaction.Transfer)
		protected.POST("/top-up", routerOpt.Transaction.TopUp)
		protected.GET("/transaction-lists", routerOpt.Transaction.GetListTransaction)
		protected.POST("/forgot-password", routerOpt.PasswordReset.ForgotPassword)
		protected.POST("/reset-password", routerOpt.PasswordReset.PasswordReset)
	}
	router.Use()

	return router
}
