package game

import "time"

var prompts = map[string]string{
	"dungeon":         "You're trapped in a dungeon with your friend.\n You see a barrel. What do you do?\n\n> ",
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
	"start":           38,
	"end":             11,
	"beach":           defaultAnimationFrames,
	"dontLeaveMeHere": defaultAnimationFrames,
	"dungeon":         defaultAnimationFrames,
	"friendHandsNote": defaultAnimationFrames,
	"friendTooWeak":   defaultAnimationFrames,
	"secretTunnel":    defaultAnimationFrames,
	"ship":            defaultAnimationFrames,
}

// milliseconds per frame
var animationFramerate = map[string]time.Duration{
	"congratulations": 500,
	"start":           100,
	"end":             100,
	"beach":           defaultAnimationMPF,
	"dontLeaveMeHere": defaultAnimationMPF,
	"dungeon":         defaultAnimationMPF,
	"friendHandsNote": defaultAnimationMPF,
	"friendTooWeak":   defaultAnimationMPF,
	"secretTunnel":    defaultAnimationMPF,
	"ship":            defaultAnimationMPF,
}
