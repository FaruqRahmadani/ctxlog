package main

import (
	"log"
	"net/http"

	"github.com/FaruqRahmadani/ctxlog"
	ctxlogecho "github.com/FaruqRahmadani/ctxlog/http/echo/v4"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(ctxlogecho.New(
		ctxlogecho.WithAdditionalFields(func(c echo.Context) map[string]any {
			return map[string]any{"app": "my-service"}
		}),

		ctxlogecho.WithOnRequest(func(c echo.Context, fields map[string]any) {
			log.Printf("IN  method=%s route=%s fields=%v",
				c.Request().Method,
				c.Path(),
				fields,
			)
		}),

		ctxlogecho.WithOnResponse(func(c echo.Context, fields map[string]any) {
			log.Printf("OUT status=%d dur=%s fields=%v",
				c.Response().Status,
				fields[ctxlog.TagLatencyHuman],
				fields,
			)
		}),
	))

	e.GET("/hello/:name", func(c echo.Context) error {
		ctxlog.Add(c.Request().Context(), "user_id", "u-123")
		return c.JSON(http.StatusOK, map[string]any{
			"hello": c.Param("name"),
		})
	})

	log.Fatal(e.Start(":8080"))
}
