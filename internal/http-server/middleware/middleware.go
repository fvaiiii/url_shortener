package middleware

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("RequestID", requestID)

		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}

func URLFormat() gin.HandlerFunc {
	return func(c *gin.Context) {
		ext := path.Ext(c.Request.URL.Path)
		if ext != "" {
			format := strings.TrimPrefix(ext, ".")
			c.Request.URL.Path = strings.TrimSuffix(c.Request.URL.Path, ext)
			c.Set("url_format", format)
		}
		c.Next()
	}
}
