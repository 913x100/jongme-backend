package routers

import (
	"bytes"
	"fmt"
	"jongme/app/api"
	"jongme/app/config"
	"jongme/app/database"
	"jongme/app/errs"
	"jongme/app/fbbot"
	"jongme/app/model"
	"jongme/app/network/fasthttp_client"

	// "jongme/app/utils"
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

func basicAuth(h RequestHandler) RequestHandler {
	return func(ctx *fasthttp.RequestCtx) error {
		auth := ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix(auth, config.BasicAuthPrefix) {
			pair := bytes.SplitN(auth, []byte(" "), 2)
			if len(pair) == 2 {
				// Delegate request to the given handle
				claims := model.AuthenticateClaims{}

				token, err := jwt.ParseWithClaims(string(pair[1]), &claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(config.TokenSecret), nil
				})
				if err != nil {
					ctx.Error("Unauthorized Access", fasthttp.StatusForbidden)
					return nil
				}
				if !token.Valid {
					ctx.Error("Unauthorized Access", fasthttp.StatusForbidden)
					return nil
				}
				ctx.Request.Header.Set("UserID", claims.UserID)
				// ctx.Request.Header.Set("AccessToken", claims.AccessToken)
				// ctx.Request.Header.Set("ServiceID", "21212")
				ctx.SetContentType("application/json")
				h(ctx)
				return nil
			}
		}
		ctx.Error("Unauthorized Access", fasthttp.StatusForbidden)
		return nil
	}
}

func Create(db *database.Mongo) *router.Router {
	client := &fasthttp.Client{MaxConnsPerHost: 2048}
	fasthttpClient := &fasthttp_client.FastHTTPClient{client}

	// bot := fbbot.New(fasthttpClient)

	authHandler := api.Auth{DB: db, FB: fbbot.FB{fasthttpClient}}
	userHandler := api.UserAPI{DB: db}
	serviceHandler := api.ServiceAPI{DB: db, Validate: validator.New()}
	pageHandler := api.PageAPI{DB: db, FB: fbbot.FB{fasthttpClient}}
	facebookHandler := api.FbAPI{DB: db, FB: fbbot.FB{fasthttpClient}}
	bookingHandler := api.BookingAPI{DB: db}

	r := router.New()

	sv := r.Group("/api")

	auth := sv.Group("/auth")
	auth.GET("/fb", rootHandler(authHandler.GenToken))
	auth.GET("/pages", rootHandler(basicAuth(authHandler.GetPages)))
	auth.GET("/page/:id", rootHandler(basicAuth(authHandler.SelectPage)))

	user := sv.Group("/user")
	user.GET("/", userHandler.GetUsers)
	user.GET("/:id", rootHandler(userHandler.GetUserByID))

	service := sv.Group("/service")
	// service.GET("/", rootHandler(serviceHandler.GetServices))
	// service.GET("/page/:id", rootHandler(serviceHandler.GetServicesByPage))
	service.GET("/", rootHandler(serviceHandler.GetServicesByFilter))
	service.GET("/slots/:id", rootHandler(serviceHandler.GetServicesSlots))
	service.POST("/", rootHandler(serviceHandler.CreateService))
	service.PUT("/:id", rootHandler(serviceHandler.UpdateServiceByID))
	service.DELETE("/:id", rootHandler(serviceHandler.DeleteServiceByID))
	// auth.GET("/", Index)

	page := sv.Group("/page")
	// page.POST("/", rootHandler(pageHandler.CreatePage))
	page.GET("/", rootHandler(pageHandler.GetPages))
	page.GET("/:id", rootHandler(pageHandler.GetPageByID))
	// page.GET("/token", rootHandler(pageHandler.GetExpireToken))
	// page.GET("/expire/:page_token", rootHandler(isExpired(pageHandler.GetExpirePage)))
	page.PUT("/:id", rootHandler(pageHandler.UpdatePage))
	// page.GET("/:id/id", rootHandler(pageHandler.GetPagesByFacebookID))

	facebook := sv.Group("/fb")
	facebook.GET("/webhook", facebookHandler.Webhook)
	facebook.POST("/webhook", facebookHandler.RecieveWebhook)
	facebook.GET("/", facebookHandler.SendSuccesMessage)

	booking := sv.Group("/booking")
	booking.POST("/", rootHandler(bookingHandler.CreateBooking))
	booking.GET("/", rootHandler(bookingHandler.GetBookingByService))
	booking.GET("/filter", rootHandler(bookingHandler.GetBookingByFilter))

	booking.PUT("/:id", rootHandler(bookingHandler.UpdateBookingByID))
	booking.DELETE("/:id", rootHandler(bookingHandler.DeleteBookingByID))
	return r
}
