package mwlogger

import (
	"log"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%d] %s %s | %s | IP: %s",
			status,
			c.Request.Method,
			c.Request.URL.Path,
			latency,
			c.ClientIP(),
		)
	}
}

func New(log *slog.Logger) gin.HandlerFunc {

	log = log.With(slog.String("component", "middleware/logger"))
	log.Info("logger middleware enabled")

	return func(c *gin.Context) {

		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("request_id", c.GetString("RequestID")),
		)

		start := time.Now()

		c.Next()

		entry.Info("request completed",
			slog.Int("status", c.Writer.Status()),
			slog.Int("bytes", c.Writer.Size()),
			slog.String("duration", time.Since(start).String()),
		)
	}
}
