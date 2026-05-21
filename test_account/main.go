package main

import (
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

	// 搜索 test_delete 授权关系
	fmt.Println("========== 搜索 test_delete ==========")
	permissions, err := client.GetALLAssetPermissions(map[string]string{
		"name": "test_delete",
	})
	if err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
		os.Exit(1)
	}

	if len(*permissions) == 0 {
		fmt.Println("❌ 未找到 test_delete 授权关系")
		os.Exit(1)
	}

	target := (*permissions)[0]
	fmt.Printf("找到: ID=%s  Name=%s\n\n", target.ID, target.Name)

	// 删除
	fmt.Println("========== 删除 test_delete ==========")
	err = client.DeletePermission(target.ID)
	if err != nil {
		fmt.Printf("❌ 删除失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 已删除: %s (%s)\n", target.Name, target.ID)
}
