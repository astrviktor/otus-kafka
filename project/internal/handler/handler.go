package handler

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"project/internal/config"
)

type Handler struct {
	log      *zap.Logger
	cfg      *config.Config
	producer *kafka.Producer
}

func NewHandler(log *zap.Logger, cfg *config.Config) (*Handler, error) {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
	}

	producer, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		log.Error("fail to create kafka producer", zap.Error(err))
		return nil, err
	}

	h := &Handler{
		log: log,
		cfg: cfg,

		producer: producer,
	}

	return h, nil
}

func randomString(n int) string {
	length := n / 2
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (h *Handler) CreateJob(ctx *fasthttp.RequestCtx) {
	job := Job{
		Id:    uuid.New().String(),
		Sleep: int64(rand.Intn(h.cfg.Receiver.MaxSleep-h.cfg.Receiver.MinSleep) + h.cfg.Receiver.MinSleep),
		Data:  randomString(h.cfg.Receiver.DataSize),
	}

	value, err := json.Marshal(job)
	if err != nil {
		h.log.Error("fail to marshal data", zap.Error(err))
		ctx.Error("fail to marshal data", http.StatusInternalServerError)
		return
	}

	err = h.producer.Produce(&kafka.Message{
		Key: []byte(uuid.New().String()),
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.cfg.Kafka.Producer.Topic,
			Partition: kafka.PartitionAny},
		Value: value,
	}, nil)

	if err != nil {
		h.log.Error("fail to write data to kafka", zap.Error(err))
		ctx.Error("fail to write data to kafka", http.StatusInternalServerError)
		return
	} else {
		ctx.Success("application/json", nil)
	}
}

func (h *Handler) Run() {

}

func (h *Handler) Stop() {
	h.producer.Close()
}
