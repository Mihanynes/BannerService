package main

import (
	"avito-banner-service/internal/http-server/router"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/redis"
	"fmt"
	"log"
	"log/slog"
	"os"
)

func main() {
	fmt.Println("Hello World!")

	db, err := postgres.New()
	if err != nil {
		slog.Error("failed to init storage")
		os.Exit(1)
	}

	redisClient, err := redis.New()
	if err != nil {
		slog.Error("failed to init redis")
		os.Exit(1)
	}

	app := router.Routes(redisClient, db)
	log.Fatal(app.Listen(":8080"))

}
