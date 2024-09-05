package commands

import (
	"fmt"
	"otdm-package/src/utils"
)

func RunUp() {
	// 開始時ログ呼び出し
	var err error
	err = utils.LogMessage(utils.INFO, "otdm up start.")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}

	// refresh.go呼び出し
	utils.CallRefresh()

	// websocket.go呼び出し
    // CallWebsocket関数でWebSocket関連の処理を行います。
    if err := utils.CallWebsocket(); err != nil {
		fmt.Println("Error during WebSocket connection:", err)
       	return
    }

	// boot.go呼び出し
	utils.CallBoot()
	
	// network.go呼び出し
	// インターフェース名を指定
    interfaceName := "otdm" 
    if err := utils.ConfigureFirewall(interfaceName); err != nil {
       	fmt.Println("Error during firewall configuration:", err)
       	return
    }

	// watchman.go呼び出し
	// 仮のIP受け渡し
	svIP := "192.168.1.2" 
	utils.CallWatchman(svIP)

	//起動終了時ログ
	err = utils.LogMessage(utils.INFO, "otdm up done.")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}
}
