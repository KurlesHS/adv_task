package postgres

import (
	"context"
	"fmt"
	"kurles/adv_task/configs"
	"kurles/adv_task/pkg/model"
	"math/rand"
	"reflect"
	"testing"
)

func TestAdvertRepo(t *testing.T) {
	cfg, err := configs.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	repo, err := New(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUserName, cfg.DBPassword)
	if err != nil {
		t.Error(err)
		return
	}
	_ = repo

	advs := make([]model.DetailedAdvert, 0)
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		adv := model.DetailedAdvert{
			Title:       fmt.Sprintf("Title %v", i+1),
			Description: fmt.Sprintf("Description %v", i+1),
			Price:       rand.Float64() * 1000,
		}

		pCnt := rand.Intn(10) + 1
		for p := 0; p < pCnt; p++ {
			adv.Photos = append(adv.Photos, fmt.Sprintf("https://photos.com/photo_%v_%v", i+1, p+1))
		}
		adv.Id, err = repo.InsertAdvert(ctx, adv)
		if err != nil {
			t.Error(err)
			return
		}
		advs = append(advs, adv)
	}

	// Тестиррование выбора одниочного объявления
	for _, adv := range advs {
		_ = adv
		rAdv, err := repo.GetAdvert(ctx, adv.Id)
		if err != nil {
			fmt.Println(err.Error())
			t.Error(err)
			return
		}
		if rAdv.Description != adv.Description {
			t.Error("Descriptions isn't equal")
			return
		}
		if rAdv.Title != adv.Title {
			t.Error("Titles isn't equal")
			return
		}
		if rAdv.Price != adv.Price {
			t.Error("Prices isn't equal")
			return
		}
		if !reflect.DeepEqual(rAdv.Photos, adv.Photos) {
			t.Error("Photos isn't equal")
			return
		}
	}
}
