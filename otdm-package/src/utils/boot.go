package utils

import (
    "fmt"
    "os/exec"
)

type Status struct {
    ClientIP   string `json:"client_ip"`
    ServerIP   string `json:"server_ip"`
    DomainName string `json:"domain_name"`
}

func saveStatusToFile(status Status) error {
    file, err := os.Create("")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(status)
}

// WireGuard トンネルの起動を行う関数
func CallBoot() error {
    // 起動時ログ
    var err error
    err = LogMessage(INFO, "boot.go start")
    if err != nil {
        fmt.Printf("Failed to log message: %v\n", err)
    }

    // ステップ1: wg-quick を使用してトンネルを起動
    err = exec.Command("sudo", "wg-quick", "up", "/etc/wireguard/otdm.conf").Run()
    if err != nil {
        return fmt.Errorf("Failed to execute wg-quick up: %v", err)
    }

    // ステップ2: systemctl を用いてトンネルをシステムサービスとして有効化
    err = exec.Command("sudo", "systemctl", "enable", "wg-quick@otdm.conf").Run()
    if err != nil {
        return fmt.Errorf("Failed to enable wg-quick service: %v", err)
    }

    // 処理終了時ログ
    err = LogMessage(INFO, "boot.go done")
    if err != nil {
        fmt.Printf("Failed to log message: %v\n", err)
    }

    fmt.Println("WireGuard tunnel is up and enabled as a service.")
    return nil
}