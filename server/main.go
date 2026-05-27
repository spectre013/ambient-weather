package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	db     *sql.DB
	logger = logrus.New()
)

func init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	// Load .env in development. In production we expect the runtime to
	// inject environment variables.
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logger.WithError(err).Fatal("error loading .env file")
		}
	}

	// Parse and validate all configuration up front.
	cfg, err := LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("invalid configuration")
	}
	config = cfg

	if cfg.LogLevel == "Debug" {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.WithField("level", logger.GetLevel()).Info("log level configured")

	db, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		logger.WithError(err).Fatal("failed to open postgres connection")
	}
	defer db.Close()

	// Sensible pool defaults; tune based on load.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("failed to ping database")
	}

	// Set up the websocket hub and broadcast loop.
	hub := newHub()
	go hub.run()
	go broadcast(hub)

	r := mux.NewRouter()
	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("websocket connection")
		serveWs(hub, w, r)
	})

	r.HandleFunc("/api/current", loggingMiddleware(current))
	r.HandleFunc("/api/app", loggingMiddleware(app))
	r.HandleFunc("/api/history", loggingMiddleware(historyHandler))
	r.HandleFunc("/api/chart/{sensor}/{time}", loggingMiddleware(chart))
	r.HandleFunc("/api/alltime/{calc}/{type}", loggingMiddleware(alltime))
	r.HandleFunc("/api/forecast", loggingMiddleware(getForecastHandler))
	r.HandleFunc("/api/climate", loggingMiddleware(getClimate))
	r.HandleFunc("/api/firstfreeze", loggingMiddleware(getFirstFreeze))
	r.HandleFunc("/api/about", loggingMiddleware(about))
	r.HandleFunc("/", loggingMiddleware(home))

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + cfg.Port,
		// Added IdleTimeout and ReadHeaderTimeout to protect
		// against slow-loris-style resource exhaustion.
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown: catch SIGINT/SIGTERM, stop accepting
	// new connections, and give in-flight requests up to 30s to finish.
	go func() {
		logger.WithField("port", cfg.Port).Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("server failed")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("server shutdown failed")
	}
	logger.Info("server stopped")
}
