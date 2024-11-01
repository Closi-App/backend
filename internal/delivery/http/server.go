package http

import (
	"context"
	"fmt"
	"github.com/Closi-App/backend/internal/delivery/http/v1"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

type Server struct {
	*fiber.App
	address string
}

func NewServer(cfg *viper.Viper, log *logger.Logger, handler *v1.Handler) *Server {
	engine := fiber.New(fiber.Config{
		AppName:      cfg.GetString("app.name"),
		ReadTimeout:  cfg.GetDuration("http.read_timeout"),
		WriteTimeout: cfg.GetDuration("http.write_timeout"),
		IdleTimeout:  cfg.GetDuration("http.idle_timeout"),
	})

	engine.Use(
		fiberzerolog.New(fiberzerolog.Config{
			Logger: &log.Logger,
		}),
		recover.New(),
	)

	handler.InitRoutes(engine.Group("/api"))

	return &Server{
		App:     engine,
		address: fmt.Sprintf("%s:%d", cfg.GetString("http.host"), cfg.GetInt("http.port")),
	}
}

func (s *Server) Start(context.Context) error {
	return s.Listen(s.address)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.ShutdownWithContext(ctx)
}
