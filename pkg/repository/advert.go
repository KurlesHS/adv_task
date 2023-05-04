package repository

import (
	"context"
	"kurles/adv_task/pkg/model"
)

type AdvertRepository interface {
	GetAdverts(ctx context.Context, page int, sortBy model.SortBy, sortOrder model.SortOrder) ([]model.Advert, error)
	GetAdvert(ctx context.Context, advId int64) (model.DetailedAdvert, error)
	InsertAdvert(ctx context.Context, advert model.DetailedAdvert) (int64, error)
}
