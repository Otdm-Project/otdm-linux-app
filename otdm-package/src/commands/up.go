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
	// インターフェース名を指定してください
    interfaceName := "otdm" 
    if err := utils.ConfigureFirewall(interfaceName); err != nil {
        fmt.Println("Error during firewall configuration:", err)
        return
    }
}
