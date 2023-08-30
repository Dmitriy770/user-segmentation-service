package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Dmitriy770/user-segmentation-service/internal/config"
	"github.com/Dmitriy770/user-segmentation-service/internal/db/postgres"
	"github.com/Dmitriy770/user-segmentation-service/internal/http-server/handlers/segment/add"
	"github.com/Dmitriy770/user-segmentation-service/internal/http-server/handlers/segment/delete"
	"github.com/Dmitriy770/user-segmentation-service/internal/http-server/handlers/user/update"
	mwLogger "github.com/Dmitriy770/user-segmentation-service/internal/http-server/middleware/logger"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/handlers/slogpretty"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/Dmitriy770/user-segmentation-service/internal/serevices/segments"
	"github.com/Dmitriy770/user-segmentation-service/internal/serevices/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	os.Exit(0)

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting user-segmentation-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.PostgreSQL)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	segmentsRep := segments.NewRepository(log, storage)
	segmentsService := segments.NewService(log, segmentsRep)
	usersRep := users.NewRepository(log, storage)
	usersService := users.NewService(log, usersRep)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/segment", add.New(log, segmentsService))
	router.Delete("/segment", delete.New(log, segmentsService))
	router.Post("/user", update.New(log, usersService))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stoped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
