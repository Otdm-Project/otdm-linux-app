package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	ping "github.com/prometheus-community/pro-bing"
)

// Watchping はサーバーへの ICMP ping を送信し、到達可能かを確認する
func CallWatchman(svip string) bool {
	LogMessage(INFO, "watchman.go-WatchPing start")
	defer LogMessage(INFO, "watchman.go-WatchPing done")

	pinger, err := ping.NewPinger(svip)
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("Ping instance creation failed: %v", err))
		return false
	}
	pinger.SetPrivileged(true)
	pinger.Count = 5
	pinger.Timeout = 10 * time.Second

	err = pinger.Run()
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("Ping execution failed: %v", err))
		return false
	}

	stats := pinger.Statistics()
	if stats.PacketsSent-stats.PacketsRecv >= 3 {
		LogMessage(WARN, "watchman.go-WatchPing: Unstable")
		return false
	}

	LogMessage(INFO, "watchman.go-WatchPing: Stable")
	return true
}

// handleSocket は "down" を受信したら終了し、"status" を受信したら応答する
func handleSocket() {
	LogMessage(INFO, "watchman.go-HandleSocket start")
	defer LogMessage(INFO, "watchman.go-HandleSocket done")

	listener, err := net.Listen("tcp", "127.0.0.1:3563")
	if err != nil {
		LogMessage(ERRO, fmt.Sprintf("Failed to start TCP listener: %v", err))
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogMessage(ERRO, fmt.Sprintf("Socket accept failed: %v", err))
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				LogMessage(ERRO, fmt.Sprintf("Socket read error: %v", err))
				return
			}

			message := string(buf[:n])
			switch message {
			case "down":
				LogMessage(INFO, "watchman.go-HandleSocket: Received 'down', exiting")
				c.Write([]byte("done"))
				os.Exit(0)
			case "status":
				LogMessage(INFO, "watchman.go-HandleSocket: Received 'status', responding 'active'")
				c.Write([]byte("active"))
			default:
				LogMessage(WARN, "watchman.go-HandleSocket: Unknown command received")
			}
		}(conn)
	}
}

// WatchMan は定期的に ping を実行し、接続が切れた場合に処理を行う
func WatchMan(svip string) {
	LogMessage(INFO, "watchman.go-WatchTunnel start")
	defer LogMessage(INFO, "watchman.go-WatchTunnel stop")

	if !CallWatchman(svip) {
		LogMessage(ERRO, "watchman.go-WatchTunnel: Initial ping failed")
		return
	}

	go handleSocket()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !CallWatchman(svip) {
				LogMessage(ERRO, "watchman.go-WatchTunnel: Ping failed, restarting tunnel")
				restartTunnel()
			}
		}
	}
}

// restartTunnel はトンネルを再起動する
func restartTunnel() {
	LogMessage(INFO, "Restarting the tunnel...")
	exec.Command("otdm", "down").Run()
	exec.Command("otdm", "up").Run()
}
