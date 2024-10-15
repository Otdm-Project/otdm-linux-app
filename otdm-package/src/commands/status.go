package commands

import "fmt"

// ShowStatus 関数で現在の状態を表示
func ShowStatus(cvIP, svIP, domainName string) {
	// トンネルの接続状況を取得
	tunnelStatus := checkTunnelStatus()

	fmt.Printf("Client IP: %s\n", cvIP)
	fmt.Printf("Server IP: %s\n", svIP)
	fmt.Printf("Domain: %s\n", domainName)
	fmt.Printf("Tunnel Status: %s\n", tunnelStatus)
}

// checkTunnelStatus ではダミーでトンネル状態を返す。実際の実装では、具体的に確認する手段を持つ。
func checkTunnelStatus() string {
	// ここでは仮に接続されている状態にする。実際の接続状況を取得するロジックを追加可能。
	return "connection"
}