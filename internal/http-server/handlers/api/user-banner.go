// user_banner/user_banner.go

package api

import (
	"banner-service/internal/models"
	"banner-service/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strconv"
)

type UserBannerRequest struct {
	TagId          int    `json:"tag_id" validate:"required"`
	FeatureId      int    `json:"feature_id" validate:"required"`
	UseLastVersion bool   `json:"use_last_version" default:"false"`
	Token          string `json:"token" default:"user_token"`
}

func GetBannerById(cache repositories.Cache, db repositories.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tagID, err := strconv.Atoi(c.Query("tag_id"))
		if err != nil {
			slog.Error("error parsing param tag_id")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid or missing 'tag_id' parameter")
		}
		featureID, err := strconv.Atoi(c.Query("feature_id"))
		if err != nil {
			slog.Error("error parsing param feature_id")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid or missing 'feature_id' parameter")
		}
		useLastRevision := c.Query("use_last_revision") == "true"

		if cache.IsEmptyBanner(tagID, featureID) {
			slog.Error("no banner info from cache")
			return c.Status(fiber.StatusNotFound).SendString("Banner not found in cache")
		}

		var banner models.UserBanner
		isCached := cache.GetBannerById(tagID, featureID, &banner)
		if useLastRevision || !isCached {
			banner, err = db.GetUserBannerByTagIdAndFeatureId(tagID, featureID)
			if err != nil {
				slog.Error("error while getting banner from database: ", err)
				cache.PutEmptyBanner(tagID, featureID)
				return c.Status(fiber.StatusNotFound).SendString("Banner not found in database")
			}
		}

		cache.PutBanner(tagID, featureID, banner)

		if !banner.IsActive && c.Get("token") != "admin_token" {
			return c.Status(fiber.StatusForbidden).SendString("Banner is not active or unauthorized access")
		}

		return c.JSON(banner.Content)
	}
}

func Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
