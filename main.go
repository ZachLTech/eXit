package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
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
	animationFrameLen   int
	currentSceneGraphic []string
	fullSceneGraphic    []string
	sceneGraphicLen     int
	userInput           string
	err                 error
}

type tickMsg time.Time
type renderTickMsg time.Time
type animationTickMsg time.Time

const (
	cursorBlinkRate = time.Millisecond * 500
	renderSpeed     = time.Millisecond * 250
)

var enteredTunnel bool = false
var isElliot bool = false
var restartCount int = 0
var cursorSymbol string = "░"

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

var animationFrames = map[string]int{
	"congratulations": 2,
	"start":           34,
	"end":             11,
}

// milliseconds per frame
var animationFramerate = map[string]time.Duration{
	"congratulations": 500,
	"start":           100,
	"end":             100,
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickAnimation("start"),
		blinkTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case renderTickMsg:
		m.fullSceneGraphic = m.currentSceneGraphic
		m.currentSceneGraphic = []string{}

		for i := 0.0; i < math.Floor(float64(m.animationFrameLen)*0.3); i++ {
			m.currentSceneGraphic = append(m.currentSceneGraphic, m.fullSceneGraphic[int(i)])
			m.fullSceneGraphic = m.fullSceneGraphic[1:]
		}

	case animationTickMsg:
		m.animationFrameLen = animationFrames[m.currentScene]
		cursorSymbol = ""
		if m.animationStep != m.animationFrameLen {
			m.animationStep++
			m.userInput = ""
			m.currentScenePrompt = m.handleAnimationInput(m.currentScene)

			m.currentSceneGraphic, err = processANSIArt("./assets/animations/" + m.currentScene + "/" + m.currentScene + strconv.Itoa(m.animationStep+1) + ".png")
			if err != nil {
				fmt.Printf("Error processing ANSI art: %v\n", err)
				os.Exit(1)
			}

			re := regexp.MustCompile(`\d`)
			m.currentScene = re.ReplaceAllString(m.currentScene, "")

			return m, tickAnimation(m.currentScene)
		} else {
			cursorSymbol = "░"
			m.animationStep = 0
			m.animationFrameLen = 0
			m.currentSceneGraphic, err = processANSIArt("./assets/animations/" + m.currentScene + "/" + m.currentScene + ".png")
			if err != nil {
				fmt.Printf("Error processing ANSI art: %v\n", err)
				os.Exit(1)
			}

			// exceptions
			if m.currentScene == "start" {
				m.currentScenePrompt = "PRESS ANY KEY TO START"
				cursorSymbol = ""
			} else if m.currentScene == "end" { // Now I can put anything here whenever the user reaches the real ending ;)... this will do for now hehehe
				fmt.Printf("Hello Elliot... Redirecting to https://www.youtube.com/watch?v=g_Miz2ZqSI4")
				openBrowser("https://www.youtube.com/watch?v=g_Miz2ZqSI4")
				os.Exit(0)
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
			if m.currentScene == "start" && m.animationFrameLen == 0 {
				cursorSymbol = "░"
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
			m.currentScene, m.currentScenePrompt, m.currentSceneGraphic, m.sceneGraphicLen, cmd = m.handleInput(m.userInput)
			m.userInput = ""
		case tea.KeyBackspace:
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
			}
		case tea.KeyRunes:
			if m.currentScene == "start" && m.animationFrameLen == 0 {
				cursorSymbol = "░"
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

func (m model) handleInput(userInput string) (string, string, []string, int, tea.Cmd) {
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	var scene string
	var prompt string
	var graphic []string
	var graphicLen int
	var err error
	var cmd tea.Cmd = nil
	customPrompt := false

	if userInput == "move the barrel" || userInput == "move barrel" && m.currentScene == "dungeon" {
		scene = "secretTunnel"
	} else if userInput == "enter the tunnel" || userInput == "enter tunnel" && m.currentScene == "secretTunnel" {
		enteredTunnel = true
		scene = "friendTooWeak"
	} else if userInput == "read the note" || userInput == "read note" && m.currentScene == "friendTooWeak" {
		scene = "friendHandsNote"
	} else if userInput == "leave" && (m.currentScene == "friendHandsNote" || m.currentScene == "friendTooWeak" || m.currentScene == "dontLeaveMeHere") {
		scene = "beach"
		if m.currentScene == "dontLeaveMeHere" || m.currentScene == "friendHandsNote" {
			if !enteredTunnel {
				customPrompt = true
				prompt = "You move the barrel, find a secret tunnel, and crawl through it.\nThe tunnel leads you to a beach. What do you do?\n\n> "
			}
		}
	} else if userInput == "look" || userInput == "look around" && m.currentScene == "beach" {
		scene = "ship"
	} else if userInput == "get on the boat" || userInput == "get on boat" || userInput == "get on" && m.currentScene == "ship" {
		scene = "congratulations"
		cmd = tickAnimation(scene)
	} else if userInput == "yes" && (m.currentScene == "congratulations" || m.currentScene == "sadEnding") {
		restartCount++
		enteredTunnel = false
		if m.currentScene == "congratulations" {
			isElliot = true
		}
		scene = "dungeon"
	} else if userInput == "no" && (m.currentScene == "congratulations" || m.currentScene == "sadEnding") {
		os.Exit(0)
	} else if userInput == "sit down next to my friend" || userInput == "sit next to friend" || userInput == "sit next to my friend" || userInput == "sit down next to friend" || userInput == "sit with friend" || userInput == "sit with my friend" && m.currentScene == "dungeon" {
		scene = "friendHandsNote"
		if m.currentScene == "dungeon" {
			customPrompt = true
			prompt = "Your friend hands you a note, but it is too dark to read.\nWhat do you do?\n\n> "
		}
	} else if userInput == "light a match" || userInput == "light match" && m.currentScene == "friendHandsNote" {
		scene = "dontLeaveMeHere"
	} else if userInput == "stay" && m.currentScene == "dontLeaveMeHere" {
		if restartCount == 1 && !enteredTunnel && isElliot {
			scene = "end"
			cmd = tickAnimation(scene)
		} else {
			scene = "sadEnding"
			customPrompt = true
			prompt = "Your friend is happy that you stayed...\nSoon you and your friend died due to starvation.\nDo you want to play again?\n\n> "
		}
	} else {
		return m.currentScene, m.currentScenePrompt, m.currentSceneGraphic, m.sceneGraphicLen, nil
	}

	if cmd != nil && scene != "sadEnding" {
		graphic, err = processANSIArt("./assets/animations/" + scene + "/" + scene + "1.png")
	} else if scene != "sadEnding" {
		graphicLen = len(graphic)
		// graphic, err = processANSIArt("./assets/" + scene + ".png")
	}
	if err != nil {
		fmt.Printf("Error processing ANSI art: %v\n", err)
		os.Exit(1)
	}

	if !customPrompt {
		prompt = prompts[scene]
	}

	return scene, prompt, graphic, graphicLen, cmd
}

func (m model) handleAnimationInput(scene string) string {
	var prompt string = ""
	if strings.Contains(scene, "congratulations") {
		return prompts["congratulations"]
	}

	return prompt
}

func blinkTick() tea.Cmd {
	return tea.Tick(cursorBlinkRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func tickAnimation(animationName string) tea.Cmd {
	return tea.Tick((time.Millisecond * animationFramerate[animationName]), func(t time.Time) tea.Msg {
		return animationTickMsg(t)
	})
}

func tickRender() tea.Cmd {
	return tea.Tick(renderSpeed, func(t time.Time) tea.Msg {
		return renderTickMsg(t)
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
	initialModel := model{
		currentScene: "start",
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
