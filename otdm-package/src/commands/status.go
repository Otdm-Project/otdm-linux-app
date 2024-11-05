package commands

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

// StatusData はJSONファイルの内容を保持するための構造体
type StatusData struct {
    CvIP      string `json:"cvIP"`
    SvIP      string `json:"svIP"`
    OtdmPubKey string `json:"otdmPubKey"`
    DomainName string `json:"domainName"`
}

// ShowStatus は /tmp/status.json の内容をCLI上に表示する
func ShowStatus() error {
    jsonFilePath := "/tmp/status.json"

    // JSONファイルを開く
    file, err := os.Open(jsonFilePath)
    if err != nil {
        fmt.Printf("Failed to open JSON file: %v\n", err)
        return err
    }
    defer file.Close()

    // JSONファイルの内容を読み込む
    byteValue, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Printf("Failed to read JSON file: %v\n", err)
        return err
    }

    // JSONデータを構造体にデコードする
    var status StatusData
    err = json.Unmarshal(byteValue, &status)
    if err != nil {
        fmt.Printf("Failed to parse JSON file: %v\n", err)
        return err
    }

    // JSONデータを表示する
    fmt.Println("Current Status:")
    fmt.Printf("Client Virtual IP (cvIP): %s\n", status.CvIP)
    fmt.Printf("Server IP (svIP): %s\n", status.SvIP)
    fmt.Printf("OTDM Public Key: %s\n", status.OtdmPubKey)
    fmt.Printf("Domain Name: %s\n", status.DomainName)

    return nil
}