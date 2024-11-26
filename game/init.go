package game

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time
type animationTickMsg time.Time
type easterEggTickMsg time.Time

const (
	cursorBlinkRate        time.Duration = time.Millisecond * 500
	defaultAnimationMPF    time.Duration = 30
	defaultAnimationFrames int           = 3
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickAnimation("start"),
		blinkTick(),
	)
}
