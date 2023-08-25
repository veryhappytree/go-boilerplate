package api

import (
	_ "embed"
	"fmt"
	"net/http"

	"go-boilerplate/config"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func ServePublicServer(cfg config.ServerConfig) {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte("boilerplate"))
	})

	httpPort := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		if err := http.ListenAndServe(httpPort, r); err != nil {
			log.Panic().AnErr("ServePublicServer http.ListenAndServe failed", err)
		}
	}()

	log.Info().Msgf("[HTTP] server is running at port: \t%d\n", cfg.Port)
}

//go:embed docs/api-docs.json
var apiDocs []byte

func ServeAPIDocs(cfg config.ServerConfig) {
	mux := http.NewServeMux()
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs", swaggerui.Handler(apiDocs)))

	httpPort := fmt.Sprintf(":%d", cfg.APIDocsPort)
	go func() {
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Panic().AnErr("ServeAPIDocs http.ListenAndServe failed", err)
		}
	}()

	log.Info().Msgf("[HTTP] api docs server is running at port: \t%d\n", cfg.APIDocsPort)
}
