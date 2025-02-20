package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"otdm-package/src/commands"
	"otdm-package/src/utils"
)

// グローバル変数
var (
	cvIP, svIP, otdmPubKey, domainName, errMessage string
)

// バージョン情報の定義
const Version = "0.0.1"

func main() {
	globalFlagSet := flag.NewFlagSet("global", flag.ExitOnError)
	httpport := globalFlagSet.Int("p", 80, "httpport")

	// rootユーザーか確認
	usr, err := user.Current()
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching user info: %v", err)
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
		cvIP, svIP, otdmPubKey, domainName, err = commands.RunUp(*httpport)
		if err != nil {
			errMessage := fmt.Sprintf("Error during up:%v", err)
			utils.ErrLogMessage(errMessage)
		} else {
			// バックグラウンドでWatchmanを起動
			go utils.CallWatchman(svIP)
		}
		break
	case "down":
		// `watchman` のプロセスを停止
		exec.Command("pkill", "-f", "otdm-watchman").Run()

		if err := commands.RunDown(); err != nil {
			utils.ErrLogMessage(fmt.Sprintf("Error during down:%v", err))
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

	select {}
}
