package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

// LogLevel はログレベルを表す
type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERRO
)

// logLevelMap は LogLevel を文字列として表現するためのマップ
var logLevelMap = map[LogLevel]string{
	INFO: "INFO",
	WARN: "WARN",
	ERRO: "ERRO",
}

// logFilePath は絶対パスで /var/log に設定
const logFilePath = "/var/log/otdm-package.log"
const maxLogLines = 1000

// LogMessage はメッセージとログレベルを受け取り、ログファイルに記録する
func LogMessage(logLevel LogLevel, message string) error {
	// 現在の時間と実行ユーザーを取得
	currentTime := time.Now().Format("Jan 02 2006 15:04:05")
	user, err := user.Current()
	if err != nil {
		errMessage := fmt.Sprintf("unable to get current user: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return err
	}

	// ログメッセージを構築
	logEntry := fmt.Sprintf("%s %s %s : %s", currentTime, logLevelMap[logLevel], user.Username, message)

	// ログファイルを開く（追記モード）
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		errMessage := fmt.Sprintf("unable to open log file: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return err
	}
	defer file.Close()

	// 既存ログ行数を数える
	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	// 古いログを削除して最大行数を維持
	if lines >= maxLogLines {
		if err := rotateLogs(); err != nil {
			errMessage := fmt.Sprintf("unable to rotate logs: %v\n", err)
			err = LogMessage(ERRO, errMessage)
			return err
		}
	}

	// 新規ログをファイルに書き込む
	_, err = file.WriteString(logEntry + "\n")
	if err != nil {
		errMessage := fmt.Sprintf("unable to write to log file: %v\n", err)
		err = LogMessage(ERRO, errMessage)
		return err
	}

	return nil
}

func rotateLogs() error {
	input, err := os.ReadFile(logFilePath)
	if err != nil {
		return fmt.Errorf("unable to read log file: %v", err)
	}

	lines := strings.Split(string(input), "\n")
	if len(lines) <= maxLogLines {
		return nil
	}

	// 最新1000行を保持し、古い行を削除
	newContent := strings.Join(lines[len(lines)-maxLogLines:], "\n")
	return os.WriteFile(logFilePath, []byte(newContent), 0644)
}

func ErrLogMessage(errMessage string) error {
	cmd := exec.Command("logger", "-p", "otdm-package", errMessage)
	fmt.Printf("journal code boot done:\n")
	if err := cmd.Run(); err != nil {
		fmt.Printf("failed to send to journal: %v\n", err)
	}
	return nil
}
