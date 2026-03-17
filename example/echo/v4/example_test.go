package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FaruqRahmadani/ctxlog"
	ctxlogecho "github.com/FaruqRahmadani/ctxlog/http/echo/v4"

	"github.com/labstack/echo/v4"
)

func TestExampleEchoV4_ResponseAndLogging(t *testing.T) {
	e := echo.New()
	e.Use(ctxlogecho.New(
		ctxlogecho.WithAdditionalFields(func(c echo.Context) map[string]any {
			return map[string]any{"app": "my-service"}
		}),
		ctxlogecho.WithOnRequest(func(c echo.Context, fields map[string]any) {
			if b, err := json.Marshal(fields); err == nil {
				t.Log("====================================")
				t.Log("============ On Request ============")
				t.Log(string(b))
				t.Log("====================================")
			} else {
				t.Errorf("OnRequest fields=<marshal_error:%v>", err)
			}
		}),
		ctxlogecho.WithOnResponse(func(c echo.Context, fields map[string]any) {
			if b, err := json.Marshal(fields); err == nil {
				t.Log("====================================")
				t.Log("============ On Response ===========")
				t.Log(string(b))
				t.Log("====================================")
			} else {
				t.Errorf("OnResponse fields=<marshal_error:%v>", err)
			}
		}),
	))

	e.GET("/hello/:name", func(c echo.Context) error {
		name := c.Param("name")
		ctx := ctxlog.New(c.Request().Context())

		ctxlog.Add(ctx, "user_id", "u-123")

		ServiceExample(ctx, name)

		return c.JSON(http.StatusOK, map[string]any{"hello": name})
	})

	req := httptest.NewRequest(http.MethodGet, "/hello/bob", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
}

func ServiceExample(ctx context.Context, name string) {
	ctxlog.Add(ctx, "name", name)
	ctxlog.Add(ctx, "random-value", "some random")
}
