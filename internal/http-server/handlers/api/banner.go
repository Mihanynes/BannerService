// banner/banner.go

package api

import (
	"banner-service/internal/models"
	"banner-service/internal/repositories"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strconv"
)

func GetBannersFiltered(cache repositories.Cache, db repositories.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("token") != "admin_token" {
			return c.Status(fiber.StatusForbidden).SendString("Unauthorized access")
		}

		tagVal := models.NilInt{}
		tagID, err := strconv.Atoi(c.Query("tag_id"))
		if err != nil {
			tagVal.Null = true
		} else {
			tagVal.Null = false
			tagVal.Value = tagID
		}

		featureVal := models.NilInt{}
		featureID, err := strconv.Atoi(c.Query("feature_id"))
		if err != nil {
			featureVal.Null = true
		} else {
			featureVal.Null = false
			featureVal.Value = featureID
		}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			offset = 0
		}

		fmt.Println("limit", limit, "offset", offset)

		data, err := cache.GetBannerGroup(tagVal, featureVal, limit, offset)
		if err == nil {
			return c.JSON(data)
		}

		banners, err := db.GetBannersFilteredByFeatureOrTagId(tagVal, featureVal, limit, offset)
		if err != nil {
			slog.Error("error while getting filtered banners")
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		cache.PutBannerGroup(tagVal, featureVal, banners, limit, offset)
		return c.JSON(banners)
	}
}

func CreateBanner(cache repositories.Cache, db repositories.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("token") != "admin_token" {
			return c.Status(fiber.StatusForbidden).SendString("Unauthorized access")
		}

		var createBanner models.CreateBannerRequest
		if err := json.Unmarshal(c.Body(), &createBanner); err != nil {
			slog.Error("error while parsing request body")
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}

		validate := validator.New()
		if err := validate.Struct(createBanner); err != nil {
			slog.Error("error while validating request body")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		banner, err := db.CreateUserBanner(createBanner)
		if err != nil {
			slog.Error("error while saving new banner")
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving new banner")
		}

		slog.Info("successfully saved new banner " + string(banner.Id))
		for _, tagID := range createBanner.TagIds {
			cache.PutBanner(tagID, banner.FeatureId, banner)
		}

		return c.Status(fiber.StatusCreated).JSON(banner.Id)
	}
}

func UpdateBanner(cache repositories.Cache, db repositories.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("token") != "admin_token" {
			return c.Status(fiber.StatusForbidden).SendString("Unauthorized access")
		}

		var bannerRequest models.CreateBannerRequest
		if err := json.Unmarshal(c.Body(), &bannerRequest); err != nil {
			slog.Error("error while parsing request body")
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}

		validate := validator.New()
		if err := validate.Struct(bannerRequest); err != nil {
			slog.Error("error while validating request body")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			fmt.Println(id)
			slog.Error("error while parsing request id", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request ID")
		}

		banner, err := db.GetBannerById(id)
		if err != nil {
			slog.Error("error while getting banner by id ", id, err)
			return c.Status(fiber.StatusNotFound).SendString("Banner not found")
		}

		updatedBanner, err := db.UpdateUserBanner(id, bannerRequest, banner)
		if err != nil {
			slog.Error("error while updating banner", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating banner")
		}

		for _, tagID := range bannerRequest.TagIds {
			cache.PutBanner(tagID, banner.FeatureId, updatedBanner)
		}

		return c.Status(fiber.StatusOK).SendString("Banner updated successfully")
	}
}

func DeleteBanner(cache repositories.Cache, db repositories.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("token") != "admin_token" {
			return c.Status(fiber.StatusForbidden).SendString("Unauthorized access")
		}

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request ID")
		}

		_, err = db.GetBannerById(id)
		if err != nil {
			slog.Error("error while getting banner by id ", id, err)
			return c.Status(fiber.StatusNotFound).SendString("Banner not found")
		}

		err = db.DeleteBannerById(id)
		if err != nil {
			slog.Error("error while deleting banner", id, err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting banner")
		}

		return c.Status(fiber.StatusNoContent).SendString("Banner deleted")
	}
}
