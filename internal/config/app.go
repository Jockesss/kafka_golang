package config

import (
	"fmt"
	"kafka-golang/internal/config/modules"

	"github.com/ilyakaznacheev/cleanenv"
)

type KafkaAppConf struct {
	Env     string `env:"APP_ENV" env-required:"true"`
	AppName string `env:"APP_NAME" env-required:"true"`

	Log   modules.LoggerCfg `env-prefix:"APP_LOG_"`
	Kafka modules.Kafka     `env-prefix:"APP_KAFKA_"`
}

func LoadKafkaAppConfConf() (*KafkaAppConf, error) {
	var (
		cfg KafkaAppConf
		err error
	)

	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("load app config: %w", err)
	}

	return &cfg, nil
}
