package app

import (
	"context"
	"fmt"
	"kafka-golang/internal/config"
	"kafka-golang/internal/infrastructure/log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap/zapcore"
)

const (
	consumerGroup = "rnis-service-kafka-app"
)

//nolint:funlen,gocritic,nolintlint
func Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadKafkaAppConfConf()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logLevel := zapcore.ErrorLevel

	if err = logLevel.UnmarshalText([]byte(cfg.Log.Level)); err != nil {
		log.Logger().Errorf("parse log level env var: %v", err)
	}

	log.SetLevel(logLevel)

	log.Logger().Info("app running")

	<-ctx.Done()

	return nil
}
