package app

type AppState int

const (
	StateInput    AppState = iota
	StateStreaming
	StateSettings
)
