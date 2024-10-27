package utils

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "io/ioutil"
    "strings"
    "github.com/gorilla/websocket"
    //"encoding/json"
)


// CallWebsocket 関数が各ステップを順に実行
func CallWebsocket() (cvIP string, svIP string, otdmPubKey string,domainName string, err error) {
    // 起動時ログ
    //var err error
    err = LogMessage(INFO, "websocket.go start")
    if err != nil {
        return "", "", "", "", err
    }


    // ステップ1: 鍵の生成
    privateKey, publicKey, err := generateKeys()
    if err != nil {
        return "", "", "", "", fmt.Errorf("Failed to generate keys: %v", err)
    }
    fmt.Printf("Generated keys: private=%s, public=%s\n", privateKey, publicKey)

    // ステップ2: 初期設定ファイル作成
    const confFilePath = "/etc/wireguard/otdm.conf"
    err = createOrEditConfig(privateKey, "", "", "", "")
    if err != nil {
        return "", "", "", "", fmt.Errorf("Failed to create/edit config: %v", err)
    }

    // ステップ3: WebSocket 通信を確立して情報を取得
    /*
    cvIP, svIP, otdmPubKey, domainName, err := getWebSocketData()
    if err != nil {
        return "", "", "", "", fmt.Errorf("Failed to retrieve data via WebSocket: %v", err)
    }
        */

    // ダミーデータの使用
    cvIP, svIP, otdmPubKey, domainName = "192.168.1.10", "10.0.0.1", "otdmPubKey", "otdm.dev"

    // ステップ4: 取得した情報を設定ファイルに追記
    err = createOrEditConfig(privateKey, cvIP, svIP, otdmPubKey, domainName)
    if err != nil {
        return "", "", "", "", fmt.Errorf("Failed to update config with received data: %v", err)
    }

    fmt.Println("Configuration setup completed.")

    // 処理終了時ログ
    err = LogMessage(INFO, "websocket.go done")
    if err != nil {
        fmt.Printf("Failed to log message: %v\n", err)
    }

    return cvIP, svIP, otdmPubKey, domainName, nil
}

// getWebSocketData はWebSocketを介してデータを取得
func getWebSocketData() (cvIP, svIP, otdmPubKey, domainName string, err error) {
    // WebSocket サーバーのURL
    url := "ws://3.93.146.222:3000"

     // WebSocket接続の確立
     c, _, err := websocket.DefaultDialer.Dial(url, nil)
     if err != nil {
         return "", "", "", "", fmt.Errorf("failed to connect to websocket server: %v", err)
     }
     defer c.Close()
 
     // メッセージの受信
     _, message, err := c.ReadMessage()
     if err != nil {
         return "", "", "", "", fmt.Errorf("failed to read message: %v", err)
     }
 
     // 平文メッセージを分割
     parts := strings.Split(string(message), ",")
     if len(parts) != 4 {
         return "", "", "", "", fmt.Errorf("received message is not valid")
     }
     
     return parts[0], parts[1], parts[2], parts[3], nil
}

// 鍵を生成する関数
func generateKeys() (privateKey, publicKey string, err error) {
    // 鍵生成のコマンド実行
    privCmd := exec.Command("wg", "genkey")
    privKey, err := privCmd.Output()
    if err != nil {
        // 具体的なエラーメッセージを含める
        return "", "", fmt.Errorf("failed to generate private key: %v", err)  
    }

    // プライベート鍵のトリム処理をして不正な文字を排除
    trimmedPrivKey := strings.TrimSpace(string(privKey))

    // 公開鍵生成のためのコマンド
    pubCmd := exec.Command("bash", "-c", fmt.Sprintf("echo %s | wg pubkey", trimmedPrivKey))
    pubKeyOutput, err := pubCmd.CombinedOutput()
    if err != nil {
        // 出力されたエラーメッセージを含める
        return "", "", fmt.Errorf("failed to generate public key: %v (%s)", err, string(pubKeyOutput))
    }

    return strings.TrimSpace(string(privKey)), strings.TrimSpace(string(pubKeyOutput)), nil
}

// otdm.confを必要なら生成または編集する
func createOrEditConfig(privateKey, cvIP, svIP, otdmPubKey, domainName string) error {
    // 設定ファイルのパスを /etc/wireguard/ に固定
	configPath := filepath.Join("/etc/wireguard", "otdm.conf")

	// /etc/wireguard ディレクトリが存在するか確認し、なければ作成
	if _, err := os.Stat("/etc/wireguard"); os.IsNotExist(err) {
		if err := os.Mkdir("/etc/wireguard", 0755); err != nil {
			return fmt.Errorf("failed to create directory /etc/wireguard: %v", err)
		}
	}

    // コンフィグテンプレートの全体
    template := fmt.Sprintf(`
[Interface]
PrivateKey = %s
Address = %s/24

[Peer]
PublicKey = %s
Endpoint = %s:51820
AllowedIPs = %s/32
PersistentKeepalive = 25
`, privateKey, cvIP, otdmPubKey, domainName, svIP)

    // ファイルが存在しなくても、初期化したい場合でも一貫してテンプレートで上書き
    return ioutil.WriteFile(configPath, []byte(template), 0644)
}
