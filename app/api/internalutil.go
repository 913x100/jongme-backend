package api

import (
	"jongme/app/errs"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func withID(ctx *fasthttp.RequestCtx, name string, f func(id primitive.ObjectID) error) error {
	// param := ctx.URI().QueryArgs().Peek(name)
	param := ctx.UserValue(name)

	id, err := primitive.ObjectIDFromHex(param.(string))

	if err == nil {
		return f(id)
	} else {
		return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	}
}
