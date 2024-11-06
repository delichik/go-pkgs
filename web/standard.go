package web

import (
	"context"
	"net/http"

	"vap/internal/das/dao"
	"vap/internal/das/utils"
	"vap/pkg/null"

	"github.com/delichik/go-pkgs/db"
)

type commonGetReq struct {
	ID null.Int `json:"id" validate:"required,min=1"`
}

type commonDeleteReq struct {
	IDs []int64 `json:"ids" validate:"required"`
}

// G Get
func G[M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(commonGetReq)
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		rsp, err := db.CommonGet[M](ctx, req.ID.Int64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}

// S Scan
func S[M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(dao.ScanReq)
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		rsp, err := db.CommonScan[M](ctx, req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}

// D Delete
func D[M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(commonDeleteReq)
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		rsp, err := db.CommonDeleteBatch[M](ctx, req.IDs)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}

// U Update
func U[REQ any, M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(utils.BatchRequest[REQ])
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		items := utils.TransferItems[REQ, M](req.Items)
		rsp, err := db.CommonUpdateBatch(ctx, items)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}

// C Create
func C[REQ any, M db.ModelInterface](g *echo.Group, path string) {
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(utils.BatchRequest[REQ])
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		items := utils.TransferItems[REQ, M](req.Items)
		err := db.CommonCreateBatch(ctx, items)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"items": items})
	}))
}

// CU Create or Update
func CU[REQ any, M db.ModelInterface](g *echo.Group, path string) {
	var m M
	g.POST(path, h(func(c echo.Context, ctx context.Context) error {
		req := new(utils.BatchRequest[REQ])
		if err := parse(c, req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		items := utils.TransferItems[REQ, M](req.Items)
		rsp, err := db.CommonUpsertBatch(ctx, items)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, rsp)
	}))
}
