package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"kafka-golang/internal/domain"
	"kafka-golang/internal/infrastructure/log"
	"math/rand"
	"time"
)

func (uc UseCase) SendNotification(ctx context.Context) error {
	l := log.FromContext(ctx)
	typeId, _ := uuid.GenerateUUID()
	key := []byte(fmt.Sprintf("%s", typeId))

	l.Info("Create notification")

	newNotification := domain.Notification{
		TypeId: typeId,
		Text:   uc.generateRandomNotification(),
	}

	message, err := json.Marshal(newNotification)
	if err != nil {
		l.Error(err, "Failed to marshal Kafka message")
		return errors.Wrapf(err, "Failed to marshal Kafka message")
	}

	for i := 0; i < 10; i++ {
		err = uc.notificationProducer.Produce(ctx, key, message)
		if err != nil {
			l.Errorf("Failed to produce notification: %v", err)
			return errors.Wrapf(err, "Failed to produce notification")
		}
	}

	return nil
}

func (uc UseCase) generateRandomNotification() string {
	notifications := []string{
		"У вас новое сообщение!",
		"Не забудьте про предстоящее событие.",
		"Ваш заказ готов к получению.",
		"Вы получили бонусные баллы!",
		"Система обнаружила подозрительную активность в вашем аккаунте.",
		"Вам поступил новый запрос на добавление в друзья.",
		"Обновление успешно установлено.",
		"Напоминание: скоро истекает срок подписки.",
		"Ваш пароль был успешно изменён.",
		"Сегодня скидка 20% на все товары!",
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return notifications[rng.Intn(len(notifications))]
}
