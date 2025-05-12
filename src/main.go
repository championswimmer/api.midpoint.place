package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/server"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
)

func main() {
	// initialize db
	appDb := db.GetAppDB()

	// initialize server
	server := server.CreateServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	go func() {
		<-quit
		// print shutdown message
		applogger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.ShutdownWithContext(ctx)

		// shutdown db
		applogger.Info("Closing db connection...")
		lo.Must0(lo.Must(appDb.DB()).Close())
	}()

	applogger.Info("Starting server...")
	lo.Must0(server.Listen(fmt.Sprintf(":%s", config.Port)))
}
