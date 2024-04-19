package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vladComan0/tasty-byte/internal/models"
)

type config struct {
	addr           string
	environment    string
	dsn            string
	debugEnabled   bool
	allowedOrigins []string
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	recipes  *models.RecipeModel
}

func main() {
	var config config
	flag.StringVar(&config.addr, "addr", ":4000", "HTTP endpoint the server should listen on.")
	flag.StringVar(&config.environment, "env", "development", "Environment the application is running in.")
	flag.BoolVar(&config.debugEnabled, "debug", false, "Enable debug mode.")
	allowedOrigins := flag.String("origins", "http://192.168.100.20:4200", "Allowed origins for CORS, comma-separated.")
	flag.Parse()

	config.dsn = fmt.Sprintf("tastybyte_user:%s@tcp(db:3306)/tastybyte?parseTime=true", os.Getenv("MYSQL_PASSWORD"))

	config.allowedOrigins = strings.Split(*allowedOrigins, ",")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(config.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// dependency injection
	app := &application{
		config:   config,
		infoLog:  infoLog,
		errorLog: errorLog,
		recipes:  &models.RecipeModel{DB: db},
	}

	server := &http.Server{
		Addr:         config.addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on port: %s", strings.Split(server.Addr, ":")[1])
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
