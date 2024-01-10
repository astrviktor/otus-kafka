package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string, development bool) (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	loggerLevel := zapcore.DebugLevel
	if level != "" {
		if err := loggerLevel.Set(level); err != nil {
			fmt.Printf("wrong logger level[%s]", level)
			return nil, err
		}
	}

	levelEnablerFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= loggerLevel
	})

	ws, _, err := zap.Open("stdout")
	if err != nil {
		return nil, err
	}

	var logger *zap.Logger
	if development {
		logger = zap.New(zapcore.NewCore(encoder, ws, levelEnablerFunc), zap.Development())
	} else {
		logger = zap.New(zapcore.NewCore(encoder, ws, levelEnablerFunc))
	}

	zap.ReplaceGlobals(logger)

	return logger, nil
}
