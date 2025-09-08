// Filename: cmd/api/main.go

package main

import (
	"flag"
	"log/slog"
	"os"
	"time"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type configuration struct {
	port int
	env  string
	vrs  string
	db   struct {
		dsn string
	}
}

// dependency injection
type application struct {
	config configuration
	logger *slog.Logger
	db *sql.DB
}



func main() {

	//Initialize configuration
	cfg := loadConfig()
	//Initialize logger
	logger := setupLogger()

	//set up DB connection
	db, err := openDB(cfg)
	if err != nil {
		logger.Error("cannot connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	//Initialize applicatioin with dependencies
	app := &application{
		config: cfg,
		logger: logger,
		db: db,
	}

	// Start the application server
	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func loadConfig() configuration {
	var cfg configuration

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.vrs, "version", "1.0.0", "Application version")
	//read in the dsn
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://quotes:quotesadvweb@localhost/quotes", "PostgresSQL DSN")
	flag.Parse()

	return cfg
}

func setupLogger() *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))


	logger.Info("database connection pool established")

	return logger
	
}

func openDB(cfg configuration) (*sql.DB, error) {
	//open a connection pool
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Second)

	defer cancel()

	//test if the connection pool was created
	//try pinging it with a 5 second timeout
	err = db.PingContext(ctx)

	if err != nil {
		db.Close()
		return nil, err
	}

	//return the connection pool (sql.DB)
	return db, nil

}
