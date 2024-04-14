package main

import (
	"banner-service/internal/http-server/router"
	"banner-service/internal/models"
	"banner-service/internal/repositories/postgres/mockdb"
	"banner-service/internal/repositories/redis/mock-cache"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	app := router.Routes(nil, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := 200
	if resp.StatusCode != statusCode {
		t.Errorf("TestPing() test returned an unexpected result: got %v want %v", resp.StatusCode, statusCode)
	}
}

func TestTokenAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := router.Routes(nil, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := 401
	if resp.StatusCode != statusCode {
		t.Errorf("TestTokenAuth() test returned an unexpected result: got %v want %v", resp.StatusCode, statusCode)
	}
}

func TestHandleGetUserBanner(t *testing.T) {
	db, err := mockdb.New()
	if err != nil {
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/user_banner?tag_id=1&feature_id=1&use_last_revision=false", nil)
	req.Header.Set("token", "user_token")
	if err != nil {
		t.Fatal(err)
	}

	app := router.Routes(redisClient, db)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := 200
	if resp.StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", resp.StatusCode, statusCode)
	}
}

func TestCreateBanner(t *testing.T) {
	db, err := mockdb.New()
	if err != nil {
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		t.Fatal(err)
	}

	ids := []int{1, 4}
	featureId := 2
	content := `{"data": "some"}`
	isActive := true

	banner := models.CreateBannerRequest{
		TagIds:    ids,
		FeatureId: featureId,
		Content:   json.RawMessage(content),
		IsActive:  isActive,
	}

	body, err := json.Marshal(banner)

	req, err := http.NewRequest("POST", "/banner", strings.NewReader(string(body)))
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(req)
	app := router.Routes(redisClient, db)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := 201
	if resp.StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", resp.StatusCode, statusCode)
	}

	err = db.DeleteBannerById(11)
	if err != nil {
		t.Fatal("created banner not deleted((")
	}
}

func TestPatchBanner(t *testing.T) {
	db, err := mockdb.New()
	if err != nil {
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		t.Fatal(err)
	}

	ids := []int{1, 4}
	featureId := 2
	content := `{"data": "some EDITED"}`
	isActive := true

	banner := models.CreateBannerRequest{
		TagIds:    ids,
		FeatureId: featureId,
		Content:   json.RawMessage(content),
		IsActive:  isActive,
	}

	body, err := json.Marshal(banner)

	req, err := http.NewRequest("PATCH", "/banner/4", strings.NewReader(string(body)))
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	app := router.Routes(redisClient, db)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := 200
	if resp.StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", resp.StatusCode, statusCode)
	}
}

func TestDeleteBanner(t *testing.T) {
	db, err := mockdb.New()
	if err != nil {
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/banner/4", nil)
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	app := router.Routes(redisClient, db)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := 204
	if resp.StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", resp.StatusCode, statusCode)
	}

	ids := []int{1, 4}
	featureId := 2
	content := `{"data": "some"}`
	isActive := true

	banner := models.CreateBannerRequest{
		TagIds:    ids,
		FeatureId: featureId,
		Content:   json.RawMessage(content),
		IsActive:  isActive,
	}

	db.CreateUserBannerWithId(banner, 4)
	if err != nil {
		t.Fatal("create banner after delete error")
	}
}
