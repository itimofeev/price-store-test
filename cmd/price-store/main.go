package main

import (
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/itimofeev/price-store-test/internal/downloader"
	"github.com/itimofeev/price-store-test/internal/handlers"
	"github.com/itimofeev/price-store-test/internal/mongo"
	"github.com/itimofeev/price-store-test/internal/pg"
	"github.com/itimofeev/price-store-test/internal/service"
	"github.com/itimofeev/price-store-test/internal/util"
)

const pgURL = "postgresql://postgres:password@db:5432/postgres?sslmode=disable"

func main() {
	dbURL := getEnvOrDefault("DB_URL", pgURL)

	log := util.NewLog()
	d := downloader.New()
	store := createStore(log, dbURL)
	log.WithField("url", dbURL).Debug("store created")
	srv := service.New(log, d, store)

	app := handlers.InitApp(srv)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}

func getEnvOrDefault(envName, defaultVal string) string {
	if val := os.Getenv(envName); val != "" {
		return val
	}
	return defaultVal
}

func createStore(log *logger.Logger, dbURL string) service.Store {
	if strings.HasPrefix(dbURL, "mongodb://") {
		return mongo.New(dbURL)
	}
	return pg.New(log, dbURL)
}
