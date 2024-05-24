package middleware

import (
	"assignment-go-rest-api/constant"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId(c *gin.Context) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Set(constant.RequestId, uuid)
	c.Next()
}
