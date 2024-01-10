package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

const ReceiverPrefix = "receiver"
const ProcessorPrefix = "processor"

type Config struct {
	Receiver  ReceiverConfig
	Processor ProcessorConfig
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

type LoggerConfig struct {
	Level string `default:"debug"`
	Debug bool   `default:"false"`
}

type KafkaConfig struct {
	BootstrapServers string `default:"localhost:9092"`

	Producer struct {
		Topic string `default:"kafka-topic"`
	}

	Consumer struct {
		GroupID string `default:"processor"`
		Topic   string `default:"kafka-topic"`
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
