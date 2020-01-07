package api

import (
	"encoding/json"
	"jongme/app/model"
	"strconv"

	"github.com/valyala/fasthttp"
)

type UserDatabase interface {
	GetUsers(page *model.Paging) ([]*model.User, error)
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
