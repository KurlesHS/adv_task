package handler

import (
	"context"
	"fmt"
	"kurles/adv_task/pkg/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BadRequestResponse struct {
	Result string `json:"result"`
}

type Service interface {
	GetAdverts(ctx context.Context, page int, sortBy model.SortBy, sortOrder model.SortOrder) ([]model.Advert, error)
	GetAdvert(ctx context.Context, advId int64) (model.DetailedAdvert, error)
	InsertAdvert(ctx context.Context, advert model.DetailedAdvert) (int64, error)
}

type Handlers struct {
	e       *echo.Echo
	port    int
	service *Service
}

func New(service *Service, port int) Handlers {
	res := Handlers{
		e:       echo.New(),
		port:    port,
		service: service,
	}
	res.e.POST("/api/adverts", res.InsertAdvert)
	res.e.GET("/api/adverts", res.GetAdverts)
	res.e.GET("/api/adverts/:id", res.GetAdvert)
	res.e.Any("*", func(c echo.Context) error {
		r := &BadRequestResponse{
			Result: "bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	})
	return res
}

func (h Handlers) InsertAdvert(ctx echo.Context) error {

	return nil
}

func (h Handlers) GetAdvert(ctx echo.Context) error {
	return nil
}

func (h Handlers) GetAdverts(ctx echo.Context) error {
	return nil
}

func (h *Handlers) Start() error {
	// TODO: make graceful shutdown
	// https://echo.labstack.com/cookbook/graceful-shutdown/
	return h.e.Start(fmt.Sprintf(":%d", h.port))
}
