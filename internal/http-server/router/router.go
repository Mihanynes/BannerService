package router

import (
	"avito-banner-service/internal/http-server/handlers/api/banner"
	user_banner "avito-banner-service/internal/http-server/handlers/api/user-banner"
	"avito-banner-service/internal/http-server/handlers/auth/token"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

//func Routes(redisClient *redis.Redis, db *postgres.Storage) *chi.Mux {
//	router := chi.NewRouter()
//
//	router.Use(middleware.RequestID)
//	router.Use(middleware.Logger)
//	router.Use(middleware.Recoverer)
//	router.Use(token.NewTokenHandler())
//
//	router.Get("/ping", user_banner.Ping())
//	router.Get("/user_banner", user_banner.GetBannerById(redisClient, db))
//	router.Get("/banner", banner.GetBannersFiltered(redisClient, db))
//	router.Post("/banner", banner.CreateBanner(redisClient, db))
//	router.Patch("/banner/{id}", banner.UpdateBanner(redisClient, db))
//	router.Delete("/banner/{id}", banner.DeleteBanner(redisClient, db))
//
//	return router
//}

func Routes(redisClient *redis.Redis, db *postgres.Storage) *fiber.App {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(token.NewTokenHandler())

	app.Get("/ping", user_banner.Ping)
	app.Get("/user_banner", user_banner.GetBannerById(redisClient, db))
	app.Get("/banner", banner.GetBannersFiltered(redisClient, db))
	app.Post("/banner", banner.CreateBanner(redisClient, db))
	app.Patch("/banner/:id", banner.UpdateBanner(redisClient, db))
	app.Delete("/banner/:id", banner.DeleteBanner(redisClient, db))

	return app
}
