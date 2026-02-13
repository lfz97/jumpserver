package models

type Node struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Type struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type AssetInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Nodes   []Node `json:"nodes"`
	Type    Type   `json:"type"`
}

type AssetsListResult []AssetInfo
