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
	log.Info().Msg("[APP] initialise app")

	config := config.LoadConfig(".")
	log.Info().Msgf("[APP] current env: \t%s", config.App.Env)

	database.Setup(config.Database)
	database.EnsureMigrations(database.Migrations)

	api.ServePublicServer(config.Server)
	api.ServeApiDocs(config.Server)

	redis.Setup(context.TODO(), config.Redis)

	rabbit.Setup(config.Rabbit)
	rabbit.Service.RegisterConsumers([]string{})

	gracefullShutdown([]func() error{
		func() error {
			return database.DbConnection.Close()
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
	})
}

func gracefullShutdown(ops []func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			op()
		}
	}
}
