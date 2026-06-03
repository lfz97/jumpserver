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

type AssetsListPaginated struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []AssetInfo `json:"results"`
}
