package ctxlogecho

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/FaruqRahmadani/ctxlog"

	echo "github.com/labstack/echo/v4"
)

func New(opts ...Option) echo.MiddlewareFunc {
	var cfg Config
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return middleware(cfg)
}

func middleware(cfg Config) echo.MiddlewareFunc {
	if cfg.RequestIDHeader == "" {
		cfg.RequestIDHeader = "X-Request-ID"
	}
	if cfg.GenerateRequestID == nil {
		cfg.GenerateRequestID = defaultRequestID
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			ctx := ctxlog.New(req.Context())

			requestID := req.Header.Get(cfg.RequestIDHeader)
			if requestID == "" && cfg.GenerateRequestID != nil {
				requestID = cfg.GenerateRequestID()
			}
			if requestID != "" {
				ctxlog.Add(ctx, ctxlog.TagRequestID, requestID)
				if cfg.RequestIDHeader != "" && c.Response().Header().Get(cfg.RequestIDHeader) == "" {
					c.Response().Header().Set(cfg.RequestIDHeader, requestID)
				}
			}

			if !cfg.DisableBaseFields {
				ctxlog.Add(ctx, ctxlog.TagHTTPMethod, req.Method)
				ctxlog.Add(ctx, ctxlog.TagHTTPPath, req.URL.Path)
				ctxlog.Add(ctx, ctxlog.TagHTTPHost, req.Host)
				if route := c.Path(); route != "" {
					ctxlog.Add(ctx, ctxlog.TagRoute, route)
				}
				if ip := c.RealIP(); ip != "" {
					ctxlog.Add(ctx, ctxlog.TagRemoteIP, ip)
				}
				if ua := req.UserAgent(); ua != "" {
					ctxlog.Add(ctx, ctxlog.TagUserAgent, ua)
				}
			}

			if cfg.AdditionalFields != nil {
				for k, v := range cfg.AdditionalFields(c) {
					ctxlog.Add(ctx, k, v)
				}
			}

			c.SetRequest(req.WithContext(ctx))

			if cfg.OnRequest != nil {
				cfg.OnRequest(c, ctxlog.Get(ctx))
			}

			err := next(c)

			if cfg.OnResponse != nil {
				status := c.Response().Status
				if status == 0 {
					status = http.StatusOK
				}
				cfg.OnResponse(c, status, time.Since(start), ctxlog.Get(ctx))
			}

			return err
		}
	}
}

var fallbackReqIDCounter uint64

func defaultRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}
	var fb [16]byte
	binary.BigEndian.PutUint64(fb[:8], uint64(time.Now().UnixNano()))
	binary.BigEndian.PutUint64(fb[8:], atomic.AddUint64(&fallbackReqIDCounter, 1))
	return hex.EncodeToString(fb[:])
}
