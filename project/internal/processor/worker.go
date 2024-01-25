package processor

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
	producer *kafka.Producer

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
		"group.id":          cfg.Kafka.ProcessorConsumer.GroupID,
	}

	consumer, err := kafka.NewConsumer(&kafkaConsumerConfig)
	if err != nil {
		log.Error("fail to create kafka consumer", zap.Error(err))
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{cfg.Kafka.ProcessorConsumer.Topic}, nil)
	if err != nil {
		log.Error("subscribe topics error", zap.Error(err))
		return nil, err
	}

	kafkaProducerConfig := kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
	}

	producer, err := kafka.NewProducer(&kafkaProducerConfig)
	if err != nil {
		log.Error("fail to create kafka producer", zap.Error(err))
		return nil, err
	}

	w := &Worker{
		log: log,
		cfg: cfg,

		consumer: consumer,
		producer: producer,
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
	msg, err := w.consumer.ReadMessage(w.cfg.Kafka.ProcessorConsumer.ReadTimeout)
	if err != nil {
		if err.(kafka.Error).IsTimeout() {
			w.log.Info("consumer read timeout", zap.Error(err))
			return resultEmpty, statusSuccess
		}
		w.log.Error("consumer error", zap.Error(err))
		return resultEmpty, statusError
	}

	var job model.Job
	err = json.Unmarshal(msg.Value, &job)
	if err != nil {
		w.log.Error("fail to unmarshal data", zap.Error(err))
		return resultEmpty, statusError
	}

	w.log.Info("job received", zap.Reflect("job", job))

	// processing job
	time.Sleep(time.Millisecond * time.Duration(job.Sleep))

	w.log.Info("job processing done", zap.Reflect("job", job))

	jobDone := model.JobDone{
		Id:         job.Id,
		Status:     model.JobStatusFinish,
		CreateDate: job.CreateDate,
		FinishDate: time.Now(),
	}

	value, err := json.Marshal(jobDone)
	if err != nil {
		w.log.Error("fail to marshal data", zap.Error(err))
		return resultEmpty, statusError
	}

	err = w.producer.Produce(&kafka.Message{
		Key: []byte(jobDone.Id),
		TopicPartition: kafka.TopicPartition{
			Topic:     &w.cfg.Kafka.ProcessorProducer.Topic,
			Partition: kafka.PartitionAny},
		Value: value,
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

	w.producer.Close()
	w.log.Info("worker was stopped")
}
