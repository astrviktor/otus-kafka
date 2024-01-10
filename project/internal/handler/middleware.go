package handler

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func (h *Handler) Middleware(handle fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h.log.Debug("request", zap.String("path", string(ctx.Path())))

		now := time.Now()
		handle(ctx)
		duration := time.Since(now)

		h.log.Debug("response", zap.Int("status", ctx.Response.StatusCode()))

		metrics.GetOrCreateHistogram(fmt.Sprintf("request_duration{method=%q,path=%q,status=%q}",
			string(ctx.Method()), string(ctx.Path()), strconv.Itoa(ctx.Response.StatusCode())),
		).Update(duration.Seconds())
	}
}
