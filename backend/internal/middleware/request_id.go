package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

const CtxRequestIDKey = "km_request_id"

func InjectRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set(CtxRequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "req-" + time.Now().Format("20060102150405")
	}
	return "req-" + hex.EncodeToString(buf)
}
