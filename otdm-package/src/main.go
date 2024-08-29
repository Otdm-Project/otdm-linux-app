package main

import (
		"fmt"
		"os"
		"otdm-package/src/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: otdm <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		commands.RunUp()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
