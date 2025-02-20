// up.go
package commands

import (
	"fmt"
	"otdm-package/src/utils"
)

// err = LogMessage(ERRO, errMessage)
// RunUpはWebSocketデータを受け取り、main.goに返す
func RunUp(httpport int) (cvIP, svIP, otdmPubKey, domainName string, err error) {
	// err 変数は既に関数の返り値で宣言されているので、新たに宣言する必要はありません
	err = utils.LogMessage(utils.INFO, "up.go start")
	// refresh.goを呼び出す
	utils.CallRefresh()

	// WebSocket.goを呼び出す
	cvIP, svIP, otdmPubKey, domainName, err = utils.CallWebsocket()
	if err != nil {
		errMessage := fmt.Sprintf("Error during WebSocket connection: %v", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}

	// 他の処理（例: CallRefresh, CallBoot など）

	if err != nil {
		errMessage := fmt.Sprintf("Error during boot: %v%v", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}

	interfaceName := "otdm"
	err = utils.ConfigureFirewall(interfaceName, cvIP, svIP, httpport)
	if err != nil {
		errMessage := fmt.Sprintf("Error during firewall configuration:", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}

	// トンネル健康状態の監視関連 (バックグラウンドで実行）
	utils.LogMessage(utils.INFO, "Starting CallWatchman in goroutine")
	//go utils.CallWatchman(svIP)
	go utils.CallWatchman(svIP)

	err = utils.LogMessage(utils.INFO, "otdm up done.")
	if err != nil {
		//fmt.Printf("Failed to log message: %v\n", err)
		errMessage := fmt.Sprintf("Failed to log message: %v\n", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}
	return cvIP, svIP, otdmPubKey, domainName, nil
}
