package banner

import (
	"avito-banner-service/internal/models"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/redis"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strconv"
)

func GetBannersFiltered(redisClient *redis.Redis, db *postgres.Storage) fiber.Handler {
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

		data, err := redis.GetBannerGroup(*redisClient, tagVal, featureVal, limit, offset)
		if err == nil {
			return c.JSON(data)
		}

		banners, err := postgres.GetBannersFilteredByFeatureOrTagId(db, tagVal, featureVal, limit, offset)
		if err != nil {
			slog.Error("error while getting filtered banners")
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		redis.PutBannerGroup(*redisClient, tagVal, featureVal, banners, limit, offset)
		return c.JSON(banners)
	}
}

//func CreateBanner(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if r.Header.Get("token") != "admin_token" {
//			resp.Send403Error(w, r)
//			return
//		}
//		var createBanner models.CreateBannerRequest
//		decoder := json.NewDecoder(r.Body)
//		decoder.DisallowUnknownFields()
//		err := decoder.Decode(&createBanner)
//		validate := validator.New()
//		err2 := validate.Struct(createBanner)
//		if err != nil || err2 != nil {
//			slog.Error("error while parsing request body")
//			resp.Send400Error(w, r)
//			return
//		}
//		slog.Info("continue updating ", err, err2)
//		banner, err := postgres.CreateUserBanner(db, createBanner)
//		if err != nil {
//			slog.Error("error while saving new banner")
//			resp.Send500Error(w, r)
//			return
//		}
//
//		slog.Info("success saved new banner " + string(banner.Id))
//		for _, tagId := range createBanner.TagIds {
//			redis.PutBanner(*redisClient, tagId, banner.FeatureId, banner)
//		}
//		//redis.PutBanner(*redisClient, )
//		w.WriteHeader(http.StatusCreated)
//		render.JSON(w, r, banner.Id)
//	}
//}

func CreateBanner(redisClient *redis.Redis, db *postgres.Storage) fiber.Handler {
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

		banner, err := postgres.CreateUserBanner(db, createBanner)
		if err != nil {
			slog.Error("error while saving new banner")
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving new banner")
		}

		slog.Info("successfully saved new banner " + string(banner.Id))
		for _, tagID := range createBanner.TagIds {
			redis.PutBanner(*redisClient, tagID, banner.FeatureId, banner)
		}

		return c.Status(fiber.StatusCreated).JSON(banner.Id)
	}
}

func UpdateBanner(redisClient *redis.Redis, db *postgres.Storage) fiber.Handler {
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
			slog.Error("error while parsing request id", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request ID")
		}

		banner, err := postgres.GetBannerById(db, id)
		if err != nil {
			slog.Error("error while getting banner by id ", id, err)
			return c.Status(fiber.StatusNotFound).SendString("Banner not found")
		}

		updatedBanner, err := postgres.UpdateUserBanner(db, id, bannerRequest, banner)
		if err != nil {
			slog.Error("error while updating banner", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating banner")
		}

		for _, tagID := range bannerRequest.TagIds {
			redis.PutBanner(*redisClient, tagID, banner.FeatureId, updatedBanner)
		}

		return c.Status(fiber.StatusOK).SendString("Banner updated successfully")
	}
}

func DeleteBanner(redisClient *redis.Redis, db *postgres.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("token") != "admin_token" {
			return c.Status(fiber.StatusForbidden).SendString("Unauthorized access")
		}

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request ID")
		}

		_, err = postgres.GetBannerById(db, id)
		if err != nil {
			slog.Error("error while getting banner by id ", id, err)
			return c.Status(fiber.StatusNotFound).SendString("Banner not found")
		}

		err = postgres.DeleteBannerById(db, id)
		if err != nil {
			slog.Error("error while deleting banner", id, err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting banner")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
