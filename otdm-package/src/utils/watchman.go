package utils

import (
	"fmt"
	"os/exec"
	"time"

	// "otdm-package/src/utils"
	probing "github.com/prometheus-community/pro-bing"
)

type messages struct {
	message    string
	errmessage string
}

// sendPing は指定したサーバーIPに対してICMP通信を行う
func SendPing(serverIP string) bool {
	// pinger インスタンスを作成
	pinger, err := probing.NewPinger(serverIP)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to create pinger: %v\n", err)
		LogMessage(ERRO, errMessage)
		return false
	}

	// ICMP送信の特権モードで有効化(Linuxで必要)
	pinger.SetPrivileged(true)

	// 送信パケット数とタイムアウトを設定
	pinger.Count = 4
	pinger.Timeout = 10 * time.Second

	// Ping 実行
	//Message := fmt.Printf("Pinging server: %v\n", serverIP)
	LogMessage(INFO, fmt.Sprintf("Starting watchman for server: %s", serverIP))
	err = LogMessage(INFO, "Pinging server")

	err = pinger.Run()
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("Ping failed: %v\n", err))
		return false
	}

	// 結果を取得
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		err = LogMessage(ERRO, "No packets received, server unreachable.")
		return false
	}

	LogMessage(INFO, fmt.Sprintf("Ping successful: %d/%d packets received.\n", stats.PacketsRecv, stats.PacketsSent))

	return true
}

// CallWatchman はサーバーとの初回ハンドシェイクとトンネル維持監視を行います
func CallWatchman(serverIP string) {
	LogMessage(INFO, "watchman.go start")

	// 初回ハンドシェイクを試みる（最大3回の再送）
	for i := 0; i < 3; i++ {
		if SendPing(serverIP) {
			LogMessage(INFO, "Handshake with server successful.")
			break
		} else if i == 2 { // 3回目も失敗した場合
			LogMessage(INFO, "Failed to establish handshake with server after 3 attempts. Exiting.")
			return
		}
		LogMessage(INFO, fmt.Sprintf("Retrying handshake (%d/3)...\n", i+2))
		time.Sleep(5 * time.Second)
	}

	// 監視ループ開始
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			LogMessage(INFO, "Checking server health...")
			if !SendPing(serverIP) {
				LogMessage(ERRO, "Server unreachable. Attempting to reset tunnel.")
				resetTunnel()
			} else {
				LogMessage(INFO, "Server is healthy.")
			}
		}
	}
}

// resetTunnel 関数はトンネルを再起動します
func resetTunnel() {
	LogMessage(INFO, "Resetting the tunnel...")
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
