package client

import (
	"log"
	"os"
)

func SetupCustomLogger() *log.Logger {
	return log.New(os.Stdout, "[ptt_crawler] ", log.LstdFlags)
}

//func CustomLoggerMiddleware() gin.HandlerFunc {
//	logger := setupCustomLogger()
//
//	return func(c *gin.Context) {
//		c.Set("logger", logger) // Store the logger in the Gin context for later use
//		c.Next()
//	}
//}
