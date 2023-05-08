package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kurles/adv_task/pkg/error_message"
	"kurles/adv_task/pkg/model"
	"net/http"
	"strconv"
	"strings"

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
	service Service
}

func New(service Service, port int) Handlers {
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
	var adv model.DetailedAdvert
	var errorMess error_message.ErrorMessage
	err := json.NewDecoder(ctx.Request().Body).Decode(&adv)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &BadRequestResponse{Result: err.Error()})
	}
	id, err := h.service.InsertAdvert(ctx.Request().Context(), adv)

	if err != nil {
		status := echo.ErrInternalServerError.Code
		if errors.As(err, &errorMess) && errorMess.Type() == error_message.NotFound {
			status = echo.ErrNotFound.Code
		}
		return ctx.JSON(status, &BadRequestResponse{Result: err.Error()})
	}
	jsonMap := make(map[string]interface{})
	jsonMap["id"] = id
	return ctx.JSON(http.StatusOK, &jsonMap)
}

func (h Handlers) GetAdvert(ctx echo.Context) error {
	var errorMess error_message.ErrorMessage
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(echo.ErrBadRequest.Code, &BadRequestResponse{Result: err.Error()})
	}
	adv, err := h.service.GetAdvert(ctx.Request().Context(), id)
	// TODO: разделить ответы на:
	// ошибка БД, запись не найдена, внутренняя ошибка
	if err != nil {
		status := echo.ErrInternalServerError.Code
		if errors.As(err, &errorMess) && errorMess.Type() == error_message.NotFound {
			status = echo.ErrNotFound.Code
		}
		return ctx.JSON(status, &BadRequestResponse{Result: err.Error()})
	}
	if strings.ToLower(ctx.QueryParam("fields")) != "true" {
		adv.Description = ""
		adv.Photos = adv.Photos[:1]
	}
	return ctx.JSON(http.StatusOK, &adv)
}

func (h Handlers) GetAdverts(ctx echo.Context) error {
	var errorMess error_message.ErrorMessage
	pageStr := ctx.QueryParam("page")
	page := int64(1)
	if len(pageStr) > 0 {
		var err error
		page, err = strconv.ParseInt(pageStr, 10, 31)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &BadRequestResponse{Result: "page is not number"})
		}
	}
	sortParam := strings.Split(ctx.QueryParam("sort"), "_")

	sortBy := model.Date
	sortOrder := model.Asc
	if len(sortParam) > 0 && sortParam[0] == "price" {
		sortBy = model.Price
	}
	if len(sortParam) > 1 && sortParam[1] == "desc" {
		sortOrder = model.Desc
	}
	advs, err := h.service.GetAdverts(ctx.Request().Context(), int(page), sortBy, sortOrder)

	if err != nil {
		if err != nil {
			var status int
			if errors.As(err, &errorMess) {
				switch errorMess.Type() {
				case error_message.BadRequest:
					status = echo.ErrBadRequest.Code
				case error_message.NotFound:
					status = echo.ErrNotFound.Code
				default:
					status = echo.ErrInternalServerError.Code
				}
			}
			return ctx.JSON(status, &BadRequestResponse{Result: err.Error()})
		}
		return ctx.JSON(http.StatusBadRequest, &BadRequestResponse{Result: err.Error()})
	}

	return ctx.JSON(http.StatusOK, &advs)
}

func (h *Handlers) Start() error {
	// TODO: make graceful shutdown
	// https://echo.labstack.com/cookbook/graceful-shutdown/
	return h.e.Start(fmt.Sprintf(":%d", h.port))
}
