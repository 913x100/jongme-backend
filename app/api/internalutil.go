package api

import (
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func withID(ctx *fasthttp.RequestCtx, name string) (primitive.ObjectID, error) {
	// param := ctx.URI().QueryArgs().Peek(name)
	param, ok := ctx.UserValue(name).(string)
	if !ok {
		return primitive.NewObjectID(), errors.New("Type assertion failed")
	}

	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return id, err
	}
	return id, nil
}
