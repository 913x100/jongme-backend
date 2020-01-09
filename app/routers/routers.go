package routers

import (
	"jongme/app/api"
	"jongme/app/auth"
	"jongme/app/database"
	"jongme/app/errs"
	"jongme/app/network/fasthttp_client"
	"jongme/app/utils"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gopkg.in/go-playground/validator.v10"
)

type Request func(*fasthttp.RequestCtx) error

func rootHandler(h Request) fasthttp.RequestHandler {
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

func Create(db *database.Mongo) *router.Router {
	client := &fasthttp.Client{MaxConnsPerHost: 2048}
	fasthttpClient := &fasthttp_client.FastHTTPClient{client}

	authHandler := auth.Auth{DB: db, FB: utils.FB{fasthttpClient}}
	userHandler := api.UserAPI{DB: db}
	serviceHandler := api.ServiceAPI{DB: db, Validate: validator.New()}

	r := router.New()

	sv := r.Group("/api")

	auth := sv.Group("/auth")
	auth.GET("/fb", rootHandler(authHandler.GenToken))

	user := sv.Group("/user")
	user.GET("/all", userHandler.GetUsers)

	service := sv.Group("/service")
	service.GET("/all", rootHandler(serviceHandler.GetServices))
	service.POST("/", rootHandler(serviceHandler.CreateService))
	service.PUT("/:id", rootHandler(serviceHandler.UpdateServiceByID))
	// auth.GET("/", Index)
	return r
}
