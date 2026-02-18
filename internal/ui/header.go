package ui

import "fmt"

const Version = "v0.1.0"

func RenderHeader(modelName string) string {
	return fmt.Sprintf("Cogito %s | %s", Version, modelName)
}
