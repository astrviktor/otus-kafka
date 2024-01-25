package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

const ReceiverPrefix = "receiver"
const ProcessorPrefix = "processor"
const ConverterPrefix = "converter"
const InformerPrefix = "informer"

type Config struct {
	Receiver  ReceiverConfig
	Processor ProcessorConfig
	Converter ConverterConfig
	Informer  InformerConfig
	Logger    LoggerConfig
	Kafka     KafkaConfig
}

type ReceiverConfig struct {
	Host string `default:"127.0.0.1"`
	Port string `default:"8081"`

	MaxRequestBodySize int `default:"5242880"`
	MinSleep           int `default:"100"`
	MaxSleep           int `default:"200"`
	DataSize           int `default:"1024"`
}

type ProcessorConfig struct {
	Host string `default:"127.0.0.1"`
	Port string `default:"8082"`
}

type ConverterConfig struct {
	Host string `default:"127.0.0.1"`
	Port string `default:"8083"`
}

type InformerConfig struct {
	Host string `default:"127.0.0.1"`
	Port string `default:"8084"`

	RequestHost string `default:"127.0.0.1"`
	RequestPort string `default:"8088"`
	RequestPath string `default:"query"`
}

type LoggerConfig struct {
	Level string `default:"debug"`
	Debug bool   `default:"false"`
}

type KafkaConfig struct {
	BootstrapServers string `default:"localhost:9091"`

	ReceiverProducer struct {
		Topic string `default:"receiver-topic"`
	}

	ProcessorConsumer struct {
		GroupID     string        `default:"processor"`
		Topic       string        `default:"receiver-topic"`
		ReadTimeout time.Duration `default:"10s"`
	}

	ProcessorProducer struct {
		Topic string `default:"processor-topic"`
	}

	ConverterConsumer struct {
		GroupID     string        `default:"converter"`
		Topic       string        `default:"processor-topic"`
		ReadTimeout time.Duration `default:"10s"`
	}

	ConverterPostgresProducer struct {
		Topic string `default:"postgres-topic"`
	}

	ConverterClickhouseProducer struct {
		Topic string `default:"clickhouse-topic"`
	}
}

func ReadConfig(prefix string) (*Config, error) {
	cfg := Config{}

	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config: %s\n", err.Error())
	}

	return &cfg, nil
}

func PrintUsage(prefix string) {
	cfg := Config{}
	err := envconfig.Usage(prefix, &cfg)
	if err != nil {
		fmt.Printf("fail to print envconfig usage: %s", err.Error())
	}
}
