package usecase

import (
	p "kafka-golang/internal/adapter/producers"
)

type UseCase struct {
	notificationProducer p.Producer
}

func NewUseCase(notificationProducer p.Producer) *UseCase {
	return &UseCase{notificationProducer: notificationProducer}
}
