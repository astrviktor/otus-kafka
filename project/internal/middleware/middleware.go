package middleware

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func Middleware(log *zap.Logger, handle fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Debug("request",
			zap.String("path", string(ctx.Path())),
			zap.String("headers", string(ctx.Request.Header.Header())),
			zap.String("body", string(ctx.Request.Body())),
		)

		now := time.Now()
		handle(ctx)
		duration := time.Since(now)

		log.Debug("response",
			zap.Int("status", ctx.Response.StatusCode()),
			zap.String("headers", string(ctx.Response.Header.Header())),
			zap.String("body", string(ctx.Response.Body())),
		)

		metrics.GetOrCreateHistogram(fmt.Sprintf("request_duration{method=%q,path=%q,status=%q}",
			string(ctx.Method()), string(ctx.Path()), strconv.Itoa(ctx.Response.StatusCode())),
		).Update(duration.Seconds())
	}
}
