package functions

import (
	"encoding/json"
	"fmt"
	"github.com/lfz97/jumpserver/models"
	"github.com/lfz97/jumpserver/utils"
	"strconv"
)

// 根据资产id和账号名获取密码
func (J_p *JMSClient) GetSecret(assetID string, account string) (*(models.Secret), error) {
	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/accounts/integration-applications/account-secret/", map[string]string{
		"asset":   assetID,
		"account": account,
	})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Source":     "jms-pam",
	}
	res_p, err := (*J_p).PAMClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, fmt.Errorf("GetSecret failed, status code: %s", strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.Secret{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}

	return &Result, nil
}
