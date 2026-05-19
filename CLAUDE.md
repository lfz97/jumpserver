# JumpServer Go SDK - CLAUDE.md

## Commands

```bash
go build ./...              # 构建检查
go run ./test_account/...   # ad-hoc 测试（无正式测试框架）
```

## Architecture

```
init.go          → 入口，创建 JMS+PAM 双 resty 客户端
functions/       → API 层：JMSfunctions.go (核心API) + PAMfunctions.go (密码检出) + SigAuth.go (签名)
logic/           → 逻辑层：校验用户/节点/授权是否存在、创建模板、插入用户
service/         → 流程层：RequestNewPermission / RequestRootPermission / CheckoutPassword
models/          → API 返回模型 + serviceModel/ (业务输出)
utils/ + mylogger/ → URL 拼接 + 双写日志（文件+stdout）
```

## Key Patterns

- **双客户端**：JMS 和 PAM 各用一个 `resty.Client`，认证 key 不同，通过 `SetPreRequestHook` 自动 HMAC-SHA256 签名
- **签名头**：`(request-target)` + `date`，PAM 额外带 `X-Source: jms-pam`
- **resty Debug 常开**：`SetDebug(true)` 在 init 中硬编码，所有请求打印完整 body

## Gotchas

- **`GetSpecifiedAccount` 返回格式**：该 API 返回数组 `[{...}]` 而非分页对象。函数已兼容两种格式（先试数组，回退到分页对象）。新增 account 相关功能时注意
- **`PermissionConfig` 更新**：update 时会将未传的字段置空。务必先从 `GetAssetPermissionDetailByID` 获取现有值再 merge
- **错误信息为中文**：所有 error 返回中文信息，边界判断时注意
- **用户同名问题**：`GetUserByName` 可能返回多条，service 层会检查唯一性

## Testing

无测试框架。测试方式：在项目外或 `test_account/` 下写 `main.go`，`go run` 执行后删掉（API key 是明文）。
