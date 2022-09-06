package services

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/docs"
	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/pkg/sentry"
	"github.com/unanoc/blockchain-indexer/pkg/viper"
)

const defaultConfigPath = "config.yml"

func Setup() {
	InitConfig()
	InitLogging()
	InitSentry()
	InitDatabase()
	InitSwaggerInfo()
}

func InitConfig() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultConfigPath
	}

	viper.Load(path, &config.Default)
}

func InitLogging() {
	logLevel, err := log.ParseLevel(config.Default.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Log level parsing error")
	}

	log.SetLevel(logLevel)
}

func InitSentry() {
	if err := sentry.SetupSentry(
		config.Default.Sentry.DSN,
		sentry.WithSampleRate(config.Default.Sentry.SampleRate),
	); err != nil {
		log.WithError(err).Fatal("Sentry init error")
	}
}

func InitDatabase() {
	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	if err := postgres.Setup(db); err != nil {
		log.WithError(err).Fatal("Database setup error")
	}
}

func InitSwaggerInfo() {
	docs.SwaggerInfo.Host = config.Default.Swagger.Hostname
}
