# JumpServer Go SDK（JMS + PAM）

一个面向 JumpServer 的 Go 开发包，封装了常用的用户、节点、资产、授权关系（权限）查询与维护，以及与 PAM 的账号密码检出能力。内置 HTTP Signature 认证与简单日志输出，便于在自动化流程中批量创建/更新授权，申请临时 ROOT 等。

**模块名**：`github.com/lfz97/jumpserver`  ·  **Go 版本**：1.25+

## 安装

```bash
go get github.com/lfz97/jumpserver@latest
```

## 快速开始

在 JumpServer 创建两组 API Key（JMS 与 PAM），然后通过 `Init()` 初始化客户端：

```go
package main

import (
	"log"
	"github.com/lfz97/jumpserver"
	"github.com/lfz97/jumpserver/service"
)

func main() {
	url := "https://jumpserver.example.com"
	jmsID, jmsSecret := "your-jms-api-id", "your-jms-api-secret"
	pamID, pamSecret := "your-pam-api-id", "your-pam-api-secret"
	logfile := "jms-sdk.log"

	client, err := jumpserver.Init(url, jmsID, jmsSecret, pamID, pamSecret, logfile)
	if err != nil {
		log.Fatal(err)
	}

	// 查询用户
	users, err := client.GetUserByName("alice")
	if err != nil {
		client.Logger_p.Println("查询用户失败:", err)
	} else {
		client.Logger_p.Println("用户数量:", len(*users))
	}

	// 申请新授权：为多个用户在多个节点创建/复用同名授权并加入用户
	err = service.RequestNewPermission(client, []string{"alice", "bob"}, []string{"组织/业务/环境/节点A", "组织/业务/环境/节点B"})
	if err != nil {
		client.Logger_p.Println("申请新授权失败:", err)
	}
}
```

## 能力一览

- **用户**：按用户名查询用户信息（精确匹配/可能返回多条）。
- **节点**：按 `full_value` 查询节点与子节点；按节点 ID 列出资产（可选是否仅当前节点）。
- **资产**：查询全部资产、按 IP/名称查询资产。
- **授权关系**：按名称查询、按 ID 获取详情、创建空授权关系、根据配置更新授权关系。
- **PAM 密码**：按资产 ID + 账号名检出密码（带 `X-Source: jms-pam`）。
- **服务流程**：
  - `RequestNewPermission`：校验用户与节点，自动创建同名授权并加入用户。
  - `RequestRootPermission`：对用户已授权范围内的目标资产申请临时 ROOT（按天数设置失效）。
  - `CheckoutPassword`：对用户已授权的目标资产批量检出指定账号密码。

## 常用示例

### 新授权（按节点同名授权）

```go
// 复用/创建与节点同名的授权关系，并把用户加入其中
err := service.RequestNewPermission(client,
	[]string{"alice", "bob"},
	[]string{"组织/业务/环境/节点A", "组织/业务/环境/节点B"},
)
if err != nil { client.Logger_p.Println(err) }
```

说明：节点匹配使用节点的 `full_value`（UI 中的完整路径），当不存在同名授权时会按默认模板自动创建，并配置账号策略与协议。

### 申请临时 ROOT

```go
// 在用户既有授权范围内，对目标资产申请 N 天有效的 ROOT
err := service.RequestRootPermission(client,
	"alice",
	[]string{"组织/业务/环境/节点A"}, // 展开为节点下资产
	[]string{"10.0.0.1", "10.0.0.2"}, // 也可直接指定 IP
	3, // 授权天数
)
if err != nil { client.Logger_p.Println(err) }
```

规则：
- 仅对“用户已授权”范围内的资产申请 ROOT；不在授权范围内的资产将被拒绝并打印到日志。
- 会自动创建以 `AUTO_ROOT_<user>_<uuid>` 命名的授权，并设置到期时间（UTC）。

### 批量检出密码（PAM）

```go
secrets, err := service.CheckoutPassword(client,
	"alice",
	[]string{"服务器一", "服务器二"},
	[]string{"root", "dominos-user"},
)
if err != nil { client.Logger_p.Println(err) }

for _, info := range secrets {
	client.Logger_p.Printf("%s(%s)\n", info.AssetName, info.AssetAddress)
	for _, s := range info.Secrets {
		client.Logger_p.Printf("  %s: %s\n", s.Account, s.Password)
	}
}
```

返回结构体：`[]serviceModel.SecretInfo`，每个元素包含资产基础信息与账号密码列表。

## 认证与日志

- **认证**：基于 HTTP Signature（`hmac-sha256`），参与签名的头为 `date` 与 `(request-target)`；PAM 请求额外带 `X-Source: jms-pam`。
- **初始化**：参见 [init.go](init.go)，内部为 JMS 与 PAM 各维护一个 `resty.Client`，均启用 `SetDebug(true)`。
- **日志**：参见 [mylogger/mylogger.go](mylogger/mylogger.go)，日志以“文件 + 标准输出”的方式写入，调用 `Init()` 传入日志文件路径即可。

## 错误与返回值

- 所有方法以 `error` 表示失败原因，并在非 2xx/预期状态码时返回明确的中文错误信息。
- 典型返回：
  - `GetUserByName()` → `*models.UserList`
  - `GetAssetNodeByFullValue()` → `*models.AssetNodeList`
  - `GetAssetsByNodeID()` / `GetAllAssets()` / `GetAssetByIP()` / `GetAssetByName()` → `*models.AssetsListResult`
  - `GetAssetPermissionByName()` → `*models.PermissionList`
  - `GetAssetPermissionDetailByID()` / `CreateEmptyPermission()` / `UpdatePermission()` → `*models.Permission`
  - `GetSecret()` → `*models.Secret`

## 兼容性与依赖

- Go 1.25+（参见 [go.mod](go.mod)）
- 主要依赖：
  - `github.com/go-resty/resty/v2`
  - `github.com/google/uuid`
  - `gopkg.in/twindagger/httpsig.v1`

## 本地构建

在模块根目录执行：

```powershell
go build ./...
```

## 目录结构

- SDK 入口与初始化：[init.go](init.go)
- 业务函数：
  - JMS API：[functions/JMSfunctions.go](functions/JMSfunctions.go)
  - PAM API：[functions/PAMfunctions.go](functions/PAMfunctions.go)
  - HTTP 签名与客户端定义：[functions/SigAuth.go](functions/SigAuth.go)
- 复合业务逻辑：[logic/logic.go](logic/logic.go)
- 服务流程封装：[service/service.go](service/service.go)
- 模型定义：[models/*](models)
- 工具与日志：[utils/parseUrl.go](utils/parseUrl.go)，[mylogger/mylogger.go](mylogger/mylogger.go)

—— 欢迎按需扩展，如需更多示例或封装，可提出 Issue 或直接在项目中新增方法。
