package main

import (
	"crypto/tls"
	"database/sql"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vladComan0/tasty-byte/internal/models"
)

type config struct {
	Addr           string   `mapstructure:"addr"`
	Environment    string   `mapstructure:"environment"`
	DSN            string   `mapstructure:"dsn"`
	DebugEnabled   bool     `mapstructure:"debug_enabled"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	recipes  models.RecipeModelInterface
}

func main() {
	var config config

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	getConfig(errorLog, &config)

	db, err := openDB(config.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	ingredientModel := &models.IngredientModel{
		DB: db,
	}

	recipeIngredientModel := &models.RecipeIngredientModel{
		DB: db,
	}

	tagModel := &models.TagModel{
		DB: db,
	}

	recipeTagModel := &models.RecipeTagModel{
		DB: db,
	}

	// dependency injection
	app := &application{
		config:   config,
		infoLog:  infoLog,
		errorLog: errorLog,
		recipes: &models.RecipeModel{
			DB:                    db,
			IngredientModel:       ingredientModel,
			RecipeIngredientModel: recipeIngredientModel,
			TagModel:              tagModel,
			RecipeTagModel:        recipeTagModel,
		},
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	server := &http.Server{
		Addr:         config.Addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}

	infoLog.Printf("Starting server on port: %s", strings.Split(server.Addr, ":")[1])
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func getConfig(errorLog *log.Logger, config *config) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		errorLog.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		errorLog.Fatalf("Unable to decode into struct, %v", err)
	}
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
