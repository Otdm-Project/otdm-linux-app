package utils

import (
	"fmt"
        "io/ioutil"
	"os"
	"os/exec"
)

// CallRefresh は、設定ファイルの存在をチェックし、
// その内容を削除した上で、VPNサービスを停止/無効化します。

func CallRefresh() {
	configPath := "config/otdm.conf"  // 設定ファイルのパスを変数に格納
	// ステップ1: 設定ファイルが存在するかを確認
	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("Config file exists. Clearing contents.")
		// Configファイルが存在する場合、内容をクリア

		// ステップ2: ファイルの内容を空にする
		err = ioutil.WriteFile(configPath, []byte{}, 0644)
		if err != nil{
			fmt.Printf("Failed to clear config file: %v\n", err)
			return
		}

		// ステップ3: トンネルを停止し、サービスを無効化するコマンドを実行
		err = stopAndDisableTunnel()
		if err != nil {
			fmt.Printf("Failed to stop or disable the tunnel: %v\n", err)
	                return
		}
		fmt.Println("Tunnel stopped and service disabled.")
	} else {
		fmt.Println("Config file does not exist or cannot be accessed.")
	}
}
		// stopAndDisableTunnel は WireGuard トンネルを停止し、無効化するためのコマンドを実行します
		func stopAndDisableTunnel() error {
			// コマンド1: wg-quick を使ってトンネルを停止 (sudo 権限が必要)
			if err := exec.Command("sudo", "wg-quick", "down", "config/otdm.conf").Run(); err != nil {
				return fmt.Errorf("failed to execute wg-quick down: %v", err)
			}
		// コマンド2: systemctl を使ってサービスを無効化 (sudo 権限が必要)
		if err := exec.Command("sudo", "systemctl", "disable", "wg-quick@otdm.conf").Run(); err != nil {
			return fmt.Errorf("failed to disable wg-quick service: %v", err)
		}

		return nil
}
