package utils

import (
    "fmt"
    "os/exec"
)

// WireGuard トンネルの起動を行う関数
func CallBoot() error {
    // ステップ 1: wg-quick を使用してトンネルを起動
    err := exec.Command("sudo", "wg-quick", "up", "config/otdm.conf").Run()
    if err != nil {
        return fmt.Errorf("Failed to execute wg-quick up: %v", err)
    }

    // ステップ 2: systemctl を用いてトンネルをシステムサービスとして有効化
    err = exec.Command("sudo", "systemctl", "enable", "wg-quick@otdm.conf").Run()
    if err != nil {
        return fmt.Errorf("Failed to enable wg-quick service: %v", err)
    }

    fmt.Println("WireGuard tunnel is up and enabled as a service.")
    return nil
}