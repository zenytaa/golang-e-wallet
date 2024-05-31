package dto

import "assignment-go-rest-api/entity"

type PaginationResponse struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func ToPaginationResponse(pagination entity.PaginationInfo) *PaginationResponse {
	return &PaginationResponse{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		TotalData: pagination.TotalData,
		TotalPage: pagination.TotalPage,
	}
}
