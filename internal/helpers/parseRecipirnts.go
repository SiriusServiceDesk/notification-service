package helpers

import (
	"fmt"
	"strings"
)

func ParseRecipients(recipients []string) string {
	if len(recipients) != 0 {
		recipientsStr := strings.Join(recipients, ",")
		return fmt.Sprintf("Sent to everyone except %s", recipientsStr)
	}
	return "Sent"
}
