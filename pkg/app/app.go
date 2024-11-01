package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	servers []Server
}

func NewApp(servers ...Server) *App {
	return &App{
		servers: servers,
	}
}

func (app *App) Run(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	for _, srv := range app.servers {
		go func(srv Server) {
			if err := srv.Start(ctx); err != nil {
				panic("error starting server: " + err.Error())
			}
		}(srv)
	}

	<-quit

	var shutdown context.CancelFunc
	ctx, shutdown = context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	for _, srv := range app.servers {
		if err := srv.Stop(ctx); err != nil {
			panic("error stopping server: " + err.Error())
		}
	}
}
