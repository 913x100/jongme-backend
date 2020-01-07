package guard

import (
	"bytes"
	"fmt"

	"jongme/app/config"
	"jongme/app/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

func Auth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		auth := ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix(auth, config.BasicAuthPrefix) {
			pair := bytes.SplitN(auth, []byte(" "), 2)
			if len(pair) == 2 {
				// Delegate request to the given handle
				claims := model.MyCustomClaims{}

				token, err := jwt.ParseWithClaims(string(pair[1]), &claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(config.TokenSecret), nil
				})
				if err != nil {
					ctx.Error("Unauthorized Access", fasthttp.StatusForbidden)
					return
				}
				if !token.Valid {
					ctx.Error("Unauthorized Access", fasthttp.StatusForbidden)
					return
				}
				ctx.Request.Header.Set("UserID", claims.UserID)
				ctx.Request.Header.Set("AccessToken", claims.AccessToken)
				ctx.Request.Header.Set("ServiceID", "21212")
				ctx.SetContentType("application/json")
				h(ctx)
				return
			}
		}
		ctx.Error("Unauthorized Access", fasthttp.StatusForbidden)
	})
}
