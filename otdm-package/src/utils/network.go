package utils

import (
	"fmt"
	"os/exec"
)

// ConfigureFirewall ファイアウォールルールを設定する関数
func ConfigureFirewall(interfaceName, localIP, remoteIP string, httpPort int) error {
	// 起動時ログ
	err := LogMessage(INFO, "network.go start")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
		return err
	}

	// ステップ1: ICMP（ping）通信を許可
	err = exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-p", "icmp", "-s", remoteIP, "-d", localIP, "-j", "ACCEPT").Run()
	if err != nil {
		return fmt.Errorf("Failed to set iptables rule for ICMP: %v", err)
	}

	// ステップ2: HTTP通信（80番ポート）を許可
	err = exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-p", "tcp", "--dport", "80", "-s", remoteIP, "-d", localIP, "-j", "ACCEPT").Run()
	if err != nil {
		return fmt.Errorf("Failed to set iptables rule for HTTP: %v", err)
	}

	// ステップ3: HTTPポートフォワーディングを設定（必要に応じて）
	err = exec.Command("sudo", "iptables", "-t", "nat", "-A", "PREROUTING", "-i", interfaceName, "-p", "tcp", "--dport", "80", "-j", "REDIRECT", "--to-port", fmt.Sprintf("%d", httpPort)).Run()
	if err != nil {
		return fmt.Errorf("Failed to set iptables rule for HTTP port forwarding: %v", err)
	}

	// ステップ4: 他の通信を拒否
	err = exec.Command("sudo", "iptables", "-A", "INPUT", "-i", interfaceName, "-d", localIP, "-j", "DROP").Run()
	if err != nil {
		return fmt.Errorf("Failed to set iptables rule to drop other traffic: %v", err)
	}

	// 成功ログ
	fmt.Println("Firewall rules applied successfully.")
	err = LogMessage(INFO, "network.go done")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}

	return nil
}
