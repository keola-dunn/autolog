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
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/keola-dunn/autolog/cmd/autolog-api/handlers/auth"
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/user"
)

var environmentConfig struct {
	DBUser     string `envconfig:"DB_USER"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBHost     string `envconfig:"DB_HOST"`
	DBPort     int64  `envconfig:"DB_PORT"`

	JWTSecret string `envconfig:"JWT_SECRET"`
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
		ConnectionConfig: postgres.ConnectionConfig{User: environmentConfig.DBUser,
			Password: environmentConfig.DBPassword,
			Host:     environmentConfig.DBHost,
			Port:     environmentConfig.DBPort,
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

	calendarSvc := calendar.NewService()

	///////////////////////
	// Service Creations //
	///////////////////////

	userSvc := user.NewService(user.ServiceConfig{
		DB:              db,
		RandomGenerator: randomSvc,
	})

	///////////////////////////
	// API Handler Creations //
	///////////////////////////

	authHandler := auth.NewAuthHandler(auth.AuthHandlerConfig{
		JWTSecret:              environmentConfig.JWTSecret,
		JWTIssuer:              "",
		JWTExpiryLengthMinutes: 1,
		CalendarService:        calendarSvc,
		RandomGenerator:        randomSvc,
		Logger:                 logger,
		UserService:            userSvc,
	})

	// create router using handlers
	router := newRouter(logger, authHandler)

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
	w.Write([]byte("Autolog Maintenance Log API!"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

func newRouter(logger *logger.Logger, authHandler *auth.AuthHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.RequestLogger)

	router.Get("/", home)
	router.Get("/health", healthCheck)

	router.Route("/v1", func(router chi.Router) {
		router.Route("/auth", func(router chi.Router) {
			router.Post("/login", authHandler.Login)
			router.Post("/signup", authHandler.SignUp)
			router.Get("/security-questions", authHandler.GetSecurityQuestions)
		})
	})

	return router
}
