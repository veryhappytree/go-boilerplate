package main

import (
	"context"
	"go-boilerplate/config"
	"go-boilerplate/pkg/api"
	"go-boilerplate/pkg/database"
	"go-boilerplate/pkg/logger"
	"go-boilerplate/pkg/rabbit"
	"go-boilerplate/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.SetupZeroLog()
	log.Info().Msg("[APP] initialize app")

	cfg := config.LoadConfig(".")
	log.Info().Msgf("[APP] current env: \t%s", cfg.App.Env)

	database.Setup(&cfg.Database)
	database.EnsureMigrations(database.Migrations)

	api.ServePublicServer(cfg.Server)
	api.ServeAPIDocs(cfg.Server)

	redis.Setup(context.TODO(), cfg.Redis)

	rabbit.Setup(cfg.Rabbit)
	rabbit.Service.RegisterConsumers([]string{})

	gracefullShutdown(
		func() error {
			return database.DBConnection.Close()
		},
		func() error {
			return redis.Client.Close()
		},
		func() error {
			return rabbit.Service.Channel.Close()
		},
		func() error {
			log.Warn().Msg("[APP] shutting down service")
			os.Exit(0)
			return nil
		},
	)
}

func gracefullShutdown(ops ...func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			if err := op(); err != nil {
				log.Panic().AnErr("gracefullShutdown op failed", err)
			}
		}
	}
}
