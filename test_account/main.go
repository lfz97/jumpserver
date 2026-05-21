package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lfz97/jumpserver"
)

func main() {
	url := ""
	jmsApiID := ""
	jmsApiSecret := ""
	pamApiID := ""
	pamApiSecret := ""

	client, err := jumpserver.Init(url, jmsApiID, jmsApiSecret, pamApiID, pamApiSecret, "./test_account.log")
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		os.Exit(1)
	}

	baseURL := url + "/api/v1/perms/asset-permissions/"

	// 加 limit=1 强制分页对象返回
	tests := []string{
		"?limit=1&offset=0",
		"?limit=1&offset=0&is_expired=true",
		"?limit=1&offset=0&is_expired=false",
	}

	for _, q := range tests {
		resp, _ := client.JMSClient_p.R().Get(baseURL + q)

		firstChar := resp.Body()[0]
		if firstChar == '[' {
			var arr []json.RawMessage
			json.Unmarshal(resp.Body(), &arr)
			fmt.Printf("%-45s  格式=数组  条数=%d\n", q, len(arr))
		} else {
			var result struct {
				Count int `json:"count"`
			}
			json.Unmarshal(resp.Body(), &result)
			fmt.Printf("%-45s  格式=对象  count=%d\n", q, result.Count)
		}
	}
}
