package dependencies

import (
	"github.com/SiriusServiceDesk/notification-service/internal/queue"
	"github.com/SiriusServiceDesk/notification-service/internal/services"
)

type Container struct {
	TemplateService services.TemplateService
	MessageService  services.MessageService
	Queue           queue.Queue
}
