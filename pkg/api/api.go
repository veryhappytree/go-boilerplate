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

func ServePublicServer(config config.ServerConfig) {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)
}

//go:embed docs/api-docs.json
var apiDocs []byte

func ServeApiDocs(config config.ServerConfig) {
	mux := http.NewServeMux()
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs", swaggerui.Handler(apiDocs)))

	httpPort := fmt.Sprintf(":%d", config.ApiDocsPort)
	go http.ListenAndServe(httpPort, mux)

	log.Info().Msgf("[HTTP] api docs server is running at port: \t%d\n", config.ApiDocsPort)
}
