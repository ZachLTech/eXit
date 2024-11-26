package game

import (
	"fmt"
	"strings"
)

// I love how simple this is :)
func (m model) View() string {
	var sb strings.Builder

	for _, line := range m.currentSceneGraphic {
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n")

	if !m.animating {
		sb.WriteString(m.currentScenePrompt + m.userInput)
	}

	if m.cursorBlink {
		sb.WriteString(m.cursorSymbol)
	} else {
		sb.WriteString(" ")
	}

	if m.err != nil {
		sb.WriteString(fmt.Sprintf("\nError: %v", m.err))
	}

	return sb.String()
}
