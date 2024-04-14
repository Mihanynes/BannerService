package router

import (
	"banner-service/internal/http-server/handlers/api"
	"banner-service/internal/http-server/handlers/auth/token"
	"banner-service/internal/repositories/postgres"
	"banner-service/internal/repositories/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Routes(redisClient *redis.Redis, db *postgres.Storage) *fiber.App {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(token.NewTokenHandler())

	app.Get("/ping", api.Ping)
	app.Get("/user_banner", api.GetBannerById(redisClient, db))
	app.Get("/banner", api.GetBannersFiltered(redisClient, db))
	app.Post("/banner", api.CreateBanner(redisClient, db))
	app.Patch("/banner/:id", api.UpdateBanner(redisClient, db))
	app.Delete("/banner/:id", api.DeleteBanner(redisClient, db))

	return app
}
