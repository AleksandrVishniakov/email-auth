package main

import (
	"context"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/handlers"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/repositories/user_repo"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/servers"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/services/email_service"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/services/user_service"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/utils/configs"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	envInit()

	cfg := configs.MustConfigs()
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting...")

	db, err := postgres.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("error while db starting: %s", err.Error())
	}
	slog.Info("db started")

	var repo = user_repo.NewUserRepository(db)
	var emailService = email_service.NewEmailService(
		os.Getenv("SENDER_EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
	)

	var userService = user_service.NewUserService(repo, emailService)
	var handler = handlers.NewHTTPHandler(userService)

	server := servers.NewHTTPServer(cfg.HTTP, handler.InitRoutes())

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		slog.Info("server started on port " + cfg.HTTP.Port)
		return server.Run()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		slog.Error("server gracefully stopped", "reason", err)
	}
}

func envInit() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}
}
