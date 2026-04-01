package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const confirmHeader = "X-Action-Confirm"

func RequireActionConfirm(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}
		value := strings.TrimSpace(c.GetHeader(confirmHeader))
		if strings.EqualFold(value, "CONFIRM") {
			c.Next()
			return
		}
		abortWithError(
			c,
			http.StatusPreconditionRequired,
			"confirmation_required",
			"write action requires confirmation",
			"set header X-Action-Confirm: CONFIRM before retrying action "+action,
		)
	}
}
