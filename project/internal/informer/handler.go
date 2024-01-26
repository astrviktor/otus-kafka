package informer

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"project/internal/config"
	"project/internal/model"
	"time"
)

type Handler struct {
	log *zap.Logger
	cfg *config.Config

	client *fasthttp.Client
}

func NewHandler(log *zap.Logger, cfg *config.Config) (*Handler, error) {
	client := &fasthttp.Client{}
	client.ReadTimeout = 30 * time.Second
	client.WriteTimeout = 30 * time.Second
	client.MaxConnsPerHost = 1024

	h := &Handler{
		log:    log,
		cfg:    cfg,
		client: client,
	}

	return h, nil
}

func (h *Handler) GetInfo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)

	ksql := fmt.Sprintf("SELECT * FROM jobs_stream WHERE id = '%s';", id)

	data, err := json.Marshal(model.KsqlData{
		Ksql:              ksql,
		StreamsProperties: model.StreamsProperties{},
	})

	if err != nil {
		h.log.Error("fail to marshal data", zap.Error(err))
		ctx.Error("fail to marshal data", http.StatusInternalServerError)
		return
	}

	requestURI := fmt.Sprintf("http://%s:%s/%s",
		h.cfg.Informer.RequestHost,
		h.cfg.Informer.RequestPort,
		h.cfg.Informer.RequestPath,
	)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(requestURI)
	req.Header.SetMethodBytes([]byte("POST"))
	req.Header.Add("content-type", "application/vnd.ksql.v1+json; charset=utf-8")
	req.SetBody(data)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = h.client.Do(req, resp)
	if err != nil {
		h.log.Error("fail to send ksql request", zap.Error(err))
		ctx.Error("fail to send ksql request", http.StatusInternalServerError)
		return
	}

	var ksqlResponse [3]interface{}
	err = json.Unmarshal(resp.Body(), &ksqlResponse)
	if err != nil {
		h.log.Error("fail to unmarshal data", zap.Error(err))
		ctx.Error("fail to unmarshal data", http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(ksqlResponse[1])
	if err != nil {
		h.log.Error("fail to marshal response data", zap.Error(err))
		ctx.Error("fail to marshal response data", http.StatusInternalServerError)
		return
	}

	informerResponse := model.InformerResponse{
		Id:     id,
		Status: string(model.JobStatusUnknown),
	}

	var response model.KsqlRow
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		h.log.Error("fail to unmarshal data", zap.Error(err))
		ctx.Error("fail to unmarshal data", http.StatusInternalServerError)
		return
	}

	if len(response.Row.Columns) == 4 {
		informerResponse.Id = response.Row.Columns[0]
		informerResponse.Status = response.Row.Columns[1]
		informerResponse.CreateDate = response.Row.Columns[2]
		informerResponse.FinishDate = response.Row.Columns[3]
	}

	body, err := json.Marshal(informerResponse)
	if err != nil {
		h.log.Error("fail to marshal response data", zap.Error(err))
		ctx.Error("fail to marshal response data", http.StatusInternalServerError)
		return
	}

	ctx.Success("application/json", body)
}

func (h *Handler) Run() {}

func (h *Handler) Stop() {}
