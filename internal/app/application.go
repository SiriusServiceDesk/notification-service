package app

import (
	"github.com/SiriusServiceDesk/notification-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/notification-service/internal/app/initializers"
	"github.com/SiriusServiceDesk/notification-service/internal/queue"
	"github.com/SiriusServiceDesk/notification-service/internal/repository"
	"github.com/SiriusServiceDesk/notification-service/internal/services"
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func InitApplication(app *fiber.App) {
	templateRepos := repository.NewTemplateRepository()
	templateService := services.NewTemplateService(templateRepos)

	emailService := services.NewEmailService()

	messageRepos := repository.NewMessageRepository()
	messageService := services.NewMessageService(messageRepos, emailService, templateService)

	q, err := queue.NewRabbitQueue()
	if err != nil {
		panic("failed to connect rabbitmq")
	}
	go messageService.ListenQueue(q)

	container := &dependencies.Container{
		TemplateService: templateService,
		MessageService:  messageService,
		Queue:           q,
	}

	grpcListener, err := initializers.InitializeGRPCListener()
	if err != nil {
		logger.Fatal("failed initializing grpc listener", zap.Error(err))
	}

	grpcServer := initializers.InitializeGRPCServer(grpcListener, container)

	initializers.SetupRoutes(app, container)

	initializers.StartGRPCServer(grpcServer)
}
