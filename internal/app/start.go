package app

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"kafka-golang/internal/config"
	"kafka-golang/internal/infrastructure/kafka"
	"kafka-golang/internal/infrastructure/log"
	"os"
	"os/signal"
	"syscall"
)

const (
	appName = "kafka-golang"
)

var defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	l := log.New(defaultLevel, os.Stdout)
	ctx = log.WithLogger(ctx, l)

	cfg, err := config.LoadKafkaAppConfConf()
	if err != nil {
		return errors.Wrapf(err, "Load config file")
	}

	l.Infof("App running")

	testProducer, err := kafka.InitKafkaProducer(ctx, appName, cfg.Kafka.TestTopic, cfg)
	if err != nil {
		return errors.Wrap(err, "Init kafka producer")
	}
	defer testProducer.Close()

	<-ctx.Done()

	return nil
}
