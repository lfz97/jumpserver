package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lfz97/jumpserver"
	"github.com/lfz97/jumpserver/models"
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

	// 生成随机后缀避免用户名冲突
	rand.Seed(time.Now().UnixNano())
	suffix := fmt.Sprintf("%04d", rand.Intn(10000))
	testName := fmt.Sprintf("fan.zhou_%s", suffix)
	testUsername := fmt.Sprintf("v-fan.zhou_%s", suffix)
	testEmail := fmt.Sprintf("v-fan.zhou_%s@freemud.com", suffix)

	mfaLevel := 0
	req := models.CreateUserRequest{
		Name:             testName,
		Username:         testUsername,
		Email:            testEmail,
		PasswordStrategy: "email",
		Source:           "radius",
		MfaLevel:         &mfaLevel,
		DateExpired:      "2096-01-17T15:32:11.000+0800",
		Comment:          "freemud",
		SystemRoles: []models.RoleParam{
			{PK: "00000000-0000-0000-0000-000000000003"},
		},
		OrgRoles: []models.RoleParam{
			{PK: "00000000-0000-0000-0000-000000000007"},
		},
	}

	fmt.Printf("测试用户名: %s / %s / %s\n", testName, testUsername, testEmail)

	fmt.Println("========== 创建用户 ==========")
	user, err := client.CreateUser(req)
	if err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		os.Exit(1)
	}

	// 美化输出
	out, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("✅ 创建成功:\n%s\n", string(out))
}
