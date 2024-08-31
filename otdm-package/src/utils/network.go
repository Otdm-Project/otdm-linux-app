package utils

import (
    "fmt"
    "os/exec"
)

// Networkファイアウォールの設定を行う関数
func ConfigureFirewall(interfaceName string) error {
    // ステップ1: SSHや既存通信を許可するルールを追加
    if err := exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-j", "ACCEPT").Run(); err != nil {
        return fmt.Errorf("Failed to set iptables rule for existing communication: %v", err)
    }

    // ステップ2: 80番ポート(HTTP)へのアクセスを許可するルールを追加
    if err := exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-p", "tcp", "--dport", "80", "-j", "ACCEPT").Run(); err != nil {
        return fmt.Errorf("Failed to set iptables rule for HTTP: %v", err)
    }

    // ステップ3: 443番ポート(HTTPS)へのアクセスを許可するルールを追加
    if err := exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-p", "tcp", "--dport", "443", "-j", "ACCEPT").Run(); err != nil {
        return fmt.Errorf("Failed to set iptables rule for HTTPS: %v", err)
    }

    // ステップ4: それ以外の通信を拒否
    if err := exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-j", "DROP").Run(); err != nil {
        return fmt.Errorf("Failed to set iptables rule for dropping other traffic: %v", err)
    }

    fmt.Println("Firewall rules applied successfully.")
    return nil
}