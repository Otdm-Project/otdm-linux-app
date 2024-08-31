package main

import (
		"fmt"
		"os"
		"os/user"
		"otdm-package/src/commands"
)

func main() {

	// rootユーザーか確認
    	usr, err := user.Current()
    	if err != nil {
        	fmt.Println("Error fetching user info:", err)
        	os.Exit(1)
    	}

    	if usr.Uid != "0" {
        	fmt.Println("This command must be run as root. Use sudo.")
        	os.Exit(1)
	    }

    	if len(os.Args) < 2 {
        	fmt.Println("Usage: otdm <command>")
        	os.Exit(1)
    	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: otdm <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		commands.RunUp()
	case "down":
        	if err := commands.RunDown(); err != nil {
            		fmt.Println("Error during down:", err)
        	}
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
