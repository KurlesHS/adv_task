package postgres

import (
	"context"
	"fmt"
	"kurles/adv_task/configs"
	"kurles/adv_task/pkg/model"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdvertRepo(t *testing.T) {
	cfg, err := configs.LoadConfig()
	require.Nil(t, err)

	repo, err := New(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUserName, cfg.DBPassword)
	require.Nil(t, err)

	ctx := context.Background()
	err = repo.ClearAllAdverts(ctx)
	require.Nil(t, err)

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
		require.Nil(t, err)

		advs = append(advs, adv)
	}

	// Тестиррование выбора одниочного объявления
	for _, adv := range advs {
		_ = adv
		rAdv, err := repo.GetAdvert(ctx, adv.Id)
		require.Nil(t, err)

		require.Equal(t, rAdv.Description, adv.Description, "Descriptions isn't equal")

		require.Equal(t, rAdv.Title, adv.Title, "Titles isn't equal")

		require.Equal(t, rAdv.Price, adv.Price, "Prices isn't equal")

		require.Equal(t, rAdv.Photos, adv.Photos, "Photos isn't equal")
	}

	// Тестирование постраничного выбора

	checkPage := func(advs []model.DetailedAdvert, sortBy model.SortBy, sortOrder model.SortOrder) error {
		for i := 0; i < int(advCnt); i += 10 {
			res, err := repo.GetAdverts(ctx, i/10+1, sortBy, sortOrder)

			require.Nil(t, err)

			require.True(t, len(res) != 0, "no adv results")

			for ai, adv := range res {
				require.False(t, i+ai >= int(advCnt), "too many adv results")

				actAdv := advs[i+ai]

				require.Equal(t, adv.Title, actAdv.Title, fmt.Sprintf("wrong adv page result (title) at page %v and #%v", i/10, ai))

				require.Equal(t, adv.Price, actAdv.Price, fmt.Sprintf("wrong adv page result (price) at page %v and #%v", i/10, ai))

				require.Equal(t, adv.MainPhoto, actAdv.Photos[0], fmt.Sprintf("wrong adv page result (photo) at page %v and #%v", i/10, ai))
			}
		}
		return nil
	}
	// сортировка по дате по возрастанию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Id < advs[j].Id
	})

	require.Nil(t, checkPage(advs, model.Date, model.Asc))

	// сортировка по дате по убыванию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Id >= advs[j].Id
	})

	require.Nil(t, checkPage(advs, model.Date, model.Desc))

	// сортировка по цене по возрастанию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Price < advs[j].Price
	})

	require.Nil(t, checkPage(advs, model.Price, model.Asc))

	// сортировка по дате по убыванию
	sort.Slice(advs, func(i, j int) bool {
		return advs[i].Price >= advs[j].Price
	})

	require.Nil(t, checkPage(advs, model.Price, model.Desc))
}
