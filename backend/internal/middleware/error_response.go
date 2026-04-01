package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, status int, code, message, hint string) {
	requestID := c.GetString(CtxRequestIDKey)
	resp := gin.H{
		"error":     message,
		"code":      code,
		"requestId": requestID,
	}
	if hint != "" {
		resp["hint"] = hint
	}
	if status == 0 {
		status = http.StatusBadRequest
	}
	c.Set("km_error", message)
	c.JSON(status, resp)
	c.Abort()
}
