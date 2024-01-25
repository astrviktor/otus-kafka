package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"project/internal/config"
	"project/internal/informer"
	"project/internal/logger"
	"project/internal/metrics"
	"project/internal/middleware"
	"syscall"
)

func main() {
	if len(os.Args) > 1 {
		config.PrintUsage(config.InformerPrefix)
		return
	}

	log, err := logger.NewLogger("debug", true)
	if err != nil {
		fmt.Println("fail to create logger")
		os.Exit(1)
	}
	log = log.Named(config.InformerPrefix)

	cfg, err := config.ReadConfig(config.InformerPrefix)
	if err != nil {
		log.Error("fail to read config", zap.Error(err))
		os.Exit(1)
	}

	defer func() {
		_ = log.Sync()
		os.Exit(1)
	}()

	h, err := informer.NewHandler(log, cfg)
	if err != nil {
		log.Error("fail to create handler", zap.Error(err))
		os.Exit(1)
	}

	h.Run()

	r := router.New()
	r.GET("/info/{id}", middleware.Middleware(log, h.GetInfo))
	r.GET("/metrics", metrics.Metrics())

	server := fasthttp.Server{
		Handler:            r.Handler,
		MaxRequestBodySize: cfg.Receiver.MaxRequestBodySize,
	}

	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Informer.Host, cfg.Informer.Port)

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

	h.Stop()
}
