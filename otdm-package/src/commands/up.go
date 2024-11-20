package commands

import (
	"fmt"
	"otdm-package/src/utils"
)

// err = LogMessage(ERRO, errMessage)
// RunUpはWebSocketデータを受け取り、main.goに返す
func RunUp() (cvIP, svIP, otdmPubKey, domainName string, err error) {
	// err 変数は既に関数の返り値で宣言されているので、新たに宣言する必要はありません
	err = utils.LogMessage(utils.INFO, "up.go start")

	utils.CallRefresh()
	// WebSocketからのデータを受け取る
	cvIP, svIP, otdmPubKey, domainName, err = utils.CallWebsocket()
	if err != nil {
		//fmt.Printf("Error during WebSocket connection: %v\n", err)
		errMessage := fmt.Sprintf("Error during WebSocket connection: %v", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}

	// 他の処理（例: CallRefresh, CallBoot など）

	err = utils.CallBoot()
	if err != nil {
		//fmt.Printf("Error during boot: %v\n", err)
		errMessage := fmt.Sprintf("Error during boot: %v%v", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}

	interfaceName := "otdm"
	err = utils.ConfigureFirewall(interfaceName)
	if err != nil {
		//fmt.Printf("Error during firewall configuration: %v\n", err)
		errMessage := fmt.Sprintf("Error during firewall configuration:", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return "", "", "", "", err
	}

	// トンネル健康状態の監視関連 (バックグラウンドで実行）
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
