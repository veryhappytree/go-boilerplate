package main

import (
	"context"
	"fmt"
	"go-boilerplate/config"
	"go-boilerplate/pkg/api"
	"go-boilerplate/pkg/database"
	"go-boilerplate/pkg/logger"
	"go-boilerplate/pkg/rabbit"
	"go-boilerplate/pkg/redis"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Setup()

	slog.Info("[APP]", "message", "Initialize app")

	cfg := config.LoadConfig(".")
	slog.Info("[APP]", "message", fmt.Sprintf("current env: %s", cfg.App.Env))

	database.Setup(&cfg.Database)
	database.EnsureMigrations(database.Migrations)

	api.ServePublicServer(cfg.Server)
	api.ServeAPIDocs(cfg.Server)

	redis.Setup(context.TODO(), cfg.Redis)

	rabbit.Setup(cfg.Rabbit)
	rabbit.Service.RegisterConsumer("queueName1", func([]byte) {})
	rabbit.Service.RegisterConsumer("queueName2", func([]byte) {})

	gracefulShutdown(
		func() error {
			return database.DBConnection.Close()
		},
		func() error {
			return redis.Client.Close()
		},
		func() error {
			return rabbit.Service.CloseChannels()
		},
		func() error {
			os.Exit(0)
			return nil
		},
	)
}

func gracefulShutdown(ops ...func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			if err := op(); err != nil {
				slog.Error("gracefulShutdown op failed", "error", err)
				panic(err)
			}
		}
	}
}
