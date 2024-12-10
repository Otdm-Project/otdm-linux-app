package utils

import (
	"fmt"
	"os/exec"
	"time"
)

// CallWatchman は、トンネルが維持されているかを監視
func CallWatchman(serverIP string) {
	err := LogMessage(INFO, "watchman.go start")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}
	pinger.SetPrivileged(true)
	ticker := time.NewTicker(5 * time.Minute) // 5分ごとに監視サイクルを開始
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking tunnel health...")

			// 最大三回の再送回数のカウンタ
			attempts := 3
			success := false

			for i := 0; i < attempts; i++ {
				fmt.Printf("Ping attempt %d...\n", i+1)

				// pingコマンドを使用して応答を確認
				err := exec.Command("ping", "-c", "4", "-w", "4", serverIP).Run()

				if err == nil {
					fmt.Println("Tunnel is healthy.")
					success = true
					break
				} else {
					fmt.Printf("Ping attempt %d failed: %v\n", i+1, err)
				}
			}

			if !success {
				fmt.Println("Tunnel is deemed unhealthy after three attempts.")

				// logger.goに異常を伝える処理
				err := LogMessage(ERRO, "Tunnel is unhealthy. Attempting to reset.")
				if err != nil {
					fmt.Printf("Failed to log message: %v\n", err)
				}
				// 異常と判断したのでトンネルを再試行、現段階
				resetTunnel()
			}

		}
	}
	err = LogMessage(INFO, "watchman.go done")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}
}

// resetTunnel 関数はトンネルを再起動します
func resetTunnel() {
	fmt.Println("Resetting the tunnel...")
	exec.Command("otdm", "down").Run()
	exec.Command("otdm", "up").Run()
}
