package api

import (
	"encoding/json"
	"fmt"
	"jongme/app/errs"
	"jongme/app/model"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v10"
)

type ServiceDatabase interface {
	CreateService(service *model.Service) (*model.Service, error)
	GetServices(page *model.Paging) ([]*model.Service, error)
	GetServiceByID(id primitive.ObjectID) (*model.Service, error)
	// GetServicesByPage(pageID string, page *model.Paging) ([]*model.Service, error)
	GetServicesAccordingFilter(filter []bson.M) ([]*model.Service, error)
	UpdateService(service *model.Service) (*model.Service, error)
	DeleteServiceByID(id primitive.ObjectID) error
}

type createServiceRequest struct {
	PageID string `json:"page_id" validate:"required"`
}

type updateServiceRequest struct {
	ID                primitive.ObjectID `json:"_id"`
	PageID            string             `json:"page_id"`
	IsActive          bool               `json:"is_active"`
	Name              string             `json:"name"`
	ImageUrl          string             `json:image_url"`
	UnitType          string             `json:"unit_type"`
	UnitQuantity      int64              `json:"unit_quantity"`
	MinimumTimeLength int64              `json:"minimum_time_length"`
	IsTimeAdjust      bool               `json:"is_time_adjust"`
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
	fmt.Println("Create")
	input := model.Service{}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}

	// if err := s.Validate.Struct(input); err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	// }

	service, err := s.DB.CreateService(input.New())
	if err != nil {
		fmt.Println(err)
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	json.NewEncoder(ctx).Encode(service)
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
	services, err := s.DB.GetServices(
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
	json.NewEncoder(ctx).Encode(services)
	return nil
}

func (s *ServiceAPI) UpdateServiceByID(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsPut() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}
	// fmt.Println("update")
	input := model.Service{}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}

	// if err := s.Validate.Struct(input); err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	// }
	_, err := withID(ctx, "id")
	if err != nil {
		return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	}

	_, err = s.DB.UpdateService(&input)

	if err != nil {
		return errs.NewHTTPError(err, 404, "service down not exists.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func (s *ServiceAPI) DeleteServiceByID(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsDelete() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	id, err := withID(ctx, "id")
	if err != nil {
		return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	}

	if err := s.DB.DeleteServiceByID(id); err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func (s *ServiceAPI) GetServicesByFilter(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json;charset=utf-8")

	ID := string(ctx.FormValue("_id"))
	pageID := string(ctx.FormValue("page_id"))
	name := string(ctx.FormValue("name"))
	startTime := string(ctx.FormValue("start_time"))
	endTime := string(ctx.FormValue("end_time"))

	filter := []bson.M{}

	if ID != "" {
		id, _ := primitive.ObjectIDFromHex(ID)
		filter = append(filter, bson.M{"_id": bson.M{"$eq": id}})
	}
	if pageID != "" {
		filter = append(filter, bson.M{"page_id": bson.M{"$eq": pageID}})
	}
	if name != "" {
		filter = append(filter, bson.M{"name": bson.M{"$eq": name}})
	}
	if startTime != "" {
		filter = append(filter, bson.M{"start_time": bson.M{"$gte": startTime}})
	}
	if endTime != "" {
		filter = append(filter, bson.M{"end_time": bson.M{"$lte": endTime}})
	}
	services, err := s.DB.GetServicesAccordingFilter(filter)
	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(services)
	return nil
}

func (s *ServiceAPI) GetServicesSlots(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	id, _ := withID(ctx, "id")
	service, err := s.DB.GetServiceByID(id)

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	startTime, _ := time.Parse("15:04:05", service.StartTime)
	endTime, _ := time.Parse("15:04:05", service.EndTime)

	// fmt.Println(startTime, endTime)

	diff := int(endTime.Sub(startTime).Minutes())

	numSlot := diff / service.MinimumTimeLength

	var slots []string

	for i := 0; i < numSlot; i++ {
		a := startTime.Add(time.Minute * time.Duration(service.MinimumTimeLength*i))

		t := fmt.Sprintf("%02d:%02d:%02d", a.Hour(), a.Minute(), a.Second())

		slots = append(slots, t)
	}

	// fmt.Println(slots)
	// fmt.Println(startTime.Add(time.Minute * time.Duration(30)))

	// fmt.Println(time.Now().Local())

	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(slots)

	return nil
}
