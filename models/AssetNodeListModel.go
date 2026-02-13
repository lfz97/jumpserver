package models

type AssetNode struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Name      string `json:"name"`
	FullValue string `json:"full_value"`
	OrgName   string `json:"org_name"`
}

type AssetNodeList []AssetNode
