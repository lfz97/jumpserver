package main

import (
	"fmt"
	"os"

	"github.com/lfz97/jumpserver"
	"github.com/lfz97/jumpserver/service"
)

func main() {
	// ========== 配置 ==========
	url := ""
	jmsApiID := ""
	jmsApiSecret := ""
	pamApiID := ""
	pamApiSecret := ""
	// ==========================

	client, err := jumpserver.Init(url, jmsApiID, jmsApiSecret, pamApiID, pamApiSecret, "./test_account.log")
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		os.Exit(1)
	}

	// ================== 集成测试: ChangeSecret ==================
	fmt.Println("========== 集成测试: ChangeSecret ==========")

	items := []service.ChangeSecretItem{
		// 资产1: 76f416af-044a-458e-85bf-cf375e67779b
		{AssetID: "76f416af-044a-458e-85bf-cf375e67779b", Account: "Bingo"},
		{AssetID: "76f416af-044a-458e-85bf-cf375e67779b", Account: "Eva.yin"},
		{AssetID: "76f416af-044a-458e-85bf-cf375e67779b", Account: "Joly"},
		// 资产2: a21ac621-fbb3-46c4-9d3c-ec31e152b35a
		{AssetID: "a21ac621-fbb3-46c4-9d3c-ec31e152b35a", Account: "dominos-tedia-readonly"},
		{AssetID: "a21ac621-fbb3-46c4-9d3c-ec31e152b35a", Account: "dominos-ops"},
		{AssetID: "a21ac621-fbb3-46c4-9d3c-ec31e152b35a", Account: "dominos_zfr"},
	}

	service.ChangeSecret(client, items)

	fmt.Println("✅ 改密完成！")
}
