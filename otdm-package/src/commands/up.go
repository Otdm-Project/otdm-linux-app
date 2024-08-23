package commands

import (
	"fmt"
	"time"

	"otdm-package/src/utils"
)

func Up() {
	fmt.Println("up.go started")
	utils.Websocket()
	time.Sleep(1 * time.Second)
	fmt.Println("websocket.go done")
	utils.Boot()
	time.Sleep(1 * time.Second)
	fmt.Println("boot.go done")
	fmt.Println("up.go done")
}
