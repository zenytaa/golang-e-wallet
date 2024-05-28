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
	AuthHandler *handler.AuthHandlerImpl
	UserHandler *handler.UserHandlerImpl
	Transaction *handler.TransactionHandler
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
	router.POST("/forgot-password", routerOpt.AuthHandler.ForgotPassword)
	router.POST("/reset-password", routerOpt.AuthHandler.ResetPassword)

	protected := router.Group("/user")
	{
		protected.Use(middleware.JWTAuthMiddleware(config))
		protected.GET("/profile", routerOpt.UserHandler.GetProfile)
		protected.POST("/transfer", routerOpt.Transaction.Transfer)
	}
	router.Use()

	return router
}
