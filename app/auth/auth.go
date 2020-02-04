package auth

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

type AuthDatabase interface {
	CreateUser(user *model.User) error
}

type AuthFb interface {
	OauthWithCode(code string) network.Response
	GetLongLiveToken(accessToken string) network.Response
	GetUserInfo(userToken string) network.Response
}

type Auth struct {
	DB AuthDatabase
	FB AuthFb
}

type tokenResponse struct {
	Token string `json:"token"`
}

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
	}

	_ = a.DB.CreateUser(user.New())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.AuthenticateClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(720 * time.Hour).Unix(), // token expire token
			Issuer:    config.TokenIssuer,
		},
	})
	tokenString, err := token.SignedString(config.TokenSecret)
	if err != nil {
		return errs.NewHTTPError(nil, 500, "Internal server error: Login failed (c1108).")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ouput := tokenResponse{Token: tokenString}
	response, _ := json.Marshal(ouput)
	ctx.Write(response)
	return nil
}
