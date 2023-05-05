package service

import (
	"context"
	"kurles/adv_task/pkg/model"
	"kurles/adv_task/pkg/repository"
)

type Service struct {
	repo repository.AdvertRepository
}

func New(repo repository.AdvertRepository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) GetAdverts(ctx context.Context, page int, sortBy model.SortBy, sortOrder model.SortOrder) ([]model.Advert, error) {
	return s.repo.GetAdverts(ctx, page, sortBy, sortOrder)
}
func (s Service) GetAdvert(ctx context.Context, advId int64) (model.DetailedAdvert, error) {
	return s.repo.GetAdvert(ctx, advId)
}
func (s Service) InsertAdvert(ctx context.Context, advert model.DetailedAdvert) (int64, error) {
	return s.repo.InsertAdvert(ctx, advert)
}
