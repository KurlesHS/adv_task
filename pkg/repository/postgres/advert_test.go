package postgres

import (
	"context"
	"fmt"
	"kurles/adv_task/configs"
	"kurles/adv_task/pkg/model"
	"math/rand"
	"reflect"
	"sort"
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

	ctx := context.Background()
	err = repo.ClearAllAdverts(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	advs := make([]model.DetailedAdvert, 0)
	advCnt := 100 + rand.Int31n(50)
	for i := 0; i < int(advCnt); i++ {
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

	// Тестирование постраничного выбора

	checkPage := func(advs []model.DetailedAdvert, sortBy model.SortBy, desc bool) error {
		for i := 0; i < int(advCnt); i += 10 {
			res, err := repo.GetAdverts(ctx, i/10+1, sortBy, desc)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("no adv results")
			}
			for ai, adv := range res {
				if i+ai >= int(advCnt) {
					return fmt.Errorf("too many adv results")
				}
				actAdv := advs[i+ai]
				if adv.Title != actAdv.Title {
					return fmt.Errorf("wrong adv page result (title) at page %v and #%v", i/10, ai)

				}
				if adv.Price != actAdv.Price {
					return fmt.Errorf("wrong adv page result (price) at page %v and #%v", i/10, ai)

				}
				if adv.MainPhoto != actAdv.Photos[0] {
					return fmt.Errorf("wrong adv page result (photo) at page %v and #%v", i/10, ai)
				}
			}
		}
		return nil
	}
	// сортировка по дате по возрастанию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Id < advs[j].Id
	})

	err = checkPage(advs, model.Date, false)
	if err != nil {
		t.Error(err)
		return
	}

	// сортировка по дате по убыванию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Id > advs[j].Id
	})

	err = checkPage(advs, model.Date, true)
	if err != nil {
		t.Error(err)
		return
	}

	// сортировка по цене по возрастанию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Price < advs[j].Price
	})

	err = checkPage(advs, model.Price, false)

	// сортировка по дате по убыванию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Price > advs[j].Price
	})

	err = checkPage(advs, model.Price, true)
	if err != nil {
		t.Error(err)
		return
	}
}
