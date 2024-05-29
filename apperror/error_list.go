package apperror

import (
	"assignment-go-rest-api/constant"
	"fmt"
	"net/http"
	"strconv"
)

func (c *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", c.Code, c.Message)
}

func ErrInternalServer() error {
	return NewCustomError(http.StatusInternalServerError, constant.ResponseMsgErrorInternalServer)
}

func ErrBadRequest() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgErrorInvalidRequest)
}

func ErrUserNotFound() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgErrorUserNotFound)
}

func ErrEmailRequired() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgEmailRequired)
}

func ErrPasswordRequired() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgPasswordRequired)
}

func ErrUnauthorized() error {
	return NewCustomError(http.StatusUnauthorized, constant.ResponseMsgErrorUnauthorized)
}

func ErrForbidden() error {
	return NewCustomError(http.StatusForbidden, constant.ResponseMsgErrorForbidden)
}

func ErrEmailAlreadyRegistered() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgEmailAlreadyRegistered)
}

func ErrPasswordShouldEightChar() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgPasswordShouldEightChar)
}

func ErrIncorrectCredentials() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgIncorrectCredentials)
}

func ErrWalletAlreadyCreated() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgWalletAlreadyCreated)
}

func ErrWalletNotFound() error {
	return NewCustomError(http.StatusNotFound, constant.ResponseMsgWalletNotFound)
}

func ErrWalletRecipientNotFound() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgRecipienWalletNotFound)
}

func ErrSourceFundNotFound() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgSourceFundNotFound)
}

func ErrInsufficientBalance() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgInsufficientBalance)

}

func ErrTopUpFailed() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgTopUpFailed)
}

func ErrTransferFailed() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgTransferFailed)
}

func ErrCantTransferToOwnWallet() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgCantTransferToOwnWallet)
}

func ErrResetTokenNotFound() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgResetTokenNotFound)
}

func ErrPasswordNotMatch() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgPasswordNotMatch)
}

func ErrLimitTopUp() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgTopUpFailed+"! top up should be between "+strconv.Itoa(constant.MinTopUp)+" and "+strconv.Itoa(constant.MaxTopUp))
}

func ErrLimitTransfer() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgTransferFailed+"! transfer should be between "+strconv.Itoa(constant.MinTransfer)+" and "+strconv.Itoa(constant.MaxTransfer))
}

func ErrTransactionNotFound() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgTransactionNotFound)
}

func ErrTokenExpired() error {
	return NewCustomError(http.StatusUnauthorized, constant.ResponseMsgErrorTokenExpired)
}

func ErrInvalidIntegerInput() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgInvalidIntInput)
}

func ErrInvalidZeroLimitInput() error {
	return NewCustomError(http.StatusBadRequest, constant.ResponseMsgInvalidZeroLimitInput)
}
