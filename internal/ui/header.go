package ui

import "fmt"

const Version = "v0.1.0"

func RenderHeader(modelName string, lastQuery string, maxWidth int) string {
	prefix := fmt.Sprintf("Cogito %s | %s", Version, modelName)

	if lastQuery == "" {
		return prefix
	}

	// Calculate available space for query text
	// prefix + " | " + query + potential "..."
	separator := " | "
	available := maxWidth - len(prefix) - len(separator) - 4 // 4 for border chars
	if available < 5 {
		return prefix
	}

	if len(lastQuery) > available {
		lastQuery = lastQuery[:available-3] + "..."
	}

	return prefix + separator + lastQuery
}
