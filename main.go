package main

import (
	"fmt"
	"os"

	"github.com/ZachLTech/ansify"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	asciiArt  string
	userInput string
}

func (m model) Init() tea.Cmd {
	// Initialize your application here
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle user input and application updates here
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			// Process the input
			fmt.Println("You entered:", m.userInput)
			m.userInput = "" // Reset input after processing
			return m, nil
		} else if msg.Type == tea.KeyRunes {
			m.userInput += string(msg.Runes)
		}
		if m.userInput == "q" {
			os.Exit(0)
		}
	}
	return m, nil
}

func (m model) View() string {
	// Render the view here	// inputHeight := terminalHeight / 5   // 20% for input area

	return ansify.GetAnsify("./assets/eXit.png")
}

func main() {
	// Load your ASCII art
	asciiArt := ansify.GetAnsify("./assets/eXit.png")

	initialModel := model{asciiArt: asciiArt}

	if _, err := tea.NewProgram(initialModel).Run(); err != nil {
		fmt.Printf("Could not start the program: %v", err)
		os.Exit(1)
	}
}
