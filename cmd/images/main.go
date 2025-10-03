package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/keola-dunn/autolog/cmd/images/internal/handlers/images"
	"github.com/keola-dunn/autolog/internal/jwt"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/image"
)

var environmentConfig struct {
	DBUser     string `envconfig:"DB_USER"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBHost     string `envconfig:"DB_HOST"`
	DBPort     int64  `envconfig:"DB_PORT"`
	DBSchema   string `envconfig:"DB_SCHEMA"`

	//AuthAPIHost string `envconfig:"AUTH_API_HOST"`

	JWKSUrl string `envconfig:"JWKS_URL"`
}

func main() {
	logger := logger.NewLogger()

	// attempt to retrieve env vars from env file. This is for local dev only
	if err := godotenv.Load(); err != nil {
		logger.Error("failed to load .env file", err)
	}

	if err := envconfig.Process("", &environmentConfig); err != nil {
		logger.Fatal("failed to process environment config", err)
	}

	///////////////////////////////////////
	// Platform and Foundational configs //
	///////////////////////////////////////
	logger.Info("connecting to the database...")
	db, err := postgres.NewConnectionPool(context.Background(), postgres.ConnectionPoolConfig{
		ConnectionConfig: postgres.ConnectionConfig{
			User:     environmentConfig.DBUser,
			Password: environmentConfig.DBPassword,
			Host:     environmentConfig.DBHost,
			Port:     environmentConfig.DBPort,
			Schema:   environmentConfig.DBSchema,
		},
		MaxConnections:        5,
		MinConnections:        1,
		MaxConnectionIdleTime: time.Minute,
	})
	if err != nil {
		logger.Fatal("failed to connect to the database", err)
	}
	defer db.Close()
	logger.Info("successfully connected to the database!")

	randomSvc := random.NewService()

	// calendarSvc := calendar.NewService()

	jwtVerifier, err := jwt.NewTokenVerifier(context.Background(), jwt.TokenVerifierConfig{
		JWKSUrl: environmentConfig.JWKSUrl,
	})
	if err != nil {
		logger.Fatal("failed to create new jwt verifier", err)
	}

	///////////////////////
	// Service Creations //
	///////////////////////

	imageSvc := image.NewService(image.ServiceConfig{
		ImagePrefix:     "images",
		DB:              db,
		RandomGenerator: randomSvc,
	})

	// userSvc := user.NewService(user.ServiceConfig{
	// 	DB:              db,
	// 	RandomGenerator: randomSvc,
	// })

	// carSvc := car.NewService(car.ServiceConfig{
	// 	DB:              db,
	// 	RandomGenerator: randomSvc,
	// })

	///////////////////////////
	// API Handler Creations //
	///////////////////////////

	authHandler, err := jwt.NewAuthHandler(jwt.AuthHandlerConfig{
		TokenVerifier: jwtVerifier,
	})
	if err != nil {
		logger.Fatal("failed to create auth handler", err)
	}

	imagesHandler, err := images.NewHandler(images.ImagesHandlerConfig{
		ImageService: imageSvc,
	})
	if err != nil {
		logger.Fatal("failed to create images handler", err)
	}

	// create router using handlers
	router := newRouter(logger, authHandler, imagesHandler)

	/////////////////////////////
	// Server config and start //
	/////////////////////////////
	server := http.Server{
		Handler: router,
		Addr:    ":8080",
	}

	go func() {
		logger.Info("starting server...")
		// TODO: convert to ListenAndServeTLS
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen and serve server error", err)
		}
		logger.Info("stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server shutdown error", err)
	}
	logger.Info("server shutdown complete")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("                _        _               _____                                            _____ _____ \n     /\\        | |      | |             |_   _|                                     /\\   |  __ \\_   _|\n    /  \\  _   _| |_ ___ | | ___   __ _    | |  _ __ ___   __ _  __ _  ___  ___     /  \\  | |__) || |  \n   / /\\ \\| | | | __/ _ \\| |/ _ \\ / _` |   | | | '_ ` _ \\ / _` |/ _` |/ _ \\/ __|   / /\\ \\ |  ___/ | |  \n  / ____ \\ |_| | || (_) | | (_) | (_| |  _| |_| | | | | | (_| | (_| |  __/\\__ \\  / ____ \\| |    _| |_ \n /_/    \\_\\__,_|\\__\\___/|_|\\___/ \\__, | |_____|_| |_| |_|\\__,_|\\__, |\\___||___/ /_/    \\_\\_|   |_____|\n                                  __/ |                         __/ |                                 \n                                 |___/                         |___/                                  \n est. 2025 - Author: keola-dunn\n"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

func robotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User-agent: *\nDisallow: /"))
}

func newRouter(logger *logger.Logger, authHandler *jwt.AuthHandler, imageHandler *images.ImagesHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.RequestLogger)

	router.Get("/", home)
	router.Get("/health", healthCheck)

	router.Get("/robots.txt", robotsTxt)

	// debug with pprof
	// debug/pprof/
	// TODO: add auth to this endpoint to prevent public access
	router.Mount("/debug", middleware.Profiler())

	router.Route("/v1", func(router chi.Router) {
		router.Route("/images", func(router chi.Router) {
			// GET image
			router.Get("/{id}", imageHandler.GetImage)

			// POST image(s)
			router.With(authHandler.RequireTokenAuthentication).Post("/", imageHandler.PostImage)
		})
	})

	return router
}
