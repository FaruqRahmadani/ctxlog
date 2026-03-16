package ctxlogecho

import (
	"time"

	echo "github.com/labstack/echo/v4"
)

type Config struct {
	RequestIDHeader   string
	GenerateRequestID func() string
	DisableBaseFields bool
	AdditionalFields  func(c echo.Context) map[string]any
	OnRequest         func(c echo.Context, fields map[string]any)
	OnResponse        func(c echo.Context, status int, dur time.Duration, fields map[string]any)
}

type Option func(*Config)

func WithRequestIDHeader(v string) Option {
	return func(cfg *Config) { cfg.RequestIDHeader = v }
}

func WithGenerateRequestID(fn func() string) Option {
	return func(cfg *Config) { cfg.GenerateRequestID = fn }
}

func WithDisableBaseFields(v bool) Option {
	return func(cfg *Config) { cfg.DisableBaseFields = v }
}

func WithAdditionalFields(fn func(c echo.Context) map[string]any) Option {
	return func(cfg *Config) { cfg.AdditionalFields = fn }
}

func WithOnRequest(fn func(c echo.Context, fields map[string]any)) Option {
	return func(cfg *Config) { cfg.OnRequest = fn }
}

func WithOnResponse(fn func(c echo.Context, status int, dur time.Duration, fields map[string]any)) Option {
	return func(cfg *Config) { cfg.OnResponse = fn }
}
