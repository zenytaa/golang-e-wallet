package constant

const (
	ResponseMsgOK                        = "ok"
	ResponseMsgErrorNotFound             = "data not found"
	ResponseMsgErrorInvalidRequest       = "invalid request, please check your request"
	ResponseMsgErrorInternalServer       = "our server encounter error, please contact us again"
	ResponseMsgErrorPageNotFound         = "page not found"
	ResponseMsgErrorMethodNotAllowed     = "method not allowed"
	ResponseMsgErrorUserNotFound         = "user not found"
	ResponseMsgErrorUnauthorized         = "unauthorized"
	ResponseMsgErrorForbidden            = "forbidden"
	ResponseMsgEmailRequired             = "email required"
	ResponseMsgEmailAlreadyRegistered    = "email already registered"
	ResponseMsgPasswordRequired          = "password required"
	ResponseMsgPasswordShouldEightChar   = "password should be minimum 8 character"
	ResponseMsgIncorrectCredentials      = "login failed! incorrect email or password"
	ResponseMsgLoginSuccess              = "login successfully!"
	ResponseMsgLoginFailed               = "failed to login"
	ResponseMsgRegisterSuccess           = "successfully registered"
	ResponseMsgRegisterFailed            = "failed to register"
	ResponseMsgWalletAlreadyCreated      = "wallet already created"
	ResponseMsgWalletNotFound            = "wallet not found"
	ResponseMsgSourceFundNotFound        = "source of fund not supported"
	ResponseMsgInsufficientBalance       = "insufficient balance"
	ResponseMsgTopUpSucces               = "successfully top up"
	ResponseMsgTopUpFailed               = "failed to top up"
	ResponseMsgTransferSucces            = "successfully transfer money"
	ResponseMsgTransferFailed            = "failed to transfer money"
	ResponseMsgCantTransferToOwnWallet   = "invalid walletnumber! you can not transfer to your own wallet"
	ResponseMsgShowProfileSucces         = "successfully show profile"
	ResponseMsgShowProfileFailed         = "failed show profile"
	ResponseMsgResetTokenNotFound        = "reset token invalid"
	ResponseMsgPasswordNotMatch          = "password not match"
	ResponseMsgGetListTransactionSuccess = "success to get your transaction list"
	ResponseMsgGetListTransactionFailed  = "failed to get your transaction list"
	ResponseMsgForgotPasswordSuccess     = "successfully forgot password"
	ResponseMsgForgotPasswordFailed      = "failed to forgot password"
	ResponseMsgResetPasswordSuccess      = "successfully reset password"
	ResponseMsgResetPasswordFailed       = "failed to reset password"
	ResponseMsgGetProfileSuccsess        = "successfully get profile"
	ResponseMsgTransactionNotFound       = "transaction not found"
	ResponseMsgErrorTokenExpired         = "token already expired"
	ResponseMsgInvalidIntInput           = "invalid integer input"
	ResponseMsgInvalidZeroLimitInput     = "limit should be more tham 0"
)
