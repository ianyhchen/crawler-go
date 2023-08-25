package client

import (
	"log"
	"os"
)

func SetupCustomLogger() *log.Logger {
	return log.New(os.Stdout, "[ptt_crawler] ", log.LstdFlags)
}
