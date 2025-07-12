package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/keola-dunn/autolog/cmd/autolog-api/handlers/auth"
	"github.com/keola-dunn/autolog/cmd/autolog-api/handlers/signup"
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

	authHandler := auth.NewAuthHandler()

	signupHandler := signup.NewSignupHandler()

	router := newRouter(authHandler, signupHandler)

	server := http.Server{
		Handler: router,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Autolog API!"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

func newRouter(authHandler *auth.AuthHandler, signupHandler *signup.SignupHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", home)
	router.Get("/health", healthCheck)

	router.Route("/v1", func(router chi.Router) {
		router.Route("/auth", func(router chi.Router) {
			router.Post("/login", authHandler.Login)
		})

		router.Post("/signup", signupHandler.SignUp)
	})

	return router
}
