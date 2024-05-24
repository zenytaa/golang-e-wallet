package utils

import "github.com/gin-gonic/gin"

func GetQuery(c *gin.Context, queryKey string) (string, bool) {
	return c.GetQuery(queryKey)
}
