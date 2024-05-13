package initializers

import (
	"github.com/SiriusServiceDesk/notification-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/notification-service/internal/web"
	"github.com/SiriusServiceDesk/notification-service/internal/web/status"
	"github.com/SiriusServiceDesk/notification-service/internal/web/swagger"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *dependencies.Container) {
	ctrl := buildRouters(container)

	for i := range ctrl {
		ctrl[i].DefineRouter(app)
	}
}

func buildRouters(container *dependencies.Container) []web.Controller {
	return []web.Controller{
		status.NewStatusController(),
		swagger.NewSwaggerController(),
	}
}
