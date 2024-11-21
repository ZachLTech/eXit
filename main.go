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
				m.userInput = ""
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

	if userInput == "move the barrel" || userInput == "move barrel" && m.currentScene == "dungeon" {
		scene = "secretTunnel"
		prompt = "The barrel rolls aside and you find a secret tunnel\nWhat do you do?\n\n>"
	} else if userInput == "enter the tunnel" || userInput == "enter tunnel" && m.currentScene == "secretTunnel" {
		scene = "friendTooWeak"
		prompt = "You start to escape but your friend is too weak to\ngo with you. They hand you a note.\nWhat do you do?\n\n>"
	} else if userInput == "read the note" || userInput == "read note" && m.currentScene == "friendTooWeak" {
		scene = "friendHandsNote"
		prompt = "It is too dark to read the note.\nWhat do you do?\n\n>"
	} else if userInput == "leave" && m.currentScene == "friendHandsNote" {
		scene = "beach"
		prompt = "You crawl through the tunnel and the tunnel leads\nyou to a beach. What do you do?\n\n>"
	} else if userInput == "look" || userInput == "look around" && m.currentScene == "beach" {
		scene = "ship"
		prompt = "In the water you see a boat.\nWhat do you do?\n\n>"
	} else if userInput == "get on the boat" || userInput == "get on boat" && m.currentScene == "ship" {
		scene = "congratulations"
		prompt = "Congratulations, you're heading to a new world!\nDo you want to play again?\n\n>"
	} else if userInput == "yes" && m.currentScene == "congratulations" {
		scene = "dungeon"
		prompt = "You're trapped in a dungeon with your friend. You see a barrel. What do you do?\n\n>"
	}

	graphic := processANSIArt("./assets/" + scene + ".png")

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
