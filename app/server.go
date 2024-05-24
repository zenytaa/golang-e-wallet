package app

import (
	"assignment-go-rest-api/handler"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/usecase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func createRouter(db *sql.DB) *gin.Engine {
	userRepository := repository.NewUserRepository()
	walletRepository := repository.NewWalletRepository()
	passwordResetRepository := repository.NewPasswordResetRepository(db)
	// sourceFundRepository := repository.NewSourceFundRepository(db)

	authUsecase := usecase.NewAuthUsecase(&usecase.AuthUsecaseOpts{
		Db:                      db,
		UserRepository:          userRepository,
		WalletRepository:        walletRepository,
		PasswordResetRepository: passwordResetRepository,
	})
	userUsecase := usecase.NewUserUsecase(&usecase.UserUsecaseOpts{
		Db:               db,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
	})
	walletUsecase := usecase.NewWalletUsecase(&usecase.WalletUsecaseOpts{
		Db:               db,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
	})

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(walletUsecase)

	return NewRouter(&RouterOpt{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		UserUsecase: userUsecase,
	})
}

func InitServer(db *sql.DB) *gin.Engine {
	router := createRouter(db)
	return router
}
