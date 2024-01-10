package handler

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/valyala/fasthttp"
)

func (h *Handler) Metrics() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		metrics.WritePrometheus(ctx.Response.BodyWriter(), true)
	}
}
