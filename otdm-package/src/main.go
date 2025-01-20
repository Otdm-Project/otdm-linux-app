// main.go
package main

import (
	"fmt"
	"os"
	"os/user"
	"otdm-package/src/commands"
	"otdm-package/src/utils"
)

// グローバルにデータを保持する変数宣言
var (
	cvIP, svIP, otdmPubKey, domainName, errMessage string
	httpPort                                       int
)

// バージョン情報の定義
const Version = "0.0.1"

func main() {
	// rootユーザーか確認
	usr, err := user.Current()
	if err != nil {
		//fmt.Println("Error fetching user info:", err)
		errMessage := fmt.Sprintf("Error fetching user info:%v", err)
		utils.ErrLogMessage(errMessage)
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

	switch os.Args[1] {
	case "up":

		// WebSocketからのデータを受け取る
		cvIP, svIP, otdmPubKey, domainName, httpPort, err = commands.RunUp()
		if err != nil {
			//fmt.Printf("Error during up: %v\n", err)
			errMessage := fmt.Sprintf("Error during up:%v", err)
			utils.ErrLogMessage(errMessage)
		} else {
			go utils.CallWatchman(svIP)
		}
		break
	case "down":
		if err := commands.RunDown(); err != nil {
			//fmt.Println("Error during down:", err)
			errMessage := fmt.Sprintf("Error during down:%v", err)
			utils.ErrLogMessage(errMessage)
		}
	case "status":
		commands.ShowStatus()
	case "version":
		fmt.Println("otdm version : ", Version)
	case "help":
		if err := commands.ShowHelp(); err != nil {
			fmt.Println("Error displaying help:", err)
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
