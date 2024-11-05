package commands

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "otdm-package/src/utils"
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
    err := utils.LogMessage(utils.INFO, "staus.go start")
	if err != nil {
		fmt.Printf("Failed to log message: %v\n", err)
	}
    jsonFilePath := filepath.Join("/tmp", fmt.Sprintf("otdm_status.json"))

    // JSONファイルの存在を確認
    if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
        return fmt.Errorf("status file not found: %s", jsonFilePath)
    }

    // ファイルを読み出し
    data, err := ioutil.ReadFile(jsonFilePath)
    if err != nil {
        return fmt.Errorf("failed to read status file: %v", err)
    }

    // JSONデータのパース
    var status map[string]interface{}
    if err := json.Unmarshal(data, &status); err != nil {
        return fmt.Errorf("failed to parse JSON: %v", err)
    }

    // CLI上に表示
    fmt.Println("Current Status:")
    for key, value := range status {
        fmt.Printf("%s: %v\n", key, value)
    }
    err = utils.LogMessage(utils.INFO, "staus.go done")
    return nil
}