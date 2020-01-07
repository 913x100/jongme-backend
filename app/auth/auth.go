package auth

import (
	"encoding/json"
	"jongme/app/config"
	"jongme/app/errs"
	"jongme/app/model"
	"jongme/app/utils"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

type AuthDatabase interface {
	CreateUser(user *model.User) error
}

type Auth struct {
	DB AuthDatabase
	FB utils.FB
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (a *Auth) GenToken(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json; charset=utf-8")

	code := ctx.URI().QueryArgs().Peek("code")

	if len(code) == 0 {

		return errs.NewHTTPError(nil, 400, "Bad request: 'code' query param is missing.")

		// ctx.SetStatusCode(fasthttp.StatusBadRequest)
		// ctx.SetBodyString(`{ "code": "c1101","message": "'code' query param is missing", "success": false }`)
		// return
	}

	res := a.FB.OauthWithCode(string(code))
	if res.Err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1102).")
		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1102","message": "Login Failed", "success": false }`)
		// return
	}
	accessToken, ok := res.Response["access_token"].(string)
	if !ok {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1103).")

		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1103","message": "Login Failed", "success": false }`)
		// return
	}

	res = a.FB.GetLongLiveToken(accessToken)
	if res.Err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1104).")

		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1104","message": "Login Failed", "success": false }`)
		// return
	}
	accessToken, ok = res.Response["access_token"].(string)
	if !ok {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1105).")

		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1105","message": "Login Failed", "success": false }`)
		// return
	}

	res = a.FB.GetUserInfo(accessToken)
	if res.Err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1106).")

		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1106","message": "Login Failed", "success": false }`)
		// return
	}

	userID, ok := res.Response["id"].(string)
	if !ok {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c11027).")

		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1107","message": "Login Failed", "success": false }`)
		// return
	}
	name, _ := res.Response["name"].(string)

	user := model.User{
		UserID:      userID,
		AccessToken: accessToken,
		Name:        name,
	}

	_ = a.DB.CreateUser(user.New())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(720 * time.Hour).Unix(), // token expire token
			Issuer:    config.TokenIssuer,
		},
	})
	tokenString, err := token.SignedString(config.TokenSecret)
	if err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1108).")

		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		// ctx.SetBodyString(`{ "code": "c1111","message": "Login Failed", "success": false }`)
		// return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ouput := tokenResponse{Token: tokenString}
	response, _ := json.Marshal(ouput)
	ctx.Write(response)
	return nil
	// ctx.SuccessString("application/json", `{"code": "c1112","success":true,"message":"Authentication Successful","token":"`+tokenString+`"}`)
}
