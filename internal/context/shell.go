package context

import (
	"fmt"
	"os"
)

func BuildSystemMessage(includeCWD bool) string {
	msg := "You are a helpful terminal assistant. Be concise and direct."

	if includeCWD {
		if cwd, err := os.Getwd(); err == nil {
			msg += fmt.Sprintf("\nThe user's current working directory is: %s", cwd)
		}
	}

	return msg
}
