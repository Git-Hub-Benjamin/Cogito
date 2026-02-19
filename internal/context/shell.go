package context

import (
	"fmt"
	"os"
	"strings"
)

func BuildSystemMessage(includeCWD bool, customInstructions string) string {
	msg := `You are Cogito, a terminal assistant. Rules:
- Be direct — no filler, greetings, or unnecessary preamble
- Give complete, useful answers — include full code examples and explanations when the question warrants it
- For simple questions, keep it brief. For complex questions, give a thorough response
- Never repeat or echo the working directory back to the user`

	if includeCWD {
		if cwd, err := os.Getwd(); err == nil {
			msg += fmt.Sprintf("\n[Context: user is in %s — do NOT mention this unless they ask]", cwd)
		}
	}

	if strings.TrimSpace(customInstructions) != "" {
		msg += "\n\nUser's custom instructions:\n" + customInstructions
	}

	return msg
}
