package models

// 创建改密计划请求
type CreateChangeSecretAutomationReq struct {
	Name           string   `json:"name"`
	Accounts       []string `json:"accounts"`
	Assets         []string `json:"assets,omitempty"`
	Nodes          []string `json:"nodes,omitempty"`
	IsActive       bool     `json:"is_active"`
	IsPeriodic     bool     `json:"is_periodic"`
	Crontab        string   `json:"crontab,omitempty"`
	Interval       string   `json:"interval,omitempty"`
	SecretStrategy string   `json:"secret_strategy"`
	SecretType     string   `json:"secret_type"`
	PasswordRules  interface{} `json:"password_rules,omitempty"`
	Secret         string   `json:"secret,omitempty"`
	Comment        string   `json:"comment,omitempty"`
}

// 改密计划响应
type ChangeSecretAutomation struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Accounts       []string      `json:"accounts"`
	Assets         []interface{} `json:"assets"`
	Nodes          []interface{} `json:"nodes"`
	IsActive       bool          `json:"is_active"`
	IsPeriodic     bool          `json:"is_periodic"`
	Crontab        interface{}   `json:"crontab"`
	Interval       interface{}   `json:"interval"`
	SecretStrategy interface{}   `json:"secret_strategy"`
	SecretType     interface{}   `json:"secret_type"`
	Secret         string        `json:"secret"`
	PasswordRules  interface{}   `json:"password_rules"`
	Recipients     []interface{} `json:"recipients"`
	OrgID          string        `json:"org_id"`
	OrgName        string        `json:"org_name"`
	Comment        string        `json:"comment"`
	DateCreated    string        `json:"date_created"`
	DateUpdated    string        `json:"date_updated"`
	CreatedBy      string        `json:"created_by"`
}

// 执行改密计划请求
type ExecuteChangeSecretReq struct {
	Automation string `json:"automation"`
}

// 执行改密计划响应
type ExecuteChangeSecretResp struct {
	Task string `json:"task"`
}
