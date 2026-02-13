package models

type Action struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Permission struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Accounts  []string `json:"accounts"`
	Protocols []string `json:"protocols"`
	Actions   []Action `json:"actions"`
	Users     []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"users"`
	User_groups []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user_groups"`
	Assets []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"assets"`
	Nodes        []Node `json:"nodes"`
	Created_by   string `json:"created_by"`
	Comment      string `json:"comment"`
	Is_active    bool   `json:"is_active"`
	Is_expired   bool   `json:"is_expired"`
	Is_valid     bool   `json:"is_valid"`
	Date_created string `json:"date_created"`
	Date_start   string `json:"date_start"`
	Date_expired string `json:"date_expired"`
}

type PermissionList []Permission
