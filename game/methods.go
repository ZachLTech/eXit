package game

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ZachLTech/ansify"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/joho/godotenv"
)

/************************************ Runner functions **************************************/
/************************ runs either the local game or SSH server **************************/

func RunGameLocal() {
	initialModel := model{
		currentScene:  "start",
		enteredTunnel: false,
		isElliot:      false,
		restartCount:  0,
		cursorSymbol:  "░",
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

func RunSSHGame() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error running program - In Loading .env to run SSH Game: %v", err)
		os.Exit(1)
	}

	var (
		host = os.Getenv("HOST")
		port = os.Getenv("PORT")
	)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),

		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

/*********************************** SSH Server Handler *************************************/

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{
		currentScene:  "start",
		enteredTunnel: false,
		isElliot:      false,
		restartCount:  0,
		cursorSymbol:  "░",
	}
	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithInput(os.Stdin)}
}

/************************************** Game Tickers ****************************************/

func blinkTick() tea.Cmd {
	return tea.Tick(cursorBlinkRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func tickElliot() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return easterEggTickMsg(t)
	})
}

func tickAnimation(animationName string) tea.Cmd {
	return tea.Tick((time.Millisecond * animationFramerate[animationName]), func(t time.Time) tea.Msg {
		return animationTickMsg(t)
	})
}

/************************************* Helpers & Utils **************************************/

func (m model) handleInput(userInput string) (string, string, []string, bool, int, bool, bool, tea.Cmd) {
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	var scene string
	var prompt string
	var graphic []string
	var animating bool
	var err error
	var cmd tea.Cmd = nil
	customPrompt := false

	if userInput == "move the barrel" || userInput == "move barrel" && m.currentScene == "dungeon" {
		scene = "secretTunnel"
	} else if userInput == "enter the tunnel" || userInput == "enter tunnel" && m.currentScene == "secretTunnel" {
		m.enteredTunnel = true
		scene = "friendTooWeak"
	} else if userInput == "read the note" || userInput == "read note" && m.currentScene == "friendTooWeak" {
		scene = "friendHandsNote"
	} else if userInput == "leave" && (m.currentScene == "friendHandsNote" || m.currentScene == "friendTooWeak" || m.currentScene == "dontLeaveMeHere") {
		scene = "beach"
		if m.currentScene == "dontLeaveMeHere" || m.currentScene == "friendHandsNote" {
			if !m.enteredTunnel {
				customPrompt = true
				prompt = "You move the barrel, find a secret tunnel, and crawl through it.\nThe tunnel leads you to a beach. What do you do?\n\n> "
			}
		}
	} else if (userInput == "look" || userInput == "look around") && (m.currentScene == "beach") {
		scene = "ship"
	} else if userInput == "get on the boat" || userInput == "get on boat" || userInput == "get on" && m.currentScene == "ship" {
		scene = "congratulations"
	} else if userInput == "yes" && (m.currentScene == "congratulations" || m.currentScene == "sadEnding") {
		m.restartCount++
		m.enteredTunnel = false
		if m.currentScene == "congratulations" {
			m.isElliot = true
		}
		scene = "dungeon"
	} else if userInput == "no" && (m.currentScene == "congratulations" || m.currentScene == "sadEnding") {
		cmd = tea.Quit
	} else if userInput == "sit down next to my friend" || userInput == "sit next to friend" || userInput == "sit next to my friend" || userInput == "sit down next to friend" || userInput == "sit with friend" || userInput == "sit with my friend" && m.currentScene == "dungeon" {
		scene = "friendHandsNote"
		if m.currentScene == "dungeon" {
			customPrompt = true
			prompt = "Your friend hands you a note, but it is too dark to read.\nWhat do you do?\n\n> "
		}
	} else if userInput == "light a match" || userInput == "light match" && m.currentScene == "friendHandsNote" {
		scene = "dontLeaveMeHere"
	} else if userInput == "stay" && m.currentScene == "dontLeaveMeHere" {
		if m.restartCount == 1 && !m.enteredTunnel && m.isElliot {
			scene = "end"
		} else {
			scene = "sadEnding"
			customPrompt = true
			prompt = "Your friend is happy that you stayed...\nSoon you and your friend died due to starvation.\nDo you want to play again?\n\n> "
		}
	} else {
		return m.currentScene, m.currentScenePrompt, m.currentSceneGraphic, m.animating, m.restartCount, m.enteredTunnel, m.isElliot, nil
	}

	if scene != "sadEnding" {
		graphic, err, _ = processANSIArt("./assets/animations/"+scene+"/"+scene+"1.png", m.termWidth)
		if err != nil {
			fmt.Printf("Error processing ANSI art: %v\n", err)
			cmd = tea.Quit
		}
		animating = true
		cmd = tickAnimation(scene)
	}

	if !customPrompt {
		prompt = prompts[scene]
	}

	return scene, prompt, graphic, animating, m.restartCount, m.enteredTunnel, m.isElliot, cmd
}

func processANSIArt(imageInput string, termWidth int) ([]string, error, tea.Cmd) {
	ansiStr, err := ansify.GetAnsifyCustomWidth(imageInput, termWidth)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		return nil, fmt.Errorf("error loading image: %v", err), tea.Quit
	}

	if ansiStr == "" {
		return nil, fmt.Errorf("empty ANSI string received"), tea.Quit
	}

	lines := strings.Split(ansiStr, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no lines found in ANSI art"), tea.Quit
	}

	return lines, nil, nil
}

type editorFinishedMsg struct{ err error }

func openBrowser(url string) tea.Cmd {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return nil
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}
