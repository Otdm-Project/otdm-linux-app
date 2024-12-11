package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

const confFilePath = "/etc/wireguard/otdm.conf"
const webscketURL = "ws://exsample.com:8080"

// メッセージ1の構造体
type WebSocketResponse struct {
	ClientPublicKey string `json:"client_public_key"`
	ServerPublicKey string `json:"server_public_key"`
	ClientVirtualIP string `json:"vpn_ip_client"`
	ServerVirtualIP string `json:"vpn_ip_server"`
	Subdomain       string `json:"subdomain"`
}

// メッセージ2の構造体
type WebSocketMessage struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// CallWebsocket 関数が各ステップを順に実行
func CallWebsocket() (cvIP string, svIP string, ServerPubKey string, domainName string, err error) {
	// 起動時ログ
	err = LogMessage(INFO, "websocket.go start")
	if err != nil {
		errMessage := fmt.Sprintf("Failed to websocket.go start: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	// ステップ1: 鍵の生成
	privateKey, publicKey, err := generateKeys()
	if err != nil {
		errMessage := fmt.Sprintf("Failed to generate keys: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}
	LogMessage(INFO, fmt.Sprintf("GeneratePrivateKey:%v\n", privateKey))
	// ステップ2: 初期設定ファイル作成
	err = createOrEditConfig(privateKey, "", "", "", "")
	if err != nil {
		errMessage := fmt.Sprintf("Failed to create/edit config: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}
	// ステップ3: WebSocket 通信を確立して情報を取得

	getWebSocketData(publicKey)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to retrieve data via WebSocket: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	// テスト用のダミーデータの挿入
	//cvIP, svIP, ServerPubKey, domainName = "192.168.1.10", "169.254.253.253", "testcodeKey", "otdm.dev"

	// ステップ4: 取得した情報を設定ファイルに追記
	err = createOrEditConfig(privateKey, cvIP, svIP, ServerPubKey, domainName)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to update config with received data: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	err = LogMessage(INFO, "Configuration setup completed.")

	// ステップ5: JSONファイルを/tmpに作成
	err = createStatusJSON(cvIP, svIP, publicKey, ServerPubKey, domainName)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to create status JSON file: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	// 処理終了時ログ
	err = LogMessage(INFO, "websocket.go done")
	if err != nil {
		errMessage := fmt.Sprintf("Failed to log message: %v\n", err)
		err = LogMessage(ERRO, errMessage)
	}

	return cvIP, svIP, ServerPubKey, domainName, nil
}

// 鍵を生成する関数
func generateKeys() (privateKey, publicKey string, err error) {
	// 鍵生成のコマンド実行
	privCmd := exec.Command("wg", "genkey")
	privKey, err := privCmd.Output()
	if err != nil {
		errMessage := fmt.Sprintf("failed to generate private key: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", err
	}

	// プライベート鍵のトリム処理をして不正な文字を排除
	trimmedPrivKey := strings.TrimSpace(string(privKey))

	// 公開鍵生成のためのコマンド
	pubCmd := exec.Command("bash", "-c", fmt.Sprintf("echo %s | wg pubkey", trimmedPrivKey))
	pubKeyOutput, err := pubCmd.CombinedOutput()
	if err != nil {
		// 出力されたエラーメッセージを含める
		errMessage := fmt.Sprintf("failed to generate public key: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", err
	}

	return strings.TrimSpace(string(privKey)), strings.TrimSpace(string(pubKeyOutput)), nil
}

// getWebSocketData はWebSocketを介してデータを取得
func getWebSocketData(otdmPubKey string) (cvIP, svIP, otdmPubKeyResult, domainName string, err error) {
	err = LogMessage(INFO, "get web socket data start")
	LogMessage(INFO, fmt.Sprintf("get web socket data start\n"))

	c, _, err := websocket.DefaultDialer.Dial(webscketURL, nil)
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("failed to connect to websocket server: %v\n", err))
		errMessage := fmt.Sprintf("failed to connect to websocket server: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}
	defer c.Close()

	// 公開鍵をWebSocketを通じて送信
	err = c.WriteMessage(websocket.TextMessage, []byte(otdmPubKey))
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("failed to send public key: %v\n", err))
		errMessage := fmt.Sprintf("failed to send public key: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	// メッセージを受信
	_, message, err := c.ReadMessage()
	if err != nil {
		LogMessage(INFO, fmt.Sprintf("failed to read message: %v\n", err))
		errMessage := fmt.Sprintf("failed to read message: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	// メッセージをJSONとしてデコード
	var response WebSocketResponse
	err = json.Unmarshal(message, &response)
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("failed to unmarshal JSON: %v\n", err))
		errMessage := fmt.Sprintf("failed to unmarshal JSON: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return "", "", "", "", err
	}

	// 必要な値を変数に代入
	cvIP = response.ClientVirtualIP
	svIP = response.ServerVirtualIP
	otdmPubKeyResult = response.ServerPublicKey
	domainName = response.Subdomain

	LogMessage(INFO, fmt.Sprintf("Received data: cvIP=%s, svIP=%s, otdmPubKey=%s, domainName=%s\n", cvIP, svIP, otdmPubKeyResult, domainName))

	// 追加でステータスメッセージを受信する場合の処理（省略可能）
	_, statusMessage, err := c.ReadMessage()
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("failed to read status message: %v\n", err))
		errMessage := fmt.Sprintf("failed to read status message: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return cvIP, svIP, otdmPubKeyResult, domainName, nil // 応急的に部分的な成功とする
	}

	var status WebSocketMessage
	err = json.Unmarshal(statusMessage, &status)
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("failed to unmarshal status JSON: %v\n", err))
		errMessage := fmt.Sprintf("failed to unmarshal status JSON: %v\n", err)
		err = LogMessage(ERRO, errMessage)
	} else {
		LogMessage(INFO, fmt.Sprintf("Received status: %s (%s)\n", status.Message, status.Status))
	}

	return cvIP, svIP, otdmPubKeyResult, domainName, nil
}

// otdm.confを必要なら生成または編集する
func createOrEditConfig(privateKey, cvIP, svIP, ServerPubKey, domainName string) error {
	// 設定ファイルのパスを /etc/wireguard/ に固定
	configPath := filepath.Join("/etc/wireguard", "otdm.conf")

	// /etc/wireguard ディレクトリが存在するか確認し、なければ作成
	if _, err := os.Stat("/etc/wireguard"); os.IsNotExist(err) {
		if err := os.Mkdir("/etc/wireguard", 0755); err != nil {
			errMessage := fmt.Sprintf("failed to create directory /etc/wireguard: %v\n", err)
			err = LogMessage(ERRO, errMessage)
			return err
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
`, privateKey, cvIP, ServerPubKey, domainName, svIP)

	// ファイルが存在しなくても、初期化したい場合でも一貫してテンプレートで上書き
	return ioutil.WriteFile(configPath, []byte(template), 0644)
}

// createStatusJSON は取得したデータをJSON形式で/tmpに保存
func createStatusJSON(cvIP, svIP, publicKey, ServerPubKey, domainName string) error {
	fileName := filepath.Join("/tmp", fmt.Sprintf("otdm_status.json"))

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData := fmt.Sprintf(`{
        "ClientIP": "%s",
        "ServerIP": "%s",
        "publicKey": "%s",
        "otdmPublicKey": "%s",
        "Domain": "%s"
    }`, cvIP, svIP, publicKey, ServerPubKey, domainName)

	_, err = file.WriteString(jsonData)
	return err
}
