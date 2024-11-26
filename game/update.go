package game

import (
	"fmt"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case animationTickMsg:
		m.animating = true
		m.animationFrameLen = animationFrames[m.currentScene]
		m.cursorSymbol = ""
		if m.animationStep != m.animationFrameLen {
			m.animationStep++
			m.userInput = ""

			m.currentSceneGraphic, err, cmd = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+strconv.Itoa(m.animationStep+1)+".png", m.termWidth)
			if err != nil {
				fmt.Printf("Error processing ANSI art: %v\n", err)
				cmd = tea.Quit
			}

			re := regexp.MustCompile(`\d`)
			m.currentScene = re.ReplaceAllString(m.currentScene, "")

			if cmd != nil {
				return m, cmd
			} else {
				return m, tickAnimation(m.currentScene)
			}
		} else {
			m.cursorSymbol = "░"
			m.animationStep = 0
			m.animationFrameLen = 0
			m.currentSceneGraphic, err, cmd = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+".png", m.termWidth)
			if err != nil {
				fmt.Printf("Error processing ANSI art: %v\n", err)
				cmd = tea.Quit
			}

			// exceptions
			if m.currentScene == "start" {
				m.currentScenePrompt = "PRESS ANY KEY TO START"
				m.cursorSymbol = ""
			} else if m.currentScene == "end" { // Now I can put anything here whenever the user reaches the real ending ;)... this will do for now hehehe
				return m, tickElliot()
			}

			m.animating = false
			return m, cmd
		}

	case tickMsg:
		m.cursorBlink = !m.cursorBlink
		return m, blinkTick()

	case easterEggTickMsg:
		fmt.Printf("Hello Elliot... Redirecting to https://www.youtube.com/watch?v=g_Miz2ZqSI4")
		cmd = tea.Batch(
			openBrowser("https://www.youtube.com/watch?v=g_Miz2ZqSI4"),
			tea.Quit,
		)

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		if m.animating {
			m.currentSceneGraphic, err, cmd = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+strconv.Itoa(m.animationStep+1)+".png", m.termWidth)
		} else {
			m.currentSceneGraphic, err, cmd = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+".png", m.termWidth)
		}
		if err != nil {
			fmt.Printf("Error processing ANSI art: %v\n", err)
			cmd = tea.Quit
		}

		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.currentScene == "start" && m.animationFrameLen == 0 {
				m.cursorSymbol = "░"
				m.userInput = ""
				m.currentScene = "dungeon"
				m.currentScenePrompt = "You're trapped in a dungeon with your friend.\nYou see a barrel. What do you do?\n\n> "
				m.currentSceneGraphic, err, _ = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+"1.png", m.termWidth)
				if err != nil {
					fmt.Printf("Error processing ANSI art: %v\n", err)
					cmd = tea.Quit
				}
				m.animating = true
				cmd = tickAnimation(m.currentScene)

				return m, cmd
			}
			m.currentScene, m.currentScenePrompt, m.currentSceneGraphic, m.animating, m.restartCount, m.enteredTunnel, m.isElliot, cmd = m.handleInput(m.userInput)
			m.userInput = ""
		case tea.KeyBackspace:
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
			}

		case tea.KeySpace:
			if m.currentScene == "start" && m.animationFrameLen == 0 {
				m.cursorSymbol = "░"
				m.userInput = ""
				m.currentScene = "dungeon"
				m.currentScenePrompt = "You're trapped in a dungeon with your friend.\nYou see a barrel. What do you do?\n\n> "
				m.currentSceneGraphic, err, _ = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+"1.png", m.termWidth)
				if err != nil {
					fmt.Printf("Error processing ANSI art: %v\n", err)
					cmd = tea.Quit
				}
				m.animating = true
				cmd = tickAnimation(m.currentScene)

				return m, cmd
			}
			m.userInput += " "

		case tea.KeyRunes:
			if m.currentScene == "start" && m.animationFrameLen == 0 {
				m.cursorSymbol = "░"
				m.userInput = ""
				m.currentScene = "dungeon"
				m.currentScenePrompt = "You're trapped in a dungeon with your friend.\nYou see a barrel. What do you do?\n\n> "
				m.currentSceneGraphic, err, _ = processANSIArt("./assets/animations/"+m.currentScene+"/"+m.currentScene+"1.png", m.termWidth)
				if err != nil {
					fmt.Printf("Error processing ANSI art: %v\n", err)
					cmd = tea.Quit
				}
				m.animating = true
				cmd = tickAnimation(m.currentScene)

				return m, cmd
			}

			m.userInput += string(msg.Runes)
		}
	}
	return m, cmd
}
