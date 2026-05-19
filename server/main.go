package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
var db *sql.DB
var logger = logrus.New()

func init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	var err error
	if os.Getenv("GO_ENV") != "production" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	logLevel := logrus.InfoLevel
	if os.Getenv("LOGLEVEL") == "Debug" {
		logLevel = logrus.DebugLevel
	}
	logger.Info("Setting Debug Level to ", logLevel)
	logger.SetLevel(logLevel)

	dburi := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
	db, err = sql.Open("postgres", dburi)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Setup Web Sockets
	hub := newHub()
	go hub.run()
	go broadcast(hub)

	r := mux.NewRouter()
	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Websocket connection")
		serveWs(hub, w, r)
	})
	//API
	r.HandleFunc("/api/current", loggingMiddleware(current))
	r.HandleFunc("/api/app", loggingMiddleware(app))
	r.HandleFunc("/api/chart/{sensor}/{time}", loggingMiddleware(chart))
	r.HandleFunc("/api/alltime/{calc}/{type}", loggingMiddleware(alltime))
	r.HandleFunc("/api/forecast", loggingMiddleware(getForecastHandler))
	r.HandleFunc("/api/climate", loggingMiddleware(getClimate))
	r.HandleFunc("/api/firstfreeze", loggingMiddleware(getFirstFreeze))
	r.HandleFunc("/api/about", loggingMiddleware(about))
	//Index
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))
	//r.Handle("/", http.FileServer(http.Dir("./public")))

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + os.Getenv("PORT"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server on " + os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logger.Infof(
			"%s\t%s\t%s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
		)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
