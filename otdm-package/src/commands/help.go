package commands

import "fmt"

// showHelp prints the help message for the otdm command.
func ShowHelp() error {
	helpMessage := `
otdm: A command-line tool for managing OTDM tunnels and connections.

Usage:
  otdm <command>

Available Commands:
  up        Start the OTDM tunnel.
  down      Stop the OTDM tunnel.
  status    Display the current OTDM tunnel status.
  version   Show the version information.
  help      Display this help message.

Examples:
  sudo otdm up       # Start the tunnel
  sudo otdm status   # Check the tunnel status
  sudo otdm down     # Stop the tunnel
  sudo otdm help     # Display this help message

Note:
  This command must be run with root privileges (use sudo).

For further help, please visit the following url ☞(https://otdm.dev)
`
	fmt.Println(helpMessage)
	return nil
}
