package serviceModel

type Secret struct {
	Account  string `json:"Account"`
	Password string `json:"Password"`
}
type SecretInfo struct {
	AssetName    string   `json:"AssetName"`
	AssetID      string   `json:"AssetID"`
	AssetAddress string   `json:"AssetAddress"`
	Secrets      []Secret `json:"Secrets"`
}
