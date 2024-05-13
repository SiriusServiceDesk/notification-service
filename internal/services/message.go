package services

import (
	"encoding/json"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"github.com/SiriusServiceDesk/notification-service/internal/queue"
	"github.com/SiriusServiceDesk/notification-service/internal/repository"
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	"github.com/google/uuid"
)

type MessageService interface {
	Create(message *models.Message) error
	GetByID(uuid string) (*models.Message, error)
	GetAll() ([]*models.Message, error)
	Update(message *models.Message) error
	Delete(uuid string) error
	Send(message *models.Message) error
}

type MessageRepositoryImpl struct {
	repos         repository.MessageRepository
	emails        EmailService
	templateRepos TemplateService
}

func NewMessageService(repos repository.MessageRepository, emails EmailService, template TemplateService) *MessageRepositoryImpl {
	return &MessageRepositoryImpl{repos: repos, emails: emails, templateRepos: template}
}

func (m *MessageRepositoryImpl) Create(message *models.Message) error {
	message.ID = uuid.New().String()
	return m.repos.CreateMessage(message)
}

func (m *MessageRepositoryImpl) GetByID(uuid string) (*models.Message, error) {
	return m.repos.GetMessageByID(uuid)
}

func (m *MessageRepositoryImpl) GetAll() ([]*models.Message, error) {
	return m.repos.GetMessages()
}

func (m *MessageRepositoryImpl) Update(message *models.Message) error {
	return m.repos.UpdateMessage(message)
}

func (m *MessageRepositoryImpl) Delete(uuid string) error {
	return m.repos.DeleteMessage(uuid)
}

func (m *MessageRepositoryImpl) Send(message *models.Message) error {
	switch message.Type {
	case "email":
		html, err := m.templateRepos.Render(message.TemplateName, message.Data)
		if err != nil {
			return err
		}
		var email = models.Email{
			From:    message.From,
			To:      message.To,
			Subject: message.Subject,
			Data:    html,
		}

		status, err := m.emails.SendEmail(&email)
		if err != nil {
			return err
		}

		message.Status = status
		if err = m.repos.UpdateMessage(message); err != nil {
			return err
		}

		return nil
	default:
		message.Status = "unknown message type"
		if err := m.repos.UpdateMessage(message); err != nil {
			return err
		}
		return nil
	}
}

func (m *MessageRepositoryImpl) ListenQueue(q queue.Queue) {
	for {
		queueMessages, err := q.Subscribe()
		if err != nil {
			panic(err)
		}

		for d := range queueMessages {
			var message models.Message
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				logger.Error("unmarshal message error", err)
				continue
			}

			err = m.Create(&message)
			if err != nil {
				logger.Error("error save message", err)
			}

			err = m.Send(&message)
			if err != nil {
				logger.Error("error send message", err)
			}
		}
	}
}
