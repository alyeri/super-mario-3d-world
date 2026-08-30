package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := startLocalNNCS(); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
}
