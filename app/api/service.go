package api

import (
	"encoding/json"
	"jongme/app/errs"
	"jongme/app/model"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v10"
)

type ServiceDatabase interface {
	CreateService(service *model.Service) error
	GetServices(page *model.Paging) ([]*model.Service, error)
	UpdateService(service *model.Service) error
	DeleteServiceByID(id primitive.ObjectID) error
}

type createServiceRequest struct {
	PageID            string `json:"page_id" validate:"required"`
	Name              string `json:"name" validate:"required"`
	Type              string `json:"type"`
	Quantity          int64  `json:"quantity"`
	MinimumTimeLength int64  `json:"minimum_time_length"`
	StartTime         int64  `json:"start_time"`
	EndTime           int64  `json:"end_time"`
}

type updateServiceRequest struct {
	ID                primitive.ObjectID `json:"_id" validate:"required"`
	PageID            string             `json:"page_id" validate:"required"`
	Name              string             `json:"name" validate:"required"`
	Type              string             `json:"type"`
	Quantity          int64              `json:"quantity"`
	MinimumTimeLength int64              `json:"minimum_time_length"`
	StartTime         int64              `json:"start_time"`
	EndTime           int64              `json:"end_time"`
}

type serviceResponse struct {
	ID                primitive.ObjectID `json:"_id"`
	PageID            string             `json:"page_id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Quantity          int64              `json:"quantity"`
	MinimumTimeLength int64              `json:"minimum_time_length"`
	StartTime         int64              `json:"start_time"`
	EndTime           int64              `json:"end_time"`
}

type ServiceAPI struct {
	DB       ServiceDatabase
	Validate *validator.Validate
}

func (s *ServiceAPI) CreateService(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json; charset=utf-8")
	if !ctx.IsPost() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}
	input := createServiceRequest{}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}

	if err := s.Validate.Struct(input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	}

	service := model.Service{
		PageID:            input.PageID,
		Name:              input.Name,
		Type:              input.Type,
		Quantity:          input.Quantity,
		MinimumTimeLength: input.MinimumTimeLength,
		StartTime:         input.StartTime,
		EndTime:           input.EndTime,
	}

	err := s.DB.CreateService(service.New())
	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

func (s *ServiceAPI) GetServices(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json;charset=utf-8")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	var (
		start int64  = 0
		end   int64  = 10
		sort  string = "_id"
		order int    = 1
	)
	if tmp := string(ctx.FormValue("_start")); tmp != "" {
		start, _ = strconv.ParseInt(tmp, 10, 64)
	}
	if tmp := string(ctx.FormValue("_end")); tmp != "" {
		end, _ = strconv.ParseInt(tmp, 10, 64)
	}
	if tmp := string(ctx.FormValue("_sort")); tmp != "" {
		sort = tmp
	}

	if sort == "id" {
		sort = "_id"
	}

	if tmp := string(ctx.FormValue("_order")); tmp != "" {
		order = -1
	}

	limit := end - start
	users, err := s.DB.GetServices(
		&model.Paging{
			Skip:      &start,
			Limit:     &limit,
			SortKey:   sort,
			SortVal:   order,
			Condition: nil,
		})

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(users)
	return nil
}

func (s *ServiceAPI) UpdateServiceByID(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsPut() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	input := updateServiceRequest{}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}

	if err := s.Validate.Struct(input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	}

	service := model.Service{
		ID:                input.ID,
		PageID:            input.PageID,
		Name:              input.Name,
		Type:              input.Type,
		Quantity:          input.Quantity,
		MinimumTimeLength: input.MinimumTimeLength,
		StartTime:         input.StartTime,
		EndTime:           input.EndTime,
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return withID(ctx, "id", func(id primitive.ObjectID) error {
		if err := s.DB.UpdateService(&service); err != nil {
			return errs.NewHTTPError(err, 404, "service down not exists.")
		}
		return nil
	})
}

func (s *ServiceAPI) DeleteServiceByID(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsDelete() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	withID(ctx, "_id", func(id primitive.ObjectID) error {
		if err := s.DB.DeleteServiceByID(id); err != nil {
			return errs.NewHTTPError(err, 500, "Internal server error.")
		}
		return nil
	})

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
