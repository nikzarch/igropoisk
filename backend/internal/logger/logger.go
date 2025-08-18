package logger

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

var (
	Logger  *slog.Logger
	logFile *os.File
)

func InitLogger() error {
	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	logFile = f
	mw := io.MultiWriter(os.Stdout, f)

	handler := slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	Logger = slog.New(handler)
	return nil
}

func CloseFile() {
	if logFile != nil {
		logFile.Close()
	}
}

func SlogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		level := slog.LevelInfo
		status := c.Writer.Status()
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		Logger.Log(c.Request.Context(), level, "HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
		)
	}
}
