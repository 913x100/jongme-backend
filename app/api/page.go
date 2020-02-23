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
)

type PageDatabase interface {
	GetUserByID(id string) (*model.User, error)
	GetUser(t interface{}) (*model.User, error)
	CreatePage(page *model.Page) (*model.Page, error)
	GetPages() ([]*model.Page, error)
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

// type createPageRequest struct {
// 	Name        string `json:"name"`
// 	AccessToken string `json:"access_token"`
// 	PageID      string `json:"id"`
// }

type getPageRequest struct {
	UserID string `json:"user_id"`
	PageID string `json:"page_id"`
}

func (p *PageAPI) GetPagesByFacebookID(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	id := ctx.UserValue("id")
	if id == "" {
		return errs.NewHTTPError(nil, 400, "Bad request: 'invalid FB id.")
	}

	req := getUserRequest{UserID: id.(string)}

	user, err := p.DB.GetUser(req)

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	user_token := user.AccessToken

	resp := p.FB.GetPages(user_token)

	pages := resp.Response
	// err = json.Unmarshal(resp.Response, &pages)

	// tmp := pages["data"].([]interface{})
	e, _ := json.Marshal(pages["data"])
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(e)
	return nil
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
