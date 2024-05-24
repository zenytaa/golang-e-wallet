package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/usecase"
	"assignment-go-rest-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandlerImpl struct {
	walletUsecase usecase.WalletUsecase
}

func NewUserHandler(walletUsecase usecase.WalletUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{
		walletUsecase: walletUsecase,
	}
}

func (h *UserHandlerImpl) GetProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(*entity.User)
	walletRequest := dto.WalletRequest{}

	walletRequest.UserId = user.Id
	wallet, err := h.walletUsecase.GetWalletByUserId(ctx, &walletRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	query := &dto.ListTransactionQuery{}
	sortBy, sortByExist := utils.GetQuery(ctx, "sort_by")
	sort, sortExist := utils.GetQuery(ctx, "sort")
	search, searchExist := utils.GetQuery(ctx, "search")

	query.SortBy = sortBy
	query.SortByExist = sortByExist
	query.Sort = sort
	query.SortExist = sortExist
	query.Search = search
	query.SearchExist = searchExist

	query.SortBy = strings.ToLower(query.SortBy)
	query.Sort = strings.ToUpper(query.Sort)

	userDetailResponse := dto.ToUserDetailResponse(*user, *wallet)

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgGetProfileSuccsess,
		Data:    userDetailResponse,
	})
}
