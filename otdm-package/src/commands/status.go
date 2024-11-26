package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"otdm-package/src/utils"
	"path/filepath"
)

// StatusData はJSONファイルの内容を保持するための構造体
type StatusData struct {
	CvIP       string `json:"ClientIP"`
	SvIP       string `json:"ServerIP"`
	MyPubKey   string `json:"publicKey"`
	OtdmPubKey string `json:"otdmPublicKey"`
	DomainName string `json:"Domain"`
}

// ShowStatus は /tmp/status.json の内容とWireGuardの起動状況をCLI上に表示する
func ShowStatus() error {
	// ログ開始
	if err := utils.LogMessage(utils.INFO, "status.go start"); err != nil {
		fmt.Printf("Failed to log start message: %v\n", err)
	}

	// JSONファイルパス
	jsonFilePath := filepath.Join("/tmp", "otdm_status.json")

	// JSONファイルの存在確認
	if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
		return fmt.Errorf("status file not found: %s", jsonFilePath)
	}

	// ファイルを読み込む
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read status file: %v", err)
	}

	// JSONデータを構造体にパース
	var status StatusData
	if err := json.Unmarshal(data, &status); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// CLIにステータスを表示
	fmt.Println("OTDM Current Status:")
	fmt.Printf("Client Virtual IP (cvIP): %s\n", status.CvIP)
	fmt.Printf("Server Virtual IP (svIP): %s\n", status.SvIP)
	fmt.Printf("OTDM Public Key: %s\n", status.OtdmPubKey)
	fmt.Printf("Domain Name: %s\n", status.DomainName)

	// WireGuard起動状況を確認
	wireguardStatus, err := checkWireGuardStatus()
	if err != nil {
		fmt.Printf("WireGuard status check failed: %v\n", err)
	} else {
		fmt.Printf("WireGuard Status: %s\n", wireguardStatus)
	}

	// ログ終了
	if err := utils.LogMessage(utils.INFO, "status.go done"); err != nil {
		fmt.Printf("Failed to log done message: %v\n", err)
	}

	return nil
}

// checkWireGuardStatus は WireGuard が起動しているかを確認する
func checkWireGuardStatus() (string, error) {
	cmd := exec.Command("wg", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "Inactive", nil // エラーの場合は非アクティブとして扱う
	}
	if len(output) == 0 {
		return "Inactive", nil
	}
	return "Active", nil
}
