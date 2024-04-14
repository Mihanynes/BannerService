package repositories

import (
	"banner-service/internal/models"
)

type Database interface {
	GetUserBannerByTagIdAndFeatureId(tagId int, featureId int) (models.UserBanner, error)
	GetBannersFilteredByFeatureOrTagId(tagVal models.NilInt, featureVal models.NilInt, limit int, offset int) ([]models.UserBannerFilteredResponse, error)
	CreateUserBanner(createBanner models.CreateBannerRequest) (models.UserBanner, error)
	GetNextUserBannerId() int
	GetBannerById(id int) (models.UserBanner, error)
	UpdateUserBanner(id int, request models.CreateBannerRequest, banner models.UserBanner) (models.UserBanner, error)
	DeleteBannerById(id int) error
}

type Cache interface {
	PutBanner(tagId int, featureId int, banner models.UserBanner)
	PutEmptyBanner(tagId int, featureId int)
	IsEmptyBanner(tagId int, featureId int) bool
	GetBannerById(tagId int, featureId int, banner interface{}) bool
	GetBannerGroup(tagVal models.NilInt, featureVal models.NilInt, limit int, offset int) ([]models.UserBannerFilteredResponse, error)
	PutBannerGroup(tagVal models.NilInt, featureVal models.NilInt, banners []models.UserBannerFilteredResponse, limit int, offset int)
}
