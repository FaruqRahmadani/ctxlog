package ctxlogecho

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"strings"
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
				ctxlog.Add(ctx, ctxlog.TagQueryParams, flattenQueryParams(req.URL.Query()))
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

			ctxlog.Add(ctx, ctxlog.TagStatus, c.Response().Status)
			ctxlog.Add(ctx, ctxlog.TagDuration, time.Since(start))
			ctxlog.Add(ctx, ctxlog.TagLatencyMS, time.Since(start).Milliseconds())
			ctxlog.Add(ctx, ctxlog.TagLatencyHuman, time.Since(start).String())

			if cfg.OnResponse != nil {
				cfg.OnResponse(c, ctxlog.Get(ctx))
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

func flattenQueryParams(values map[string][]string) map[string]string {
	if len(values) == 0 {
		return nil
	}

	out := make(map[string]string, len(values))
	for k, vals := range values {
		if len(vals) == 0 {
			out[k] = ""
			continue
		}
		if len(vals) == 1 {
			out[k] = vals[0]
			continue
		}
		out[k] = strings.Join(vals, ",")
	}

	return out
}
