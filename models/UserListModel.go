package models

type userInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserList []userInfo

// RoleParam 创建用户时传递角色的参数结构
type RoleParam struct {
	PK string `json:"pk"`
}

// CreateUserRequest 创建用户请求体
type CreateUserRequest struct {
	Name                 string      `json:"name"`                              // 名称（必选）
	Username             string      `json:"username"`                          // 用户名（必选）
	Email                string      `json:"email"`                             // 邮箱（必选）
	Password             string      `json:"password,omitempty"`                // 密码（password_strategy=custom时必选）
	Wechat               string      `json:"wechat,omitempty"`                  // 微信
	Phone                string      `json:"phone,omitempty"`                   // 手机
	Groups               []string    `json:"groups,omitempty"`                  // 用户组id
	NeedUpdatePassword   bool        `json:"need_update_password,omitempty"`    // 是否下一次登录修改密码
	PublicKey             string      `json:"public_key,omitempty"`              // SSH公钥
	SystemRoles          []RoleParam `json:"system_roles"`                      // 系统角色（必选）
	OrgRoles             []RoleParam `json:"org_roles"`                         // 组织角色（必选）
	PasswordStrategy     string      `json:"password_strategy,omitempty"`       // 密码策略：email / custom
	Source               string      `json:"source,omitempty"`                  // 来源：local、ldap、openid 等
	MfaLevel             *int        `json:"mfa_level,omitempty"`               // MFA：0=禁用 1=启用 2=强制启用
	DateExpired          string      `json:"date_expired,omitempty"`            // 用户失效时间，格式：2023-02-04T00:54:39.000Z
	Comment              string      `json:"comment,omitempty"`                 // 备注
}

// UserInfo 用户信息（创建用户接口返回）
type UserInfo struct {
	ID                      string     `json:"id"`
	Name                    string     `json:"name"`
	Username                string     `json:"username"`
	Email                   string     `json:"email"`
	Wechat                  string     `json:"wechat"`
	Phone                   *string    `json:"phone"`
	MfaLevel                LabelValue `json:"mfa_level"`
	Source                  LabelValue `json:"source"`
	WecomID                 *string    `json:"wecom_id"`
	DingtalkID              *string    `json:"dingtalk_id"`
	FeishuID                *string    `json:"feishu_id"`
	CreatedBy               string     `json:"created_by"`
	UpdatedBy               string     `json:"updated_by"`
	Comment                 *string    `json:"comment"`
	AvatarURL               string     `json:"avatar_url"`
	Groups                  []string   `json:"groups"`
	SystemRoles             []RoleInfo `json:"system_roles"`
	OrgRoles                []OrgRole  `json:"org_roles"`
	IsServiceAccount        bool       `json:"is_service_account"`
	IsValid                 bool       `json:"is_valid"`
	IsExpired               bool       `json:"is_expired"`
	IsActive                bool       `json:"is_active"`
	IsOtpSecretKeyBound     bool       `json:"is_otp_secret_key_bound"`
	CanPublicKeyAuth        bool       `json:"can_public_key_auth"`
	MfaEnabled              bool       `json:"mfa_enabled"`
	NeedUpdatePassword      bool       `json:"need_update_password"`
	MfaForceEnabled         bool       `json:"mfa_force_enabled"`
	IsFirstLogin            bool       `json:"is_first_login"`
	LoginBlocked            bool       `json:"login_blocked"`
	DateExpired             string     `json:"date_expired"`
	DateJoined              string     `json:"date_joined"`
	LastLogin               *string    `json:"last_login"`
	DateUpdated             string     `json:"date_updated"`
	DatePasswordLastUpdated string     `json:"date_password_last_updated"`
}

// LabelValue 通用 label/value 结构（mfa_level、source 等字段）
// Value 可能是 int（如 mfa_level）或 string（如 source），用 interface{} 兼容
type LabelValue struct {
	Value interface{} `json:"value"`
	Label string      `json:"label"`
}

// RoleInfo 系统角色信息
type RoleInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// OrgRole 组织角色信息
type OrgRole struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}
