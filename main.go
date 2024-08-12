package main

import (
	"fmt"
	"os"
	"strings"

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
	}
	return m, nil
}

func (m model) View() string {
	// Render the view here
	terminalHeight := 20                // Assume a fixed height for simplicity
	artHeight := terminalHeight * 4 / 5 // 80% for ASCII art
	// inputHeight := terminalHeight / 5   // 20% for input area

	art := strings.Split(m.asciiArt, "\n")
	if len(art) > artHeight {
		art = art[:artHeight] // Trim the art if it's too tall
	}
	artArea := strings.Join(art, "\n")

	inputArea := fmt.Sprintf("\n%s", m.userInput) // Add a newline to separate from art

	return artArea + inputArea
}

func main() {
	// Load your ASCII art
	asciiArt := `Your ASCII art here`

	initialModel := model{asciiArt: asciiArt}

	if _, err := tea.NewProgram(initialModel).Run(); err != nil {
		fmt.Printf("Could not start the program: %v", err)
		os.Exit(1)
	}
}
