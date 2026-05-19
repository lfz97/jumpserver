package models

type AccountList struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Account   `json:"results"`
}

type Account struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	SecretType struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"secret_type"`
	SpecInfo struct {
	} `json:"spec_info"`
	CreatedBy string      `json:"created_by"`
	Comment   string      `json:"comment"`
	SuFrom    interface{} `json:"su_from"`
	Asset     struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Address string `json:"address"`
		Type    struct {
			Value string `json:"value"`
			Label string `json:"label"`
		} `json:"type"`
		Category struct {
			Value string `json:"value"`
			Label string `json:"label"`
		} `json:"category"`
		Platform struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"platform"`
		AutoConfig struct {
			SuEnabled      bool `json:"su_enabled"`
			DomainEnabled  bool `json:"domain_enabled"`
			AnsibleEnabled bool `json:"ansible_enabled"`
			ID             int  `json:"id"`
			AnsibleConfig  struct {
				AnsibleConnection string `json:"ansible_connection"`
			} `json:"ansible_config"`
			PingEnabled bool   `json:"ping_enabled"`
			PingMethod  string `json:"ping_method"`
			PingParams  struct {
			} `json:"ping_params"`
			GatherFactsEnabled bool   `json:"gather_facts_enabled"`
			GatherFactsMethod  string `json:"gather_facts_method"`
			GatherFactsParams  struct {
			} `json:"gather_facts_params"`
			ChangeSecretEnabled bool   `json:"change_secret_enabled"`
			ChangeSecretMethod  string `json:"change_secret_method"`
			ChangeSecretParams  struct {
			} `json:"change_secret_params"`
			PushAccountEnabled bool   `json:"push_account_enabled"`
			PushAccountMethod  string `json:"push_account_method"`
			PushAccountParams  struct {
				Sudo   string `json:"sudo"`
				Shell  string `json:"shell"`
				Groups string `json:"groups"`
			} `json:"push_account_params"`
			VerifyAccountEnabled bool   `json:"verify_account_enabled"`
			VerifyAccountMethod  string `json:"verify_account_method"`
			VerifyAccountParams  struct {
			} `json:"verify_account_params"`
			GatherAccountsEnabled bool   `json:"gather_accounts_enabled"`
			GatherAccountsMethod  string `json:"gather_accounts_method"`
			GatherAccountsParams  struct {
			} `json:"gather_accounts_params"`
			Platform int `json:"platform"`
		} `json:"auto_config"`
	} `json:"asset"`
	Version int `json:"version"`
	Source  struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"source"`
	SourceID     interface{} `json:"source_id"`
	Connectivity struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"connectivity"`
	OrgID       string `json:"org_id"`
	OrgName     string `json:"org_name"`
	HasSecret   bool   `json:"has_secret"`
	Privileged  bool   `json:"privileged"`
	IsActive    bool   `json:"is_active"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}
