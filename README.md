# ctxlog

Small helper to attach structured key/value fields to a `context.Context`, plus adapters (starting with Echo v4) to initialize common request fields and emit them via callbacks.

## Status

- Go module path: `github.com/FaruqRahmadani/ctxlog` (see [go.mod](./go.mod)).

## Install

```bash
go get github.com/FaruqRahmadani/ctxlog
```

## Core usage (pure ctxlog)

API:

- `ctxlog.New(ctx)` ensures `ctx` has an internal map store (returns a context)
- `ctxlog.Add(ctx, key, value)` adds/overwrites a field
- `ctxlog.Copy(ctx)` returns a snapshot `map[string]any` of current fields

Example:

```go
package main

import (
	"context"
	"fmt"

	"github.com/FaruqRahmadani/ctxlog"
)

func main() {
	ctx := ctxlog.New(context.Background())

	ctxlog.Add(ctx, "user_id", "u-123")
	ctxlog.Add(ctx, "plan", "pro")

	fields := ctxlog.Copy(ctx)
	fmt.Printf("%v\n", fields)
}
```

## Echo v4 middleware

Package: [http/echo/v4](./http/echo/v4)

This middleware:

- ensures ctxlog is initialized on the request context
- populates common request fields
- optionally adds custom fields
- calls `OnRequest` before the handler and `OnResponse` after the handler

### Usage (functional options)

```go
package main

import (
	"log"
	"time"

	ctxlogecho "github.com/FaruqRahmadani/ctxlog/http/echo/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(ctxlogecho.New(
		ctxlogecho.WithAdditionalFields(func(c echo.Context) map[string]any {
			return map[string]any{"app": "my-service"}
		}),
		ctxlogecho.WithOnResponse(func(c echo.Context, status int, dur time.Duration, fields map[string]any) {
			log.Printf("status=%d dur=%s fields=%v", status, dur, fields)
		}),
	))

	log.Fatal(e.Start(":8080"))
}
```

### Default fields

The middleware sets these tags (constants live in [tags.go](./tags.go)):

- `request_id` (from `X-Request-ID`, generated if missing; also written back to response header)
- `http_method`
- `http_path`
- `http_host`
- `route` (Echo route pattern, `c.Path()`)
- `remote_ip` (`c.RealIP()`)
- `user_agent` (when present)

### Example output

Running the included example test:

```bash
go test -v ./...
```

Produces logs like:

```text
============ On Request ============
{"app":"my-service","http_host":"example.com","http_method":"GET","http_path":"/hello/bob","remote_ip":"192.0.2.1","request_id":"5231c780a39082ea4ccf4afda2ec417e","route":"/hello/:name"}
============ On Response ===========
{"app":"my-service","http_host":"example.com","http_method":"GET","http_path":"/hello/bob","name":"bob","random-value":"some random","remote_ip":"192.0.2.1","request_id":"5231c780a39082ea4ccf4afda2ec417e","route":"/hello/:name","user_id":"u-123"}
```

## Notes

- `ctxlog.Add` mutates an internal `map[string]any` stored inside the context. Treat a single request context as request-scoped; avoid concurrent writes from multiple goroutines unless you add your own synchronization.
- Use `ctxlog.Copy` when you want a stable snapshot for logging/serialization.

## Contributing

- Open an issue for bugs/feature requests.
- Send a PR with a focused change and tests (when applicable).

## License

MIT. See [LICENSE](./LICENSE).
