package commands

import (
	"fmt"
	"os/exec"
	"otdm-package/src/utils"
)

// RunDown は、WireGuard トンネルを停止し、無効化する
func RunDown() error {
	utils.LogMessage(utils.INFO, "down.go start")
	// ステップ1: wg-quick を使用してトンネルをダウン
	err := exec.Command("sudo", "wg-quick", "down", "etc/wireguard/otdm.conf").Run()
	if err != nil {
		errMessage := fmt.Sprintf("Failed to execute wg-quick down:  %v", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return nil
	}

	// ステップ2: systemctl を用いてトンネルのサービスを無効化
	err = exec.Command("sudo", "systemctl", "disable", "wg-quick@otdm.conf").Run()
	if err != nil {
		errMessage := fmt.Sprintf("Failed to disable wg-quick service: %v", err)
		utils.LogMessage(utils.ERRO, errMessage)
		return nil
	}

	fmt.Println("WireGuard tunnel is down and disabled.")
	utils.LogMessage(utils.INFO, "down.go start")
	return nil
}
