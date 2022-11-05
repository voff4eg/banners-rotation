package main

import (
	"banners-rotation/internal/config"
	"banners-rotation/internal/rmq"
	"banners-rotation/internal/server"
	"banners-rotation/internal/storage"
	"context"
	"flag"
	"github.com/jackc/pgx/v4"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config/config.yaml", "Путь до файла настроек")
}

func main() {
	flag.Parse()

	cfg := config.NewConfig(configFile)
	ctx := context.Background()

	svr := server.NewServer(*cfg)

	rabbit, err := rmq.NewRabbit(ctx, cfg.Rabbit.Dsn, cfg.Rabbit.Exchange, cfg.Rabbit.Queue, cfg.Rabbit.Tag)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgx.Connect(ctx, cfg.Database.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	strg := storage.NewStorage(ctx, db)
	defer strg.CloseDb()

	err = svr.Run(server.NewRouters(strg, rabbit))
	if err != nil {
		log.Fatal(err)
	}
}
