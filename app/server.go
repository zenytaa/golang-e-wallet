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
	sourceFundRepository := repository.NewSourceFundRepository(db)
	transactionRepository := repository.NewTransactionRepository()

	authUsecase := usecase.NewAuthUsecase(db, userRepository, walletRepository, passwordResetRepository)
	userUsecase := usecase.NewUserUsecase(db, userRepository, walletRepository, transactionRepository)
	walletUsecase := usecase.NewWalletUsecase(db, userRepository, walletRepository)
	transactionUsecase := usecase.NewTransactionUsecase(db, transactionRepository, walletRepository, sourceFundRepository, userRepository)

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(walletUsecase, transactionUsecase)
	transactionHandler := handler.NewTransactionHandler(&handler.TransactionHandlerConfig{
		TransactionUsecase: transactionUsecase,
		UserUsecase:        userUsecase,
		AuthUsecase:        authUsecase,
		WalletUsecase:      walletUsecase,
	})

	return NewRouter(&RouterOpt{
		AuthHandler:        authHandler,
		TransactionHandler: transactionHandler,
		UserHandler:        userHandler,
		UserUsecase:        userUsecase,
	})
}

func InitServer(db *sql.DB) *gin.Engine {
	router := createRouter(db)
	return router
}
