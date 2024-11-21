package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZachLTech/ansify"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	currentScene        string
	currentScenePrompt  string
	currentSceneGraphic []string
	userInput           string
	err                 error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			m.currentScene, m.currentScenePrompt, m.currentSceneGraphic = m.handleInput(m.userInput)
			m.userInput = ""
		case tea.KeyBackspace:
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
			}
		case tea.KeyRunes:
			if m.currentScene == "eXit" {
				m.currentScene = "dungeon"
				m.currentScenePrompt = "You're trapped in a dungeon with your friend.\nYou see a barrel. What do you do?"
				m.currentSceneGraphic = processANSIArt("./assets/" + m.currentScene + ".png")
			}

			m.userInput += string(msg.Runes)
		}
	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	for _, line := range m.currentSceneGraphic {
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n")

	sb.WriteString(m.currentScenePrompt + m.userInput)

	if m.err != nil {
		sb.WriteString(fmt.Sprintf("\nError: %v", m.err))
	}

	return sb.String()
}

func (m model) handleInput(userInput string) (string, string, []string) {
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	var scene string
	var prompt string
	var graphic []string

	if userInput == "move the barrel" || userInput == "move barrel" && m.currentScene == "dungeon" {
		scene = "secretTunnel"
		prompt = "The barrel rolls aside and you find a secret tunnel\nWhat do you do?"
		graphic = processANSIArt("./assets/" + scene + ".png")
	} else if userInput == "move the barrel" || userInput == "move barrel" && m.currentScene == "secretTunnel" {

	}

	return scene, prompt, graphic
}

func processANSIArt(imageInput string) []string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return []string{fmt.Sprintf("Error getting terminal size: %v", err)}
	}

	var lines []string
	var currentLine string
	count := 0

	ansiStr := ansify.GetAnsify(imageInput)

	for _, char := range ansiStr {
		currentLine += string(char)
		if char == '█' {
			count++
			if count == width {
				lines = append(lines, currentLine)
				currentLine = ""
				count = 0
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func main() {
	initialScene := processANSIArt("./assets/eXit.png")

	initialModel := model{
		currentScene:        "eXit",
		currentScenePrompt:  "PRESS ANY KEY TO START",
		currentSceneGraphic: initialScene,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
