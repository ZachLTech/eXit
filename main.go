package main

import (
	"fmt"
	"os"
	"strconv"
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
	animationStep       int
	currentSceneGraphic []string
	userInput           string
	err                 error
}

type tickMsg time.Time
type animationTickMsg time.Time

const (
	blinkRate         = time.Millisecond * 500
	cursorSymbol      = "░"
	animationStepRate = time.Millisecond * 200
)

var prompts = map[string]string{
	"dungeon":         "You're trapped in a dungeon with your friend. You see a barrel. What do you do?\n\n> ",
	"secretTunnel":    "The barrel rolls aside and you find a secret tunnel\nWhat do you do?\n\n> ",
	"friendTooWeak":   "You start to escape but your friend is too weak to\ngo with you. They hand you a note.\nWhat do you do?\n\n> ",
	"friendHandsNote": "It is too dark to read the note.\nWhat do you do?\n\n> ",
	"beach":           "You crawl through the tunnel and the tunnel leads\nyou to a beach. What do you do?\n\n> ",
	"ship":            "In the water you see a boat.\nWhat do you do?\n\n> ",
	"congratulations": "Congratulations, you're heading to a new world!\nDo you want to play again?\n\n> ",
	"dontLeaveMeHere": "The note says, \"Don't leave me here.\"\nDo you leave your friend or stay?\n\n> ",
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		blinkTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case animationTickMsg:
		if m.animationStep != 2 {
			m.animationStep++
			m.userInput = ""
			m.currentScene = m.currentScene + strconv.Itoa(m.animationStep)
			m.currentScenePrompt = m.handleAnimationInput(m.currentScene)

			sceneNameLen := len(m.currentScene)
			m.currentScene = m.currentScene[:sceneNameLen-1]

			m.currentSceneGraphic, err = processANSIArt("./assets/" + m.currentScene + strconv.Itoa(m.animationStep+1) + ".png")
			if err != nil {
				fmt.Printf("Error processing ANSI art: %v\n", err)
				os.Exit(1)
			}

			return m, tickAnimation()
		} else {
			m.animationStep = 0
			sceneNameLen := len(m.currentScene)

			if sceneNameLen > 0 && m.currentScene[sceneNameLen-1] == '3' {
				m.currentScene = m.currentScene[:sceneNameLen-1]
			}

			return m, nil
		}

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
			m.currentScene, m.currentScenePrompt, m.currentSceneGraphic, cmd = m.handleInput(m.userInput)
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
	return m, cmd
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

func (m model) handleInput(userInput string) (string, string, []string, tea.Cmd) {
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	var scene string
	var prompt string
	var graphic []string
	var err error
	var cmd tea.Cmd = nil

	customPrompt := "You move the barrel, find a secret tunnel, and crawl through it.\nThe tunnel leads you to a beach. What do you do?\n\n> "

	if userInput == "move the barrel" || userInput == "move barrel" && m.currentScene == "dungeon" {
		scene = "secretTunnel"
	} else if userInput == "enter the tunnel" || userInput == "enter tunnel" && m.currentScene == "secretTunnel" {
		scene = "friendTooWeak"
	} else if userInput == "read the note" || userInput == "read note" && m.currentScene == "friendTooWeak" {
		scene = "friendHandsNote"
	} else if userInput == "leave" && (m.currentScene == "friendHandsNote" || m.currentScene == "friendTooWeak") {
		scene = "beach"
	} else if userInput == "look" || userInput == "look around" && m.currentScene == "beach" {
		scene = "ship"
	} else if userInput == "get on the boat" || userInput == "get on boat" || userInput == "get on" && m.currentScene == "ship" {
		scene = "congratulations"
		cmd = tickAnimation()
	} else if userInput == "yes" && m.currentScene == "congratulations" {
		scene = "dungeon"
	} else if userInput == "no" && m.currentScene == "congratulations" {
		os.Exit(0)
	} else if userInput == "sit down next to my friend" || userInput == "sit down next to friend" || userInput == "sit with friend" || userInput == "sit with my friend" && m.currentScene == "dungeon" {
		scene = "friendHandsNote"
	} else if userInput == "light a match" || userInput == "light match" && m.currentScene == "friendHandsNote" {
		scene = "dontLeaveMeHere"
	} else if userInput == "stay" && m.currentScene == "dontLeaveMeHere" {
		openBrowser("https://www.youtube.com/watch?v=g_Miz2ZqSI4")
		os.Exit(0)
	} else if userInput == "leave" && m.currentScene == "dontLeaveMeHere" {
		scene = "beach"
		prompt = customPrompt
	} else {
		return m.currentScene, m.currentScenePrompt, m.currentSceneGraphic, nil
	}

	if cmd != nil {
		graphic, err = processANSIArt("./assets/" + scene + "1.png")
	} else {
		graphic, err = processANSIArt("./assets/" + scene + ".png")
	}
	if err != nil {
		fmt.Printf("Error processing ANSI art: %v\n", err)
		os.Exit(1)
	}

	if prompt != customPrompt {
		prompt = prompts[scene]
	}

	return scene, prompt, graphic, cmd
}

func (m model) handleAnimationInput(scene string) string {
	var prompt string = ""
	if strings.Contains(scene, "congratulations") {
		return prompts["congratulations"]
	}

	return prompt
}

func blinkTick() tea.Cmd {
	return tea.Tick(blinkRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func tickAnimation() tea.Cmd {
	return tea.Tick(animationStepRate, func(t time.Time) tea.Msg {
		return animationTickMsg(t)
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
