package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct{
	port int
	env string
}

type application struct{
	config config
	logger *log.Logger
}

var _ = godotenv.Load(".env")

var (
	connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	      os.Getenv("host"),
	      os.Getenv("port"),
	      os.Getenv("user"),
	      os.Getenv("password"),
	      os.Getenv("dbname"),
		)
)

func main(){
	var cfg config 

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}


	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %d", cfg.env, cfg.port)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
