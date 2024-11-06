package web

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

const defaultTimeout = 4 * time.Second

func h(fn func(echo.Context, context.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		timeout := defaultTimeout

		t, err := strconv.ParseInt(c.QueryParam("timeout"), 10, 64)
		if err == nil {
			if t > 0 {
				timeout = time.Duration(t) * time.Second
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return fn(c, ctx)
	}
}

func parse(c echo.Context, v any) error {
	if err := c.Bind(v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
