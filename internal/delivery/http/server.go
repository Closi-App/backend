package http

import (
	"context"
	"fmt"
	"github.com/Closi-App/backend/internal/config"
	"github.com/Closi-App/backend/internal/delivery/http/v1"
	"github.com/Closi-App/backend/internal/logger"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Server struct {
	*fiber.App
	address string
}

func NewServer(cfg *config.Config, log logger.Logger, handler *v1.Handler) *Server {
	engine := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	})

	zerologLogger := log.Logger().(zerolog.Logger)

	engine.Use(
		fiberzerolog.New(fiberzerolog.Config{ // TODO: implement custom logger middleware
			Logger: &zerologLogger,
		}),
		recover.New(),
	)

	handler.InitRoutes(engine.Group("/api"))

	return &Server{
		App:     engine,
		address: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (s *Server) Start(context.Context) error {
	if err := s.Listen(s.address); err != nil {
		return errors.Wrap(err, "error starting http server")
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.ShutdownWithContext(ctx); err != nil {
		return errors.Wrap(err, "error stopping http server")
	}
	return nil
}
