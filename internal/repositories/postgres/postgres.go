package postgres

import (
	"banner-service/internal/models"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "avito_db"
)

type Storage struct {
	Db *sql.DB
}

func New() (*Storage, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("ready to connect psql")

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	storage := &Storage{Db: db}

	if err := storage.initialize(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) initialize() error {
	// Run initialization queries for database setup
	initQuery := s.readInitFile()
	rows, err := s.Db.Query(initQuery)
	if err != nil {
		return err
	}

	for rows.Next() {
		var version string
		rows.Scan(&version)
		fmt.Println(version)
	}
	fmt.Println("success read init sql file")

	// Fill database with test data
	fillQuery := s.readFillFile()
	_, err = s.Db.Query(fillQuery)
	if err != nil {
		fmt.Println("error while fill test data. continue with that data or rebuild postgres container")
		//log.Fatal(err)
		//panic(err)
	}

	return nil
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func (s *Storage) readInitFile() string {
	file, err := os.Open("resources/database/init.sql")
	slog.Info("in read init", err)
	CheckError(err)
	b, err := io.ReadAll(file)
	CheckError(err)
	//fmt.Println(b)
	return string(b)
}

func (s *Storage) readFillFile() string {
	file, err := os.Open("resources/database/fill.sql")
	CheckError(err)
	b, err := io.ReadAll(file)
	CheckError(err)
	//fmt.Println(b)
	return string(b)
}

func (s *Storage) GetUserBannerByTagIdAndFeatureId(tagId int, featureId int) (models.UserBanner, error) {
	var banner models.UserBanner
	err := s.Db.QueryRow(`SELECT ub.id, ub.content, ub.is_active 
										FROM user_banners_tags join user_banners ub on ub.id = user_banners_tags.banner_id 
										where feature_id = $1 and tag_id = $2;`, featureId, tagId).Scan(&banner.Id, &banner.Content, &banner.IsActive)
	return banner, err
}

func (s *Storage) GetBannersFilteredByFeatureOrTagId(tagVal models.NilInt, featureVal models.NilInt, limit int, offset int) ([]models.UserBannerFilteredResponse, error) {
	rows, err := s.Db.Query(`SELECT ub.id, ub.content, ub.is_active, ub.feature_id, ub.created_at, ub.updated_at, tag_id
				FROM user_banners_tags join user_banners ub on ub.id = user_banners_tags.banner_id
				where (feature_id = $1 OR $1 IS NULL) and (tag_id = $2 OR $2 IS NULL) 
				ORDER BY created_at
				limit $3 
				offset $4;`, featureVal.GetValue(), tagVal.GetValue(), limit, offset)
	if err != nil {
		slog.Info("error while getting banners data", err)
		return []models.UserBannerFilteredResponse{}, err
	}

	banners := []models.UserBannerFilteredResponse{}
	bannerIdPositionMap := make(map[int]int)
	for rows.Next() {
		var banner models.UserBannerFilteredResponse
		var tagId int
		rows.Scan(&banner.BannerId, &banner.Content, &banner.IsActive, &banner.FeatureId, &banner.CreatedAt, &banner.UpdatedAt, &tagId)

		if bannerIdPositionMap[banner.BannerId] == 0 {
			banners = append(banners, banner)
			bannerIdPositionMap[banner.BannerId] = len(banners)
			banners[bannerIdPositionMap[banner.BannerId]-1].TagIds = []int{tagId}
		} else {
			banners[bannerIdPositionMap[banner.BannerId]-1].TagIds = append(banners[bannerIdPositionMap[banner.BannerId]-1].TagIds, tagId)
		}
	}
	//fmt.Println(bannerIdPositionMap)
	//fmt.Println(banners)
	return banners, nil
}

func (s *Storage) CreateUserBanner(createBanner models.CreateBannerRequest) (models.UserBanner, error) {

	var banner = models.UserBanner{
		Id:        s.GetNextUserBannerId(),
		Content:   createBanner.Content,
		IsActive:  createBanner.IsActive,
		FeatureId: createBanner.FeatureId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	slog.Info("new banner")
	fmt.Println(banner)

	ctx := context.Background()
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_banners (id, content, is_active, feature_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6);`, banner.Id, banner.Content, banner.IsActive, banner.FeatureId, banner.CreatedAt, banner.UpdatedAt)
	if err != nil {
		slog.Error("error while saving new banner ", err)
		tx.Rollback()
		return models.UserBanner{}, err
	}

	for _, tagId := range createBanner.TagIds {
		_, err := tx.ExecContext(ctx, `INSERT INTO user_banners_tags (banner_id, tag_id) values ($1, $2);`, banner.Id, tagId)
		if err != nil {
			slog.Error("error adding tag id", tagId, err)
			tx.Rollback()
			return models.UserBanner{}, err
		}
	}

	tx.Commit()

	return banner, nil
}

func (s *Storage) CreateUserBannerWithId(createBanner models.CreateBannerRequest, id int) (models.UserBanner, error) {

	var banner = models.UserBanner{
		Id:        id,
		Content:   createBanner.Content,
		IsActive:  createBanner.IsActive,
		FeatureId: createBanner.FeatureId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	slog.Info("new banner")
	fmt.Println(banner)

	ctx := context.Background()
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_banners (id, content, is_active, feature_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6);`, banner.Id, banner.Content, banner.IsActive, banner.FeatureId, banner.CreatedAt, banner.UpdatedAt)
	if err != nil {
		slog.Error("error while saving new banner ", err)
		tx.Rollback()
		return models.UserBanner{}, err
	}

	for _, tagId := range createBanner.TagIds {
		_, err := tx.ExecContext(ctx, `INSERT INTO user_banners_tags (banner_id, tag_id) values ($1, $2);`, banner.Id, tagId)
		if err != nil {
			slog.Error("error adding tag id", tagId, err)
			tx.Rollback()
			return models.UserBanner{}, err
		}
	}

	tx.Commit()

	return banner, nil
}

func (s *Storage) GetNextUserBannerId() int {
	var maxId int
	err := s.Db.QueryRow(`select MAX(id) FROM user_banners;`).Scan(&maxId)
	if err != nil {
		slog.Error("error while saving new banner ", err)
		return -1
	}
	return maxId + 1
}

func (s *Storage) GetBannerById(id int) (models.UserBanner, error) {
	var banner models.UserBanner
	err := s.Db.QueryRow(`SELECT ub.id, ub.content, ub.is_active, ub.feature_id, ub.created_at, ub.updated_at
										FROM  user_banners ub where id = $1`, id).Scan(&banner.Id, &banner.Content, &banner.IsActive, &banner.FeatureId, &banner.CreatedAt, &banner.UpdatedAt)
	return banner, err
}

func (s *Storage) UpdateUserBanner(id int, request models.CreateBannerRequest, banner models.UserBanner) (models.UserBanner, error) {
	ctx := context.Background()
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx, `UPDATE user_banners set content = $1, is_active = $2, feature_id = $3, created_at = $4, updated_at = $5 where id = $6;`, request.Content, request.IsActive, request.FeatureId, banner.CreatedAt, time.Now(), id)

	if err != nil {
		slog.Error("error while saving new banner ", err)
		tx.Rollback()
		return models.UserBanner{}, err
	}

	//Удалить если обновление только добавляет тэги, а не удаляет старые
	_, err = tx.ExecContext(ctx, `delete from user_banners_tags where banner_id = $1`, id)

	if err != nil {
		slog.Error("error while deleting old banner tags ", err)
		tx.Rollback()
		return models.UserBanner{}, err
	}
	//до сюда

	for _, tagId := range request.TagIds {
		_, err := tx.ExecContext(ctx, `INSERT INTO user_banners_tags (banner_id, tag_id) values ($1, $2);`, banner.Id, tagId)
		if err != nil {
			slog.Error("error adding tag id", tagId, err)
			tx.Rollback()
			return models.UserBanner{}, err
		}
	}
	tx.Commit()
	slog.Info("success update banner")

	return models.UserBanner{
		Id:        id,
		Content:   request.Content,
		FeatureId: request.FeatureId,
		IsActive:  request.IsActive,
	}, nil
}

func (s *Storage) DeleteBannerById(id int) error {
	ctx := context.Background()
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = tx.ExecContext(ctx, `delete from user_banners_tags where banner_id = $1`, id)
	if err != nil {
		slog.Error("error while deleting banner tags ", err)
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `delete from user_banners where id = $1`, id)
	if err != nil {
		slog.Error("error while deleting banner ", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	slog.Info("success delete banner")
	return nil
}
