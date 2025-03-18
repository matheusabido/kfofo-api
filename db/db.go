package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Instance *pgxpool.Pool

func SetupDB() {
	url := os.Getenv("DB_URL")

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signalChan
		fmt.Println("Fechando conexões...")
		cancel()
		fmt.Println("Conexões fechadas.")
		os.Exit(0)
	}()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatal("Could not connect to the database.")
	}

	Instance = pool
}
