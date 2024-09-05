package utils

import (
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
	"strings"
    // "github.com/gorilla/websocket"
)

// CallWebsocket 関数が各ステップを順に実行
func CallWebsocket() error {
    // 起動時ログ
    var err error
    err = LogMessage(INFO, "websocket.go start")
    if err != nil {
        fmt.Printf("Failed to log message: %v\n", err)
    }

    // ステップ1: 鍵の生成
    privateKey, publicKey, err := generateKeys()
    if err != nil {
        return fmt.Errorf("Failed to generate keys: %v", err)
    }
    fmt.Printf("Generated keys: private=%s, public=%s\n", privateKey, publicKey)

    // ステップ2: 初期設定ファイル作成
    err = createOrEditConfig(privateKey, "", "", "", "")
    if err != nil {
        return fmt.Errorf("Failed to create/edit config: %v", err)
    }

    // ステップ3: WebSocket 通信を確立して情報を取得 (ダミーデータ使用)
    cvIP, svIP, otdmPubKey, domainName := getWebSocketData()

    // ステップ4: 取得した情報を設定ファイルに追記
    err = createOrEditConfig(privateKey, cvIP, svIP, otdmPubKey, domainName)
    if err != nil {
        return fmt.Errorf("Failed to update config with received data: %v", err)
    }

    fmt.Println("Configuration setup completed.")

    // 処理終了時ログ
    err = LogMessage(INFO, "websocket.go done")
    if err != nil {
        fmt.Printf("Failed to log message: %v\n", err)
    }

    return nil
}

// getWebSocketData はWebSocketを介してデータを取得する (ダミーデータを使用)
func getWebSocketData() (cvIP, svIP, otdmPubKey, domainName string) {
    // 実際のWebSocket接続コード（コメントアウト）
    /*
    c, _, err := websocket.DefaultDialer.Dial("ws://example.com/socket", nil)
    if err != nil {
        return "", "", "", "", fmt.Errorf("failed to connect websocket: %v", err)
    }
    defer c.Close()

    _, message, err := c.ReadMessage()
    if err != nil {
        return "", "", "", "", fmt.Errorf("failed to read message: %v", err)
    }
    */

    // ダミーデータの利用(テスト用)
    return "192.168.1.2", "192.168.1.1", "dummyOtdmPublicKey", "example.com"
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
    configPath := "config/otdm.conf"

    // 初期ファイルの生成
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
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

        return ioutil.WriteFile(configPath, []byte(template), 0644)
    }

    // 既存ファイルの追記
    additionalPeer := fmt.Sprintf(`
[Peer]
PublicKey = <対向の公開鍵>
AllowedIPs = <許可されたIPアドレス>
`)
    content, err := ioutil.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("failed to read existing config: %v", err)
    }

    newContent := string(content) + additionalPeer

    // 処理終了時ログ
    err = LogMessage(INFO, "websocket.go done")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}

    return ioutil.WriteFile(configPath, []byte(newContent), 0644)
}