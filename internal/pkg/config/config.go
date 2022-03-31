package config

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/klaital/stock-portfolio-api/datalayer"
	"github.com/klaital/stock-portfolio-api/datalayer/mysql"
	"github.com/klaital/stock-portfolio-api/stockfetcher"
	"github.com/klaital/stock-portfolio-api/stockfetcher/finnhub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Config struct {
	// Meta
	Environment    string
	AuthBasePath   string
	StocksBasePath string
	Timeout        time.Duration

	// Auth
	CookieStoreSecret  string
	CookieMaxAge       int
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackUrl  string
	HashCost           int

	// Databases
	store            datalayer.StockStore
	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     int
	DatabaseName     string

	// External APIs
	stockFetcher  stockfetcher.StockFetcher
	FinnHubApiKey string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./run/")
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Debug("Failed to read config file")
	}

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "localreader")
	viper.SetDefault("DB_PASSWORD", "nopassword")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_NAME", "test")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_PRETTY", false)
	viper.SetDefault("ENV", "local")
	viper.SetDefault("FINNHUB_API_KEY", "")
	viper.SetDefault("HASH_COST", bcrypt.DefaultCost)
	viper.SetDefault("COOKIE_MAX_AGE", 86400*30) // 30 days
	viper.SetDefault("GOOGLE_CALLBACK_URL", "http://localhost:3000/auth/google/callback")
	viper.SetDefault("TIMEOUT", 2*time.Second)

	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("LOG_PRETTY")
	viper.BindEnv("REALM")
	viper.BindEnv("FINNHUB_API_KEY")
	viper.BindEnv("COOKIE_STORE_SECRET")
	viper.BindEnv("GOOGLE_CLIENT_ID")
	viper.BindEnv("GOOGLE_CLIENT_SECRET")
	viper.BindEnv("HASH_COST")
	viper.BindEnv("COOKIE_MAX_AGE")
	viper.BindEnv("GOOGLE_CALLBACK_URL")
	viper.BindEnv("TIMEOUT")

	logLevel, err := log.ParseLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		log.WithField("level", viper.GetString("LOG_LEVEL")).Error("Unable to parse log level")
		logLevel = log.DebugLevel
	}
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: viper.GetBool("LOG_PRETTY"),
	})
	log.SetLevel(logLevel)
	log.SetReportCaller(true)

	return &Config{
		Environment:        viper.GetString("ENV"),
		CookieStoreSecret:  viper.GetString("COOKIE_STORE_SECRET"),
		CookieMaxAge:       viper.GetInt("COOKIE_MAX_AGE"),
		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackUrl:  viper.GetString("GOOGLE_CALLBACK_URL"),
		FinnHubApiKey:      viper.GetString("FINNHUB_API_KEY"),
		HashCost:           viper.GetInt("HASH_COST"),
		DatabaseHost:       viper.GetString("DB_HOST"),
		DatabaseUser:       viper.GetString("DB_USER"),
		DatabasePassword:   viper.GetString("DB_PASSWORD"),
		DatabasePort:       viper.GetInt("DB_PORT"),
		DatabaseName:       viper.GetString("DB_NAME"),
		Timeout:            viper.GetDuration("TIMEOUT"),
	}
}

// GetDatalayer connects to the database and returns a valid StockStore or a connection error.
//The connection pool is cached in the config object.
func (c *Config) GetDatalayer() (datalayer.StockStore, error) {
	if c.store != nil {
		return c.store, nil
	}

	var err error
	c.store, err = mysql.New(context.Background(), c.DatabaseHost, c.DatabaseUser, c.DatabasePassword, c.DatabaseName, c.DatabasePort, c.HashCost)
	return c.store, err
}

// GetStockFetcher generates the API client for fetching stock prices.
//The client is cached in the config object.
func (c *Config) GetStockFetcher() stockfetcher.StockFetcher {
	if c.stockFetcher != nil {
		return c.stockFetcher
	}

	c.stockFetcher = finnhub.New(c.FinnHubApiKey)
	return c.stockFetcher
}

func (c *Config) IsProd() bool {
	return c.Environment == "prod"
}

// Validate ensures that all required env vars are set. Does not validate that
//the values are valid, such as whether the DB password will work.
func (c *Config) Validate() error {
	if len(c.CookieStoreSecret) == 0 {
		return fmt.Errorf("missing COOKIE_STORE_SECRET")
	}
	if len(c.GoogleClientID) == 0 {
		return fmt.Errorf("missing GOOGLE_CLIENT_ID")
	}
	if len(c.GoogleClientSecret) == 0 {
		return fmt.Errorf("missing GOOGLE_CLIENT_SECRET")
	}

	if len(c.DatabaseHost) == 0 {
		return fmt.Errorf("missing DB_HOST")
	}
	if len(c.DatabaseUser) == 0 {
		return fmt.Errorf("missing DB_USER")
	}
	if len(c.DatabasePassword) == 0 {
		return fmt.Errorf("missing DB_PASSWORD")
	}
	return nil
}

func (c *Config) GetCookieStore() *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(c.CookieStoreSecret))
	store.MaxAge(c.CookieMaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	return store
}
