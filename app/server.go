package app

import (
	"assignment-go-rest-api/handler"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/usecase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func createRouter(db *sql.DB) *gin.Engine {
	userRepository := repository.NewUserRepository(&repository.UserRepoOpts{Db: db})
	walletRepository := repository.NewWalletRepository(&repository.WalletRepoOpts{Db: db})
	passwordResetRepository := repository.NewPasswordResetRepository(db)
	// sourceFundRepository := repository.NewSourceFundRepository(db)

	authUsecase := usecase.NewAuthUsecase(&usecase.AuthUsecaseOpts{
		UserRepository:          userRepository,
		WalletRepository:        walletRepository,
		PasswordResetRepository: passwordResetRepository,
		Transactor:              repository.NewTransactor(db),
	})
	userUsecase := usecase.NewUserUsecase(&usecase.UserUsecaseOpts{
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
	})
	walletUsecase := usecase.NewWalletUsecase(&usecase.WalletUsecaseOpts{
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
