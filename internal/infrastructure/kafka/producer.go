package kafka

import (
	"context"
	"fmt"
	producer "kafka-golang/internal/adapter/producers"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"kafka-golang/internal/config"
	"kafka-golang/internal/infrastructure/log"
)

func InitKafkaProducer(ctx context.Context, appName, topic string, cfg *config.KafkaAppConf) (producer.Producer, error) {
	uuidForClient, _ := uuid.NewUUID()
	clientIDStr := fmt.Sprintf("%s-%s", appName, uuidForClient)

	l := log.FromContext(ctx)

	cfgMap := sarama.NewConfig()
	cfgMap.ClientID = clientIDStr
	cfgMap.Producer.Return.Successes = true
	cfgMap.Producer.Return.Errors = true
	cfgMap.Producer.RequiredAcks = sarama.WaitForAll

	syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, cfgMap)
	if err != nil {
		l.Errorf(err.Error(), "Failed to create Kafka producer")
		return nil, errors.Wrap(err, "Failed to create producer")
	}

	kafkaProducer, err := producer.NewKafkaProducer(syncProducer, topic)
	if err != nil {
		l.Errorf(err.Error(), "Failed to create Kafka producer")
		return nil, errors.Wrap(err, "Failed to create producer")
	}

	l.Infof("Kafka producer created")

	return kafkaProducer, nil
}
