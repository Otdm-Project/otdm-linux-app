package main

import (
	"fmt"
	"os"

	"otdm-package/src/commands"
)

func main() {
	if len(os.Args) < 2 {
		commands.Help()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		commands.Up()
	case "down":
		commands.Down()
	case "status":
		commands.Status()
	case "--help", "-h":
		commands.Help()
	case "--version", "-v":
		commands.Version()
	default:
		fmt.Println("Invalid command. Use 'otdm --help' for usage.")
		os.Exit(1)
	}
}

