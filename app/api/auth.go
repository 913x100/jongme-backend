package api

import (
	"encoding/json"
	"fmt"
	"jongme/app/config"
	"jongme/app/errs"
	"jongme/app/fbbot"
	"jongme/app/model"
	"jongme/app/network"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasthttp"
)

type AuthDatabase interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	CreatePage(page *model.Page) (*model.Page, error)
	GetPageByID(id string) (*model.Page, error)
	UpdatePage(page interface{}) (*model.Page, error)
}

type AuthFb interface {
	OauthWithCode(code string) network.Response
	GetLongLiveToken(accessToken string) network.Response
	SubscribeWebhook(pageID, accessToken string) network.Response
	GetUserInfo(userToken string) network.Response
	GetPages(userToken string) network.Response
	GetPageToken(pageID, userToken string) network.Response
	EnabledGetStarted(pageAccessToken string) network.Response
	AddPersistentMenus(pageAccessToken string, menus ...*fbbot.Menu) network.Response
}

type AuthBot interface {
}

type Auth struct {
	DB AuthDatabase
	FB AuthFb
	// Bot AuthBot
}

type tokenResponse struct {
	Token string `json:"token"`
}

type pageTokenResponse struct {
	Name        string `json:"name"`
	PageID      string `json:"id"`
	AccessToken string `json:"access_token"`
}

// type Payload struct {
// 	StepID  int    `json:"step_id"`
// 	Payload string `json:"payload"`
// }

func (a *Auth) GenToken(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json; charset=utf-8")

	code := ctx.QueryArgs().Peek("code")
	if len(code) == 0 {
		return errs.NewHTTPError(nil, 400, "Bad request: 'code' query param is missing.")
	}

	res := a.FB.OauthWithCode(string(code))
	if res.Err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1102).")
	}
	accessToken, ok := res.Response["access_token"].(string)
	if !ok {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1103).")
	}

	res = a.FB.GetLongLiveToken(accessToken)
	if res.Err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1104).")
	}
	accessToken, ok = res.Response["access_token"].(string)
	if !ok {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1105).")
	}

	res = a.FB.GetUserInfo(accessToken)
	if res.Err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1106).")
	}

	userID, ok := res.Response["id"].(string)
	if !ok {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c11027).")
	}
	name, _ := res.Response["name"].(string)

	user := model.User{
		UserID:      userID,
		AccessToken: accessToken,
		Name:        name,
		UpdatedOn:   time.Now(),
	}

	err := a.DB.CreateUser(user.New())

	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.AuthenticateClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(720 * time.Hour).Unix(), // token expire token
			Issuer:    config.TokenIssuer,
		},
	})

	tokenString, err := token.SignedString(config.TokenSecret)
	fmt.Println(tokenString)
	if err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1108).")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	output := tokenResponse{Token: tokenString}
	response, _ := json.Marshal(output)
	ctx.Write(response)
	return nil
}

func (a *Auth) GetPages(ctx *fasthttp.RequestCtx) error {
	//TODO Add error

	ctx.SetContentType("application/json; charset=utf-8")

	userID := string(ctx.Request.Header.Peek("UserID"))
	user, err := a.DB.GetUserByID(userID)

	if err != nil {
		return nil
	}

	result := a.FB.GetPages(user.AccessToken).Response

	pages, _ := json.Marshal(result["data"])
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(pages)

	return nil
}

func (a *Auth) SelectPage(ctx *fasthttp.RequestCtx) error {
	//TODO Add error
	ctx.SetContentType("application/json; charset=utf-8")

	pageID := ctx.UserValue("id").(string)

	userID := string(ctx.Request.Header.Peek("UserID"))
	user, err := a.DB.GetUserByID(userID)

	if err != nil {
		return nil
	}

	result := a.FB.GetPageToken(pageID, user.AccessToken)

	if result.Err != nil {
		return errs.NewHTTPError(result.Err, 500, "Internal server error")
	}

	var pageWithToken pageTokenResponse

	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &pageWithToken,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	decoder.Decode(result.Response)

	page, err := a.DB.GetPageByID(pageWithToken.PageID)
	// if err != nil {
	// 	return errs.NewHTTPError(err, 500, "Internal server error")
	// }
	if page == nil {
		fmt.Println("new")
		newPage := model.Page{
			Name:        pageWithToken.Name,
			PageID:      pageWithToken.PageID,
			AccessToken: pageWithToken.AccessToken,
			UpdatedOn:   time.Now(),
		}
		page, err = a.DB.CreatePage(newPage.New())
		//TODO check error

	} else {

		p := &model.UpdatePageToken{
			PageID:      pageWithToken.PageID,
			AccessToken: pageWithToken.AccessToken,
			Name:        pageWithToken.Name,
		}

		// page.PageID = pageWithToken.PageID
		// page.Name = pageWithToken.Name
		// page.AccessToken = pageWithToken.AccessToken
		// page.UpdatedOn = time.Now()
		_, err = a.DB.UpdatePage(p)
	}

	//Subscription webhook
	result = a.FB.SubscribeWebhook(page.PageID, page.AccessToken)

	//Enable get started button
	result = a.FB.EnabledGetStarted(page.AccessToken)
	//Add persistent menu
	menu := fbbot.NewMenu()

	menu.AddMenuItems(
		fbbot.NewPostbackMenuItem("จองบริการ", fmt.Sprintf(`{ "step_id":%d, "page_id":"%s"}`, 1, page.PageID)),
		fbbot.NewPostbackMenuItem("ยกเลิกการจอง", fmt.Sprintf(`{ "step_id":%d, "page_id":"%s"}`, -1, page.PageID)),
	)
	result = a.FB.AddPersistentMenus(page.AccessToken, menu)

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.AuthenticateClaims{
		UserID: userID,
		PageID: pageID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(720 * time.Hour).Unix(), // token expire token
			Issuer:    config.TokenIssuer,
		},
	})

	tokenString, err := token.SignedString(config.TokenSecret)

	if err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1108).")
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)

	output := tokenResponse{Token: tokenString}
	response, _ := json.Marshal(output)
	ctx.Write(response)

	return nil
}
