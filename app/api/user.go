package api

import (
	"encoding/json"
	"jongme/app/errs"
	"jongme/app/model"
	"strconv"

	"github.com/valyala/fasthttp"
)

type UserDatabase interface {
	GetUsers(page *model.Paging) ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUser(t interface{}) (*model.User, error)
}

type UserAPI struct {
	DB UserDatabase
}

func (u *UserAPI) GetUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	var (
		start int64  = 0
		end   int64  = 10
		sort  string = "_id"
		order int    = 1
	)
	if tmp := string(ctx.FormValue("_start")); tmp != "" {
		start, _ = strconv.ParseInt(tmp, 10, 64)
	}
	if tmp := string(ctx.FormValue("_end")); tmp != "" {
		end, _ = strconv.ParseInt(tmp, 10, 64)
	}
	if tmp := string(ctx.FormValue("_sort")); tmp != "" {
		sort = tmp
	}

	if sort == "id" {
		sort = "_id"
	}

	if tmp := string(ctx.FormValue("_order")); tmp != "" {
		order = -1
	}

	limit := end - start
	users, err := u.DB.GetUsers(
		&model.Paging{
			Skip:      &start,
			Limit:     &limit,
			SortKey:   sort,
			SortVal:   order,
			Condition: nil,
		})
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

	// id, err := withID(ctx, "id")
	id := ctx.UserValue("id").(string)
	// if err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	// }

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
