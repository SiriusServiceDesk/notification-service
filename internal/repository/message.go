package repository

import (
	"fmt"
	"github.com/SiriusServiceDesk/notification-service/internal/config"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message *models.Message) error
	GetMessageByID(id string) (*models.Message, error)
	UpdateMessage(message *models.Message) error
	DeleteMessage(id string) error
	GetMessages() ([]*models.Message, error)
}

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func (m *MessageRepositoryImpl) CreateMessage(message *models.Message) error {
	result := m.db.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *MessageRepositoryImpl) GetMessageByID(id string) (*models.Message, error) {
	var message models.Message
	result := m.db.Where("id =?", id).First(&message)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

func (m *MessageRepositoryImpl) UpdateMessage(message *models.Message) error {
	result := m.db.Save(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *MessageRepositoryImpl) DeleteMessage(id string) error {
	result := m.db.Delete(&models.Message{}, "id =?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *MessageRepositoryImpl) GetMessages() ([]*models.Message, error) {
	var messages []*models.Message
	result := m.db.Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func NewMessageRepository() MessageRepository {
	cfg := config.GetConfig().Db
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	pgSvc := &MessageRepositoryImpl{db: db}
	err = db.AutoMigrate(&models.Message{})
	if err != nil {
		panic(err)
	}
	return pgSvc
}
