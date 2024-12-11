package utils

import (
	"fmt"
	"os/exec"
	"time"

	// "otdm-package/src/utils"
	probing "github.com/prometheus-community/pro-bing"
)

// sendPing は指定したサーバーIPに対してICMP通信を行う
func SendPing(serverIP string) bool {
	// pinger インスタンスを作成
	pinger, err := probing.NewPinger(serverIP)
	if err != nil {
		fmt.Printf("Failed to create pinger: %v\n", err)
		return false
	}

	// ICMP送信の特権モードで有効化(Linuxで必要)
	pinger.SetPrivileged(true)

	// 送信パケット数とタイムアウトを設定
	pinger.Count = 4
	pinger.Timeout = 10 * time.Second

	// Ping 実行
	fmt.Printf("Pinging server: %s\n", serverIP)
	err = pinger.Run()
	if err != nil {
		fmt.Printf("Ping failed: %v\n", err)
		return false
	}

	// 結果を取得
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		fmt.Println("No packets received, server unreachable.")
		return false
	}

	fmt.Printf("Ping successful: %d/%d packets received.\n", stats.PacketsRecv, stats.PacketsSent)

	return true
}

// CallWatchman はサーバーとの初回ハンドシェイクとトンネル維持監視を行います
func CallWatchman(serverIP string) {
	fmt.Println("Starting Watchman...")

	// 初回ハンドシェイクを試みる（最大3回の再送）
	for i := 0; i < 3; i++ {
		if SendPing(serverIP) {
			fmt.Println("Handshake with server successful.")
			break
		} else if i == 2 { // 3回目も失敗した場合
			fmt.Println("Failed to establish handshake with server after 3 attempts. Exiting.")
			return
		}
		fmt.Printf("Retrying handshake (%d/3)...\n", i+2)
		time.Sleep(5 * time.Second)
	}

	// 監視ループ開始
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking server health...")
			if !SendPing(serverIP) {
				fmt.Println("Server unreachable. Attempting to reset tunnel.")
				resetTunnel()
			} else {
				fmt.Println("Server is healthy.")
			}
		}
	}
}

// resetTunnel 関数はトンネルを再起動します
func resetTunnel() {
	fmt.Println("Resetting the tunnel...")
	err := exec.Command("otdm", "down").Run()
	if err != nil {
		fmt.Printf("Failed to execute 'otdm down': %v\n", err)
		LogMessage(ERRO, fmt.Sprintf("Failed to execute 'otdm down': %v", err))
	}
	err = exec.Command("otdm", "up").Run()
	if err != nil {
		fmt.Printf("Failed to execute 'otdm up': %v\n", err)
		LogMessage(ERRO, fmt.Sprintf("Failed to execute 'otdm up': %v", err))
	}
	exec.Command("otdm", "down").Run()
	exec.Command("otdm", "up").Run()
}
