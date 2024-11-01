package app

import (
	"context"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	name    string
	log     *logger.Logger
	servers []Server
}

func NewApp(cfg *viper.Viper, log *logger.Logger, servers ...Server) *App {
	return &App{
		name:    cfg.GetString("app.name"),
		log:     log,
		servers: servers,
	}
}

func (app *App) Run(ctx context.Context) {
	app.log.Info().
		Msgf("🚀 Starting %s", app.name)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	for _, srv := range app.servers {
		go func(srv Server) {
			if err := srv.Start(ctx); err != nil {
				app.log.Error().
					Err(err).
					Msg("error starting server")
			}
		}(srv)
	}

	<-quit

	var shutdown context.CancelFunc
	ctx, shutdown = context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	for _, srv := range app.servers {
		if err := srv.Stop(ctx); err != nil {
			app.log.Error().
				Err(err).
				Msg("error stopping server")
		}
	}
}
