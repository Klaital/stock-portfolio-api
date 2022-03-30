package main

import (
	"context"
	"flag"
	"github.com/klaital/stock-portfolio-api/datalayer/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type config struct {
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     int
}

func main() {

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./run/")
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Fatal("Failed to read config file")
	}

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "localreader")
	viper.SetDefault("DB_PASSWORD", "nopassword")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_NAME", "test")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_PRETTY", false)
	viper.SetDefault("REALM", "local")
	viper.SetDefault("NASDAQ_API_KEY", "")

	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("LOG_PRETTY")
	viper.BindEnv("REALM")
	viper.BindEnv("FINNHUB_API_KEY")

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

	datastore, err := mysql.New(context.Background(),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetInt("DB_PORT"),
		bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Fatal("Failed to init datastore")
	}

	switch os.Args[1] {
	case "add-stock":
		addPosition(datastore)
	case "list-positions":
		listPositionsCmd := flag.NewFlagSet("list-positions", flag.ExitOnError)
		listParams := struct {
			userId uint64
		}{}
		listPositionsCmd.Uint64Var(&listParams.userId, "user", 0, "User ID")
		listPositionsCmd.Parse(os.Args[2:])

		listPositions(datastore, listParams.userId)

	case "add-user":
		cmd := flag.NewFlagSet("add-user", flag.ExitOnError)
		params := struct {
			email    string
			password string
		}{}
		cmd.StringVar(&params.email, "email", "", "User email")
		cmd.StringVar(&params.password, "pass", "", "User password")
		cmd.Parse(os.Args[2:])

		addUser(datastore, params.email, params.password)
	}

}
