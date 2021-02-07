package main

import (
	"github.com/itimofeev/price-store-test/cmd/internal/downloader"
	"github.com/itimofeev/price-store-test/cmd/internal/handlers"
	"github.com/itimofeev/price-store-test/cmd/internal/pg"
	"github.com/itimofeev/price-store-test/cmd/internal/service"
	"github.com/itimofeev/price-store-test/cmd/internal/util"
)

const pgURL = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func main() {
	log := util.NewLog()
	d := downloader.New()
	store := pg.New(log, pgURL)
	srv := service.New(log, d, store)

	app := handlers.InitApp(srv)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
