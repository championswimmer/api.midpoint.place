package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/championswimmer/api.midpoint.place/src/server"
	"github.com/samber/lo"
)

func main() {
	server := server.CreateServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	go func() {
		<-quit
		// print shutdown message
		log.Println("Shutting down server...") // TODO: use appLogger

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.ShutdownWithContext(ctx)
	}()

	lo.Must0(server.Listen(":3000")) // TODO: pick from env.PORT
}
