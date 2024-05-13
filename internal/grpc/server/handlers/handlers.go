package handlers

import (
	"context"
	"encoding/json"
	"github.com/SiriusServiceDesk/gateway-service/pkg/notification_v1"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"github.com/SiriusServiceDesk/notification-service/internal/queue"
	"github.com/SiriusServiceDesk/notification-service/internal/services"
	"net/http"
)

type Handlers struct {
	notification_v1.UnimplementedNotificationV1Server
	Message services.MessageService
	Q       queue.Queue
}

func NewHandler(message services.MessageService, q queue.Queue) *Handlers {
	return &Handlers{
		Message: message,
		Q:       q,
	}
}

func (h Handlers) CreateMessage(
	ctx context.Context,
	request *notification_v1.CreateMessageRequest) (
	*notification_v1.CreateMessageResponse, error) {

	messageData := request
	var data models.JSONB

	err := json.Unmarshal([]byte(messageData.GetData()), &data)
	if err != nil {
		return nil, err
	}

	message := models.Message{
		Subject:      messageData.GetSubject(),
		To:           messageData.GetTo(),
		Data:         data,
		Type:         messageData.GetType(),
		TemplateName: messageData.GetTemplateName(),
	}

	err = h.Message.Create(&message)
	if err != nil {
		return nil, err
	}

	err = h.Q.Publish(message)
	if err != nil {
		return nil, err
	}

	return &notification_v1.CreateMessageResponse{
		Status:  http.StatusOK,
		Message: "message send successfully",
	}, nil
}
