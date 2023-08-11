package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGHUP)
	defer cancel()

	if err := execute(ctx); err != nil {
		log.Fatalf("Error running the App. %v \n", err)
	}
}
