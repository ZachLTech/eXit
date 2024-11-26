package game

type model struct {
	currentScene        string
	currentScenePrompt  string
	cursorBlink         bool
	animating           bool
	animationStep       int
	animationFrameLen   int
	currentSceneGraphic []string
	userInput           string
	err                 error
	termWidth           int

	enteredTunnel bool
	isElliot      bool
	restartCount  int
	cursorSymbol  string
}
