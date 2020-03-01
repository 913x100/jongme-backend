package api

import (
	"encoding/json"

	"jongme/app/config"
	"jongme/app/errs"
	"jongme/app/model"
	"jongme/app/network"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PageDatabase interface {
	GetPageByID(id string) (*model.Page, error)
	CreatePage(page *model.Page) (*model.Page, error)
	GetPages() ([]*model.Page, error)
	UpdatePage(page interface{}) (*model.Page, error)
}

type PageFb interface {
	GetPages(userToken string) network.Response
	GetPageToken(pageID, userToken string) network.Response
}

type PageAPI struct {
	DB PageDatabase
	FB PageFb
	// Validate *validator.Validate
}

type getUserRequest struct {
	UserID string `json:"user_id"`
}

type getPageRequest struct {
	UserID string `json:"user_id"`
	PageID string `json:"page_id"`
}

func (p *PageAPI) GetPages(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json;charset=utf-8")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	pages, err := p.DB.GetPages()

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(pages)
	return nil
}

func (p *PageAPI) GetPageByID(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json;charset=utf-8")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	pageID := ctx.UserValue("id").(string)
	// fmt.Println(pageID)

	page, err := p.DB.GetPageByID(pageID)

	type output struct {
		ID     primitive.ObjectID `son:"_id"`
		PageID string             `json:"page_id,omitempty"`
		Name   string             `json:"name,omitempty"`
		// PageHours []*model.PageHours `json:"page_hours,omitempty"`
		StartTime  string    `json:"start_time"`
		EndTime    string    `json:"end_time"`
		IsActive   bool      `json:"is_active"`
		IsBreak    bool      `json:"is_break"`
		BreakStart string    `json:"break_start"`
		BreakEnd   string    `json:"break_end"`
		Sun        bool      `json:"sun"`
		Mon        bool      `json:"mon"`
		Tue        bool      `json:"tue"`
		Wed        bool      `json:"wed"`
		Thu        bool      `json:"thu"`
		Fri        bool      `json:"fri"`
		Sat        bool      `json:"sat"`
		UpdatedOn  time.Time `json:"updated_on,omitempty"`
		CreatedOn  time.Time `json:"created_on,omitempty"`
	}

	pageOutput := &output{
		ID:         page.ID,
		PageID:     page.PageID,
		Name:       page.Name,
		StartTime:  page.StartTime,
		EndTime:    page.EndTime,
		BreakStart: page.BreakStart,
		BreakEnd:   page.BreakEnd,
		Sun:        page.Sun,
		Mon:        page.Mon,
		Tue:        page.Tue,
		Wed:        page.Wed,
		Thu:        page.Thu,
		Fri:        page.Fri,
		Sat:        page.Sat,
		// DayOfWeek:  page.DayOfWeek,
		// PageHours: page.PageHours,
		IsActive:  page.IsActive,
		IsBreak:   page.IsBreak,
		UpdatedOn: page.UpdatedOn,
	}

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(pageOutput)
	return nil
}

func (p *PageAPI) GetExpireToken(ctx *fasthttp.RequestCtx) error {
	expirationTime := time.Now().Add(time.Minute)
	claims := &model.ExpirePageClaims{
		UserID: "1234",
		PageID: "2345",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.TokenSecret)

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	ctx.Write([]byte(tokenString))
	return nil
}

func (p *PageAPI) GetExpirePage(ctx *fasthttp.RequestCtx) error {

	ctx.Write([]byte("Hi"))
	return nil
}

func (p *PageAPI) UpdatePage(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsPut() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}
	// fmt.Println("Update")
	input := model.UpdatePage{}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}

	_ = ctx.UserValue("id").(string)

	_, err := p.DB.UpdatePage(&input)

	if err != nil {
		return errs.NewHTTPError(err, 404, "service down not exists.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
