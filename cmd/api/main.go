// Filename: cmd/api/main.go

package main

import (
	"flag"
	"log/slog"
	"os"
	"time"
	"context"
	"database/sql"
	"strings"
	_ "github.com/lib/pq"
	"github.com/kelseyaban/qod/internal/data"
)

type configuration struct {
	port int
	env  string
	vrs  string
	db   struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
	limiter struct {
		rps float64 //request per second
		burst int  //initial requests posssible 
		enabled bool  //enable or disable rate limiter 
	}
}

// dependency injection
type application struct {
	config configuration
	logger *slog.Logger
	quoteModel data.QuoteModel
	// quotes *data.QuoteModel
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
		quoteModel: data.QuoteModel{DB: db},

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

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2,"Rate Limiter maximum requests per second")

    flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5,"Rate Limiter maximum burst")

    flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true,"Enable rate limiter")


	// We will build a custom command-line flag.  This flag will allow us to access space-separated origins. 
	//We will then put those origins in our slice. Again notsomething we can do with the flag functions that we have seen so far. 
	// strings.Fields() splits string (origins) on spaces
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
			
		cfg.cors.trustedOrigins = strings.Fields(val)
			return nil
		})

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
