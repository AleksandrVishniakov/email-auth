package main

import (
	"context"
	"github.com/AleksandrVishniakov/email-auth/app/internal/handlers"
	"github.com/AleksandrVishniakov/email-auth/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/email-auth/app/internal/repositories/user_repo"
	"github.com/AleksandrVishniakov/email-auth/app/internal/servers"
	"github.com/AleksandrVishniakov/email-auth/app/internal/services/email_service"
	"github.com/AleksandrVishniakov/email-auth/app/internal/services/user_service"
	"github.com/AleksandrVishniakov/email-auth/app/utils/configs"
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
	cfg.Email.Password = os.Getenv("EMAIL_PASSWORD")

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting...")

	db, err := postgres.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("error while db starting: %s", err.Error())
	}
	slog.Info("db started")

	var repo = user_repo.NewUserRepository(db)
	var emailService = email_service.NewEmailService(cfg.Email)
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
