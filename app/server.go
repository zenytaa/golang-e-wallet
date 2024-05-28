package app

import (
	"assignment-go-rest-api/handler"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/usecase"
	"assignment-go-rest-api/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func createRouter(config utils.Config, db *sql.DB) *gin.Engine {
	userRepository := repository.NewUserRepository(&repository.UserRepoOpts{Db: db})
	walletRepository := repository.NewWalletRepository(&repository.WalletRepoOpts{Db: db})
	passwordResetRepository := repository.NewPasswordResetRepository(db)
	sourceFundRepository := repository.NewSourceFundRepository(db)
	transactionRepository := repository.NewTransactionRepository(&repository.TransactionRepoOpts{Db: db})

	authUsecase := usecase.NewAuthUsecase(&usecase.AuthUsecaseOpts{
		UserRepository:          userRepository,
		WalletRepository:        walletRepository,
		PasswordResetRepository: passwordResetRepository,
		Transactor:              repository.NewTransactor(db),
		AuthTokenProvider:       utils.NewJwtProvider(config),
	})
	walletUsecase := usecase.NewWalletUsecase(&usecase.WalletUsecaseOpts{
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
	})
	transactionUsecase := usecase.NewTransactionUsecase(&usecase.TransactionUsecaseOpts{
		TransactionRepo: transactionRepository,
		SourceFundRepo:  sourceFundRepository,
		WalletRepo:      walletRepository,
		Transactor:      repository.NewTransactor(db),
	})

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(walletUsecase)
	transactionHandler := handler.NewTransactionHandler(&handler.TransactionHandlerOpts{TransactionUsecase: transactionUsecase})

	return NewRouter(&RouterOpt{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		Transaction: transactionHandler,
	}, config)
}

func InitServer(config utils.Config, db *sql.DB) *gin.Engine {
	router := createRouter(config, db)
	return router
}
