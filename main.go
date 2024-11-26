package main

import (
	"eXit/game"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// Loads .env and SSH_SERVER_ENABLED value
	err := godotenv.Load()
	check(err, "Loading .env in main... Make sure you have your .env file in the root directory of this program", true)

	// Starts either the TUI SSH session or TUI program depending on if the flag is true or false.
	SSHEnabled, err := strconv.ParseBool(os.Getenv("SSH_SERVER_ENABLED"))
	check(err, "Parsing .env SSH_SERVER_ENABLED bool in main", true)

	if SSHEnabled {
		game.RunSSHGame()
	} else {
		game.RunGameLocal()
	}
}

func check(e error, check string, fatal bool) {
	if e != nil {
		fmt.Printf("Error running program - In %v: %v", check, e)
		if fatal {
			os.Exit(1)
		}
	}
}
