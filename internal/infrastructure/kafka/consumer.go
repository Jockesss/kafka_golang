package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"kafka-golang/internal/config"
	"kafka-golang/internal/infrastructure/log"
)

func InitKafkaConsumer(ctx context.Context, appName string, cfg *config.KafkaAppConf) (sarama.ConsumerGroup, error) {
	uuidForClient, _ := uuid.NewUUID()
	clientIDStr := fmt.Sprintf("%s-%s", appName, uuidForClient)

	l := log.FromContext(ctx)

	l.Infof("Generate UUID for client consumer")

	cfgMap := sarama.NewConfig()
	cfgMap.ClientID = clientIDStr
	cfgMap.Consumer.Group.InstanceId = uuidForClient.String()
	cfgMap.Consumer.Offsets.AutoCommit.Enable = false
	cfgMap.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, clientIDStr, cfgMap)
	if err != nil {
		l.Errorf(err.Error(), "Failed to create Kafka consumer")
		return nil, errors.Wrap(err, "Create Kafka consumer with sarama")
	}

	l.Infof("Kafka consumer successfully created with ID '%s'", clientIDStr)

	return consumer, nil
}
