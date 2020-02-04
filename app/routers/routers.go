package routers

import (
	"jongme/app/api"
	"jongme/app/auth"
	"jongme/app/config"
	"jongme/app/database"
	"jongme/app/errs"
	"jongme/app/model"
	"jongme/app/network/fasthttp_client"
	"jongme/app/utils"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gopkg.in/go-playground/validator.v10"
)

type RequestHandler func(*fasthttp.RequestCtx) error

func rootHandler(h RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		err := h(ctx)

		if err == nil {
			return
		}

		log.Printf("An error accured: %v", err) // Log the error.

		clientError, ok := err.(errs.ClientError) // Check if it is a ClientError.

		if !ok {
			ctx.SetStatusCode(500)
			return
		}

		body, err := clientError.ResponseBody() // Try to get response body of ClientError.

		if err != nil {
			log.Printf("An error accured: %v", err)
			ctx.SetStatusCode(500)
			return
		}
		status, headers := clientError.ResponseHeaders() // Get http status code and headers.

		for k, v := range headers {
			ctx.Response.Header.Set(k, v)
		}
		ctx.SetStatusCode(status)
		ctx.Write(body)
	})
}

func isExpired(h RequestHandler) RequestHandler {
	return func(ctx *fasthttp.RequestCtx) error {
		tokenString := ctx.UserValue("page_token").(string)
		claims := &model.ExpirePageClaims{}

		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.TokenSecret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return errs.NewHTTPError(err, fasthttp.StatusUnauthorized, "Unauthorized.")
			}
			return errs.NewHTTPError(err, 400, "Bad request.")
		}

		if !tkn.Valid {
			return errs.NewHTTPError(err, fasthttp.StatusUnauthorized, "Unauthorized.")
		}

		h(ctx)
		return nil
	}
}

func Create(db *database.Mongo) *router.Router {
	client := &fasthttp.Client{MaxConnsPerHost: 2048}
	fasthttpClient := &fasthttp_client.FastHTTPClient{client}

	authHandler := auth.Auth{DB: db, FB: utils.FB{fasthttpClient}}
	userHandler := api.UserAPI{DB: db}
	serviceHandler := api.ServiceAPI{DB: db, Validate: validator.New()}
	pageHandler := api.PageAPI{DB: db, FB: utils.FB{fasthttpClient}}
	facebookHandler := api.FbAPI{DB: db}

	r := router.New()

	sv := r.Group("/api")

	auth := sv.Group("/auth")
	auth.GET("/fb", rootHandler(authHandler.GenToken))

	user := sv.Group("/user")
	user.GET("/", userHandler.GetUsers)
	user.GET("/:id/fb", rootHandler(userHandler.GetUser))
	user.GET("/:id", rootHandler(userHandler.GetUserByID))

	service := sv.Group("/service")
	service.GET("/", rootHandler(serviceHandler.GetServices))
	service.POST("/", rootHandler(serviceHandler.CreateService))
	service.PUT("/:id", rootHandler(serviceHandler.UpdateServiceByID))
	// auth.GET("/", Index)

	page := sv.Group("/page")
	page.POST("/", rootHandler(pageHandler.CreatePage))
	page.GET("/", rootHandler(pageHandler.GetPages))
	page.GET("/token", rootHandler(pageHandler.GetExpireToken))
	page.GET("/expire/:page_token", rootHandler(isExpired(pageHandler.GetExpirePage)))
	// page.GET("/:id/id", rootHandler(pageHandler.GetPagesByFacebookID))

	facebook := sv.Group("/fb")
	facebook.GET("/webhook", facebookHandler.Webhook)
	facebook.POST("/webhook", facebookHandler.RecieveWebhook)
	return r
}
