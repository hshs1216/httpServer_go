package main

import (
	"fmt"
	"log"
	"os"
)

func setupLogger() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
}
