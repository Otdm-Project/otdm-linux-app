package commands

import (
	"fmt"
	"otdm-package/src/utils"
)

func RunUp() {
	utils.CallBoot()
    // CallWebsocket関数でWebSocket関連の処理を行います。
    if err := utils.CallWebsocket(); err != nil {
        fmt.Println("Error during WebSocket connection:", err)
        return
    }
	utils.CallBoot()
	utils.CallNetwork()
}
