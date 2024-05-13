package cli

import (
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	"github.com/spf13/cobra"
)

// ExecuteRootCmd prepares all CLI commands
func ExecuteRootCmd() {
	c := cobra.Command{}

	c.AddCommand(NewServeCmd())

	if err := c.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}
