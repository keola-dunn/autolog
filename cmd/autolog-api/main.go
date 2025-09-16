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
	"github.com/keola-dunn/autolog/cmd/autolog-api/handlers/auth"
	"github.com/keola-dunn/autolog/cmd/autolog-api/handlers/cars"
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/jwt"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/car"
	"github.com/keola-dunn/autolog/internal/service/user"
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

	calendarSvc := calendar.NewService()

	jwtVerifier, err := jwt.NewTokenVerifier(context.Background(), jwt.TokenVerifierConfig{
		JWKSUrl: environmentConfig.JWKSUrl,
	})
	if err != nil {
		logger.Fatal("failed to create new jwt verifier", err)
	}

	///////////////////////
	// Service Creations //
	///////////////////////

	userSvc := user.NewService(user.ServiceConfig{
		DB:              db,
		RandomGenerator: randomSvc,
	})

	carSvc := car.NewService(car.ServiceConfig{
		DB:              db,
		RandomGenerator: randomSvc,
	})

	///////////////////////////
	// API Handler Creations //
	///////////////////////////

	authHandler, err := auth.NewAuthHandler(auth.AuthHandlerConfig{
		CalendarService: calendarSvc,
		RandomGenerator: randomSvc,
		Logger:          logger,
		UserService:     userSvc,
		TokenVerifier:   jwtVerifier,
	})
	if err != nil {
		logger.Fatal("failed to create auth handler", err)
	}

	carsHandler, err := cars.NewCarsHandler(cars.CarsHandlerConfig{
		CalendarService: calendarSvc,
		RandomGenerator: randomSvc,
		Logger:          logger,

		UserService:   userSvc,
		CarService:    carSvc,
		TokenVerifier: jwtVerifier,
	})
	if err != nil {
		logger.Fatal("failed to create cars handler", err)
	}

	// create router using handlers
	router := newRouter(logger, authHandler, carsHandler)

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
	w.Write([]byte("                _        _                        _____ _____ \n     /\\        | |      | |                 /\\   |  __ \\_   _|\n    /  \\  _   _| |_ ___ | | ___   __ _     /  \\  | |__) || |  \n   / /\\ \\| | | | __/ _ \\| |/ _ \\ / _` |   / /\\ \\ |  ___/ | |  \n  / ____ \\ |_| | || (_) | | (_) | (_| |  / ____ \\| |    _| |_ \n /_/    \\_\\__,_|\\__\\___/|_|\\___/ \\__, | /_/    \\_\\_|   |_____|\n                                  __/ |                       \n                                 |___/                        \n est. 2025 - Author: keola-dunn\n"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

func robotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User-agent: *\nDisallow: /"))
}

func newRouter(logger *logger.Logger, authHandler *auth.AuthHandler, carsHandler *cars.CarsHandler) *chi.Mux {
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
		router.Route("/users", func(router chi.Router) {
			// GET user details
			// authenticated only
			router.With(authHandler.RequireTokenAuthentication).Get("/", nil)
		})

		router.Route("/shops", func(router chi.Router) {
			// GET search for shops
			// public
			// ex. signing up for the right shop
			router.Get("/search", nil)
		})

		router.Route("/cars", func(router chi.Router) {
			// GET user's cars
			// authenticated only
			router.With(authHandler.RequireTokenAuthentication).Get("/", nil)

			// GET search for car
			// search by vin, ID, plate, etc.
			// public or authenticated
			router.Get("/lookup", nil)

			// PUT car if acquired
			// authenticated only
			router.With(authHandler.RequireTokenAuthentication).Put("/", carsHandler.CreateCar)

			router.Route("/{carId}", func(router chi.Router) {

				// GET car details and logs
				// public or authenticated
				router.Get("/", nil)

				// POST car update (if sold, etc.)
				// authenticated only
				router.With(authHandler.RequireTokenAuthentication).Post("/", nil)

				// POST maintence log
				// authenticated only
				router.With(authHandler.RequireTokenAuthentication).Post("/maintenance-log", nil)
			})

		})
	})

	return router
}
