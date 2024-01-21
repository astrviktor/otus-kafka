package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"project/internal/config"
	"project/internal/converter"
	"project/internal/logger"
	"project/internal/metrics"
	"syscall"
)

func main() {
	if len(os.Args) > 1 {
		config.PrintUsage(config.ConverterPrefix)
		return
	}

	log, err := logger.NewLogger("debug", true)
	if err != nil {
		fmt.Println("fail to create logger")
		os.Exit(1)
	}
	log = log.Named(config.ConverterPrefix)

	cfg, err := config.ReadConfig(config.ConverterPrefix)
	if err != nil {
		log.Error("fail to read config", zap.Error(err))
		os.Exit(1)
	}

	defer func() {
		_ = log.Sync()
		os.Exit(1)
	}()

	w, err := converter.NewWorker(log, cfg)
	if err != nil {
		log.Error("fail to create worker", zap.Error(err))
		os.Exit(1)
	}

	w.Run()

	r := router.New()
	r.GET("/metrics", metrics.Metrics())

	server := fasthttp.Server{
		Handler: r.Handler,
	}

	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Converter.Host, cfg.Converter.Port)

		log.Info("start server", zap.String("server addr", addr))
		err = server.ListenAndServe(addr)
		if err != nil {
			log.Error("fail to listen", zap.Error(err))
			os.Exit(1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	log.Info("caught os signal to stop server")

	err = server.Shutdown()
	if err != nil {
		log.Error("fail to shutdown service", zap.Error(err))
	} else {
		log.Info("server was successfully stopped")
	}

	w.Stop()
}
