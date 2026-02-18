package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type streamChunkMsg struct {
	chunk string
}

type streamDoneMsg struct{}

type streamErrMsg struct {
	err error
}

type configSavedMsg struct{}

type configSaveErrMsg struct {
	err error
}

// listenForChunks reads one chunk from the channel and returns the appropriate message.
func listenForChunks(ch <-chan string, errCh <-chan error) tea.Cmd {
	return func() tea.Msg {
		select {
		case chunk, ok := <-ch:
			if !ok {
				select {
				case err := <-errCh:
					if err != nil {
						return streamErrMsg{err: err}
					}
				default:
				}
				return streamDoneMsg{}
			}
			return streamChunkMsg{chunk: chunk}
		case err := <-errCh:
			if err != nil {
				return streamErrMsg{err: err}
			}
			// Drain remaining chunks
			for chunk := range ch {
				_ = chunk
			}
			return streamDoneMsg{}
		}
	}
}
