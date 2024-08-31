package commands

import (
    "fmt"
    "os/exec"
)

// RunDown は、WireGuard トンネルを停止し、無効化する
func RunDown() error {
    // ステップ1: wg-quick を使用してトンネルをダウン
    err := exec.Command("sudo", "wg-quick", "down", "config/otdm.conf").Run()
    if err != nil {
        return fmt.Errorf("Failed to execute wg-quick down: %v", err)
    }

    // ステップ2: systemctl を用いてトンネルのサービスを無効化
    err = exec.Command("sudo", "systemctl", "disable", "wg-quick@otdm.conf").Run()
    if err != nil {
        return fmt.Errorf("Failed to disable wg-quick service: %v", err)
    }

    fmt.Println("WireGuard tunnel is down and disabled.")
    return nil
}
