package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"os/exec"
	"runtime"

	"github.com/ZachLTech/ansify"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentScene        string
	currentScenePrompt  string
	cursorBlink         bool
	currentSceneGraphic []string
	userInput           string
	err                 error
}

type tickMsg time.Time

const (
	blinkRate    = time.Millisecond * 500
	cursorSymbol = "░"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		blinkTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error

	switch msg := msg.(type) {
	case tickMsg:
		m.cursorBlink = !m.cursorBlink
		return m, blinkTick()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.currentScene == "eXit" {
				m.userInput = ""
				m.currentScene = "dungeon"
				m.currentScenePrompt = "You're trapped in a dungeon with your friend.\nYou see a barrel. What do you do?\n\n> "
				m.currentSceneGraphic, err = processANSIArt("./assets/" + m.currentScene + ".png")
				if err != nil {
					fmt.Printf("Error processing ANSI art: %v\n", err)
					os.Exit(1)
				}

				return m, nil
			}
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
				m.currentScenePrompt = "You're trapped in a dungeon with your friend.\nYou see a barrel. What do you do?\n\n> "
				m.currentSceneGraphic, err = processANSIArt("./assets/" + m.currentScene + ".png")
				if err != nil {
					fmt.Printf("Error processing ANSI art: %v\n", err)
					os.Exit(1)
				}

				return m, nil
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

	if m.cursorBlink {
		sb.WriteString(cursorSymbol)
	} else {
		sb.WriteString(" ")
	}

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
		prompt = "The barrel rolls aside and you find a secret tunnel\nWhat do you do?\n\n> "
	} else if userInput == "enter the tunnel" || userInput == "enter tunnel" && m.currentScene == "secretTunnel" {
		scene = "friendTooWeak"
		prompt = "You start to escape but your friend is too weak to\ngo with you. They hand you a note.\nWhat do you do?\n\n> "
	} else if userInput == "read the note" || userInput == "read note" && m.currentScene == "friendTooWeak" {
		scene = "friendHandsNote"
		prompt = "It is too dark to read the note.\nWhat do you do?\n\n> "
	} else if userInput == "leave" && (m.currentScene == "friendHandsNote" || m.currentScene == "friendTooWeak") {
		scene = "beach"
		prompt = "You crawl through the tunnel and the tunnel leads\nyou to a beach. What do you do?\n\n> "
	} else if userInput == "look" || userInput == "look around" && m.currentScene == "beach" {
		scene = "ship"
		prompt = "In the water you see a boat.\nWhat do you do?\n\n> "
	} else if userInput == "get on the boat" || userInput == "get on boat" || userInput == "get on" && m.currentScene == "ship" {
		scene = "congratulations"
		prompt = "Congratulations, you're heading to a new world!\nDo you want to play again?\n\n> "
	} else if userInput == "yes" && m.currentScene == "congratulations" {
		scene = "dungeon"
		prompt = "You're trapped in a dungeon with your friend. You see a barrel. What do you do?\n\n> "
	} else if userInput == "no" && m.currentScene == "congratulations" {
		os.Exit(0)
	} else if userInput == "sit down next to my friend" || userInput == "sit down next to friend" || userInput == "sit with friend" || userInput == "sit with my friend" && m.currentScene == "dungeon" {
		scene = "friendHandsNote"
		prompt = "Your friend hands you a note.\nWhat do you do?\n\n> "
	} else if userInput == "light a match" || userInput == "light match" && m.currentScene == "friendHandsNote" {
		scene = "dontLeaveMeHere"
		prompt = "The note says, \"Don't leave me here.\"\nDo you leave your friend or stay?\n\n> "
	} else if userInput == "stay" && m.currentScene == "dontLeaveMeHere" {
		openBrowser("https://www.youtube.com/watch?v=g_Miz2ZqSI4")
		os.Exit(0)
	} else if userInput == "leave" && m.currentScene == "dontLeaveMeHere" {
		scene = "beach"
		prompt = "You move the barrel, find a secret tunnel, and crawl through it.\nThe tunnel leads you to a beach. What do you do?\n\n> "
	} else {
		return m.currentScene, m.currentScenePrompt, m.currentSceneGraphic
	}

	graphic, err := processANSIArt("./assets/" + scene + ".png")
	if err != nil {
		fmt.Printf("Error processing ANSI art: %v\n", err)
		os.Exit(1)
	}

	return scene, prompt, graphic
}

func blinkTick() tea.Cmd {
	return tea.Tick(blinkRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func processANSIArt(imageInput string) ([]string, error) {
	ansiStr, err := ansify.GetAnsify(imageInput)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		os.Exit(1)
	}

	if ansiStr == "" {
		return nil, fmt.Errorf("empty ANSI string received")
	}

	lines := strings.Split(ansiStr, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no lines found in ANSI art")
	}

	return lines, nil
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return fmt.Errorf("error opening browser: %v", err)
	}
	return nil
}

func main() {
	initialScene, err := processANSIArt("./assets/eXit.png")
	if err != nil {
		fmt.Printf("Error processing ANSI art: %v\n", err)
		os.Exit(1)
	}

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
