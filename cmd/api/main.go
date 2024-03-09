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
	db struct{
		dsn string
	}
}

type application struct{
	config config
	logger *log.Logger
}

var _ = godotenv.Load(".env")

func main(){
	var cfg config 

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development|staging|production)")
	flag.Parse()

	cfg.db.dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil{
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("database connection pool establisted")

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
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(cfg config)(*sql.DB, error){
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil{
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil{
		return nil, err
	}

	return db, nil
}
