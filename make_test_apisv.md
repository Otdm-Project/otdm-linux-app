main.rsを以下のように変更して作成している為注意
```rust
// getWebSocketData はWebSocketを介して平文でデータを取得
func getWebSocketData() (cvIP, svIP, otdmPubKey, domainName string, err error) {
    // WebSocket サーバーのURL
    url := "ws://10.2.141.46:3000"

    // WebSocket接続の確立
    c, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return "", "", "", "", fmt.Errorf("failed to connect to websocket server: %v", err)
    }
    defer c.Close()

    // メッセージの受信
    _, message, err := c.ReadMessage()
    if err != nil {
        return "", "", "", "", fmt.Errorf("failed to read message: %v", err)
    }

    // 平文メッセージを分割
    parts := strings.Split(string(message), ",")
    if len(parts) != 4 {
        return "", "", "", "", fmt.Errorf("received message is not valid")
    }

    return parts[0], parts[1], parts[2], parts[3], nil
}
```