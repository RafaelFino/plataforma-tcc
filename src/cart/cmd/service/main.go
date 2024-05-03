package main

import (
	"cart/internal/server"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"

	"cart/internal/config"
)

func main() {
	PrintLogo()
	if len(os.Args) < 2 {
		fmt.Print("Usage: cart <config_file>\n")
		os.Exit(1)
	}

	configFile := os.Args[1]

	cfg, err := config.ConfigClientFromFile(configFile)
	if err != nil {
		fmt.Printf("Error loading config file: %s", err)
		os.Exit(1)
	}

	err = initLogger(cfg.LogPath)
	if err != nil {
		fmt.Printf("Error opening log file: %s, using stdout", err)
		log.SetOutput(os.Stdout)
	}

	log.Printf("[main] Starting with config: %s", cfg.ToJSON())

	fmt.Printf("\nStarting...\n")

	server := server.NewServer(cfg)
	go server.Run()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	log.Print("Stopping...\n")
}

func initLogger(path string) error {
	if err := os.Mkdir(path, 0755); !os.IsExist(err) {
		fmt.Printf("Error creating directory %s: %s", path, err)
		return err
	}

	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/cart-%s.log", path, "%Y%m%d"),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationCount(30), //30 days
	)
	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}

	multi := io.MultiWriter(writer, os.Stdout)
	log.SetOutput(multi)

	return nil
}

func PrintLogo() {
	fmt.Print(`
################################################################
#
# CART SERVICE 										   
#
################################################################	

`)
}
