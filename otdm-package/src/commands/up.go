package commands

import (
    "fmt"
    "otdm-package/src/utils"
)

// RunUpはWebSocketデータを受け取り、main.goに返す
func RunUp() (cvIP, svIP, domainName string, err error) {
    // err 変数は既に関数の返り値で宣言されているので、新たに宣言する必要はありません

    // WebSocketからのデータを受け取る
    cvIP, svIP, domainName, err = utils.CallWebsocket()
    if err != nil {
        fmt.Printf("Error during WebSocket connection: %v\n", err)
        return "", "", "", err
    }

    // 他の処理（例: CallRefresh, CallBoot など）
    utils.CallRefresh()

    err = utils.CallBoot()
    if err != nil {
        fmt.Printf("Error during boot: %v\n", err)
        return "", "", "", err
    }

    interfaceName := "otdm"
    err = utils.ConfigureFirewall(interfaceName)
    if err != nil {
        fmt.Printf("Error during firewall configuration: %v\n", err)
        return "", "", "", err
    }

    // トンネル健康状態の監視関連 (バックグラウンドで実行）
    go utils.CallWatchman(svIP)

    err = utils.LogMessage(utils.INFO, "otdm up done.")
    if err != nil {
        fmt.Printf("Failed to log message: %v\n", err)
        return "", "", "", err
    }

    return cvIP, svIP, domainName, nil
}