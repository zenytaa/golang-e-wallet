package app

// import (
// 	"assignment-go-rest-api/repository"
// 	"assignment-go-rest-api/usecase"
// 	"assignment-go-rest-api/utils"
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"time"
// )

// type RouterOptSeeder struct {
// 	Authusecase        usecase.AuthUsecase
// 	Transactionusecase usecase.TransactionUsecase
// 	Userusecase        usecase.UserUsecase
// }

// func RegisSeeder(routerSedd *RouterOptSeeder) {
// 	ctx := context.Background()
// 	for _, v := range UserData() {
// 		res, err := routerSedd.Authusecase.Register(ctx, v)
// 		// time.Sleep(1 * time.Second)
// 		fmt.Println("user : ", res)
// 		utils.IfErrorLogPrint(err)
// 	}

// 	for _, v := range TransactionTopUpData() {
// 		res, err := routerSedd.Transactionusecase.TopUp(ctx, v)
// 		time.Sleep(1 * time.Second)
// 		fmt.Println("topup : ", res)
// 		utils.IfErrorLogPrint(err)
// 	}

// 	for _, v := range TransactionTransferData() {
// 		res, err := routerSedd.Transactionusecase.Transfer(ctx, v)
// 		time.Sleep(1 * time.Second)
// 		fmt.Println("transfer : ", res)
// 		utils.IfErrorLogPrint(err)
// 	}

// }

// func RunSeeder(db *sql.DB) {
// 	userRepository := repository.NewUserRepository()
// 	walletRepository := repository.NewWalletRepository()
// 	passwordResetRepository := repository.NewPasswordResetRepository(db)
// 	sourceFundRepository := repository.NewSourceFundRepository(db)
// 	transactionRepository := repository.NewTransactionRepository()

// 	authUsecase := usecase.NewAuthUsecase(db, userRepository, walletRepository, passwordResetRepository)
// 	userUsecase := usecase.NewUserUsecase(db, userRepository, walletRepository, transactionRepository)
// 	// walletUsecase := usecase.NewWalletUsecase(db, userRepository, walletRepository)
// 	transactionUsecase := usecase.NewTransactionUsecase(db, transactionRepository, walletRepository, sourceFundRepository, userRepository)

// 	RegisSeeder(&RouterOptSeeder{
// 		Authusecase:        authUsecase,
// 		Userusecase:        userUsecase,
// 		Transactionusecase: transactionUsecase,
// 	})
// }
