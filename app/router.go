package app

import (
	"assignment-go-rest-api/handler"
	"assignment-go-rest-api/middleware"
	"assignment-go-rest-api/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouterOpt struct {
	AuthHandler *handler.AuthHandlerImpl
	UserHandler *handler.UserHandlerImpl
	UserUsecase usecase.UserUsecase
}

func NewRouter(routerOpt *RouterOpt) *gin.Engine {
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

	protected := router.Group("/api")
	{
		protected.Use(middleware.JwtAuthMiddleware(routerOpt.UserUsecase))
		protected.GET("/profile", routerOpt.UserHandler.GetProfile)
	}
	router.Use()

	return router
}
