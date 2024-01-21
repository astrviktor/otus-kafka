package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
	"project/internal/config"
	"project/internal/model"
	"sync"
	"time"
)

type Worker struct {
	log *zap.Logger
	cfg *config.Config

	consumer *kafka.Consumer

	producerPostgres   *kafka.Producer
	producerClickhouse *kafka.Producer

	ctx  context.Context
	done context.CancelFunc
	wg   sync.WaitGroup
}

func NewWorker(
	log *zap.Logger,
	cfg *config.Config,
) (*Worker, error) {
	kafkaConsumerConfig := kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
		"group.id":          cfg.Kafka.ConverterConsumer.GroupID,
	}

	consumer, err := kafka.NewConsumer(&kafkaConsumerConfig)
	if err != nil {
		log.Error("fail to create kafka consumer", zap.Error(err))
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{cfg.Kafka.ConverterConsumer.Topic}, nil)
	if err != nil {
		log.Error("subscribe topics error", zap.Error(err))
		return nil, err
	}

	kafkaProducerConfig := kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
	}

	producerPostgres, err := kafka.NewProducer(&kafkaProducerConfig)
	if err != nil {
		log.Error("fail to create kafka producer postgres", zap.Error(err))
		return nil, err
	}

	producerClickhouse, err := kafka.NewProducer(&kafkaProducerConfig)
	if err != nil {
		log.Error("fail to create kafka producer clickhouse", zap.Error(err))
		return nil, err
	}

	w := &Worker{
		log: log,
		cfg: cfg,

		consumer: consumer,

		producerPostgres:   producerPostgres,
		producerClickhouse: producerClickhouse,
	}

	w.ctx, w.done = context.WithCancel(context.Background())

	return w, nil
}

func (w *Worker) worker() {
	for {
		select {
		case <-w.ctx.Done():
			w.wg.Done()
			return
		default:
			w.middleware(w.processing)()
		}
	}
}

type FuncHandler func()
type FuncHandlerError func() (string, string)

const (
	resultEmpty   = "empty"
	resultDone    = "done"
	statusError   = "error"
	statusSuccess = "success"
)

func (w *Worker) middleware(handle FuncHandlerError) FuncHandler {
	return func() {
		now := time.Now()

		result, status := handle()

		if result == resultEmpty && status == statusSuccess {
			return
		}

		duration := time.Since(now)
		metrics.GetOrCreateHistogram(fmt.Sprintf("processor_duration{result=%q, status=%q}", result, status)).Update(duration.Seconds())
	}
}

func (w *Worker) processing() (string, string) {
	msg, err := w.consumer.ReadMessage(w.cfg.Kafka.ConverterConsumer.ReadTimeout)
	if err != nil {
		if err.(kafka.Error).IsTimeout() {
			w.log.Info("consumer read timeout", zap.Error(err))
			return resultEmpty, statusSuccess
		}
		w.log.Error("consumer error", zap.Error(err))
		return resultEmpty, statusError
	}

	var jobDone model.JobDone
	err = json.Unmarshal(msg.Value, &jobDone)
	if err != nil {
		w.log.Error("fail to unmarshal data", zap.Error(err))
		return resultEmpty, statusError
	}
	w.log.Info("jobDone received", zap.Reflect("jobDone", jobDone))

	// to postgres
	w.log.Info("jobDone to postgres message")

	valuePostgres, err := w.JobDoneToPostgresMessage(jobDone)
	if err != nil {
		w.log.Error("fail to convert to postgres message", zap.Error(err))
		return resultEmpty, statusError
	}

	err = w.producerPostgres.Produce(&kafka.Message{
		Key: []byte(jobDone.Id),
		TopicPartition: kafka.TopicPartition{
			Topic:     &w.cfg.Kafka.ConverterPostgresProducer.Topic,
			Partition: kafka.PartitionAny},
		Value: valuePostgres,
	}, nil)

	if err != nil {
		w.log.Error("fail to write data to kafka", zap.Error(err))
		return resultEmpty, statusError
	}

	// to clickhouse
	w.log.Info("jobDone to clickhouse message")

	valueClickhouse, err := w.JobDoneToPostgresMessage(jobDone)
	if err != nil {
		w.log.Error("fail to convert to clickhouse message", zap.Error(err))
		return resultEmpty, statusError
	}

	err = w.producerClickhouse.Produce(&kafka.Message{
		Key: []byte(jobDone.Id),
		TopicPartition: kafka.TopicPartition{
			Topic:     &w.cfg.Kafka.ConverterClickhouseProducer.Topic,
			Partition: kafka.PartitionAny},
		Value: valueClickhouse,
	}, nil)

	if err != nil {
		w.log.Error("fail to write data to kafka", zap.Error(err))
		return resultEmpty, statusError
	}

	return resultDone, statusSuccess
}

func (w *Worker) Run() {
	w.wg.Add(1)
	go w.worker()

	w.log.Info("worker running")
}

func (w *Worker) Stop() {
	w.done()
	w.wg.Wait()

	err := w.consumer.Close()
	if err != nil {
		w.log.Error("fail to close consumer", zap.Error(err))
	}

	w.producerPostgres.Close()
	w.producerClickhouse.Close()
	w.log.Info("worker was stopped")
}

func (w *Worker) JobDoneToPostgresMessage(jobDone model.JobDone) ([]byte, error) {
	msg := model.PostgresMessage{
		Id:         jobDone.Id,
		Status:     jobDone.Status,
		DateCreate: jobDone.DateCreate,
		DateDone:   jobDone.DateDone,
	}

	return json.Marshal(msg)
}

func (w *Worker) JobDoneToClickhouseMessage(jobDone model.JobDone) ([]byte, error) {
	msg := model.PostgresMessage{
		Id:         jobDone.Id,
		Status:     jobDone.Status,
		DateCreate: jobDone.DateCreate,
		DateDone:   jobDone.DateDone,
	}

	return json.Marshal(msg)
}
