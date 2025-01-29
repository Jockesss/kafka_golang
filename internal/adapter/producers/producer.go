package producer

import (
	"context"
	"kafka-golang/internal/infrastructure/log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const messageIDHeader = "MessageId"

type Producer interface {
	Produce(ctx context.Context, key []byte, message []byte) error
	ProduceWithMsgID(ctx context.Context, msgID, key, message []byte) error
	Close(ctx context.Context)
	ProduceWithHeaders(ctx context.Context, key, message []byte, headers map[string]string) error
}

type eventProducer struct {
	topicName string
	producer  sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer, topicName string) (Producer, error) {
	return &eventProducer{
		topicName: topicName,
		producer:  producer,
	}, nil
}

func (p *eventProducer) ProduceWithMsgID(ctx context.Context, msgID, key, message []byte) error {
	return p.produce(ctx, msgID, key, message, nil)
}

func (p *eventProducer) ProduceWithHeaders(ctx context.Context, key, message []byte, headers map[string]string) error {
	uuidP, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Generate uuid")
	}

	return p.produce(ctx, []byte(uuidP.String()), key, message, headers)
}

func (p *eventProducer) Produce(ctx context.Context, key, message []byte) error {
	uuidP, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Generate uuid")
	}

	return p.produce(ctx, []byte(uuidP.String()), key, message, nil)
}

func (p *eventProducer) produce(ctx context.Context, msgID, key, message []byte, headers map[string]string) error {
	saramaHeaders := []sarama.RecordHeader{
		{
			Key:   []byte(messageIDHeader),
			Value: msgID,
		},
	}

	for k, v := range headers {
		saramaHeaders = append(saramaHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic:   p.topicName,
		Key:     sarama.ByteEncoder(key),
		Value:   sarama.ByteEncoder(message),
		Headers: saramaHeaders,
	}

	l := log.FromContext(ctx)

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		l.Errorf("Failed to send message to Kafka: %v", err)
		return errors.Wrapf(err, "Failed to produce message to topic %s", p.topicName)
	}

	l.Infof("Message sent successfully to topic %s partition %d offset %d", p.topicName, partition, offset)
	return nil
}

func (p *eventProducer) Close(ctx context.Context) {
	if err := p.producer.Close(); err != nil {
		log.FromContext(ctx).Errorf("Error closing Kafka producer: %v\n", err)
	} else {
		log.FromContext(ctx).Infof("Kafka producer stopped")
	}
}
