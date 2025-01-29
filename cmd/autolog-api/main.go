package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/keola-dunn/autolog/cmd/autolog-api/handler"
	"github.com/sirupsen/logrus"
)

var environmentConfig struct {
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.Info("Hello, world!")

	if err := envconfig.Process("", &environmentConfig); err != nil {
		logger.WithError(err).Fatal("failed to process environment config")
	}

	h := handler.NewHandler(handler.HandlerConfig{})

	router := newRouter(h)

	server := http.Server{
		Handler: router,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

func newRouter(h *handler.Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", h.Home)
	router.Get("/health", h.HealthCheck)

	return router
}
