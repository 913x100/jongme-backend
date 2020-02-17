package api

import (
	"encoding/json"
	"jongme/app/errs"
	"jongme/app/model"

	"github.com/valyala/fasthttp"
)

type UserDatabase interface {
	GetUsers() ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUser(t interface{}) (*model.User, error)
}

type UserAPI struct {
	DB UserDatabase
}

func (u *UserAPI) GetUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	users, err := u.DB.GetUsers()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Cannot get users from database")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(users)
}

func (u *UserAPI) GetUserByID(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	id := ctx.UserValue("id").(string)

	user, err := u.DB.GetUserByID(id)

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(user)

	return nil
}

func (u *UserAPI) GetUser(ctx *fasthttp.RequestCtx) error {
	// TODO edit error
	// TODO add struct validation
	ctx.SetContentType("application/json")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	// id, err := withID(ctx, "id")
	// if err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	// }

	// user, err := u.DB.GetUserByID(id)
	// id := ctx.UserValue("id")

	// req := make(map[string]string)
	// req["user_id"] = id.(string)
	// // req := getUserRequest{UserID: id.(string)}

	// user, err := u.DB.GetUser(req)

	// if err != nil {
	// 	return errs.NewHTTPError(err, 500, "Internal server error.")
	// }
	// ctx.SetStatusCode(fasthttp.StatusOK)
	// json.NewEncoder(ctx).Encode(user)

	return nil

	// return withID(ctx, "id", func(id primitive.ObjectID) error {
	// 	user, err := u.DB.GetUserByID(id)

	// 	if err != nil {
	// 		return errs.NewHTTPError(err, 500, "Internal server error.")
	// 	}
	// 	return user
	// })
}
