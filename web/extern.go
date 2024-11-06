package web

import (
	"context"
	"net/http"

	"vap/internal/das/utils"

	"github.com/delichik/go-pkgs/db"
)

type commonGetBatchReq struct {
	IDs []int64 `json:"ids" validate:"required,min=1,dive"`
}

// GB Get Batch
func GB[M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(commonGetBatchReq)
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		rsp, err := db.CommonGetBatch[M](ctx, req.IDs)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}

// UW Update Where
func UW[REQ any, M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(utils.UpdateWhereRequest[REQ])
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		m := utils.CopyContent[REQ, M](req.Item)
		rsp, err := db.CommonUpdateWhere(ctx, m, req.Wheres)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}

// DW Delete Where
func DW[M db.ModelInterface](g *echo.Group, path string) {
	var m M
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(utils.DeleteWhereRequest)
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		rsp, err := db.CommonDeleteWhere[M](ctx, req.Wheres)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}
