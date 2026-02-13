package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lfz97/jumpserver/models"
	"github.com/lfz97/jumpserver/utils"
	"strconv"
)

// 根据用户名获取指定用户信息
func (J_p *JMSClient) GetUserByName(Name string) (*(models.UserList), error) {

	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/users/users/", map[string]string{
		"username": Name,
	})
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询用户失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	result := models.UserList{}
	err = json.Unmarshal(res_p.Body(), &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("未查询到用户: " + Name)
	}

	return &result, nil
}

// 根据节点full_value获取节点及子节点信息
func (J_p *JMSClient) GetAssetNodeByFullValue(fullValue string) (*(models.AssetNodeList), error) {

	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/assets/nodes/", map[string]string{
		"search": fullValue,
	})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}

	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询节点失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.AssetNodeList{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}

	if len(Result) == 0 {
		return nil, errors.New("未查询到组: " + fullValue)
	}

	return &Result, nil
}

// 根据节点id获取节点中的资产
func (J_p *JMSClient) GetAssetsByNodeID(Nid string, currentAsset bool) (*(models.AssetsListResult), error) {

	//设置是否仅显示当前节点资产
	show_current_asset := 0
	if currentAsset == true {
		show_current_asset = 1

	} else {
	}

	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/assets/assets/", map[string]string{
		"node_id":            Nid,
		"show_current_asset": strconv.Itoa(show_current_asset),
	})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询节点资产失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.AssetsListResult{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}
	if len(Result) == 0 {
		return nil, errors.New("本节点无资产：" + Nid)
	}

	return &Result, nil
}

// 根据用户id获取用户所有已授权的资产
func (J_p *JMSClient) GetUserAssetsListByUid(Uid string) (*(models.AssetsListResult), error) {

	url := fmt.Sprintf("%s/api/v1/perms/users/%s/assets/", (*J_p).Url, Uid)

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询用户授权资产失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	result := models.AssetsListResult{}
	err = json.Unmarshal(res_p.Body(), &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("未查询到用户授权资产，Uid: " + Uid)
	}

	return &result, nil
}

// 获取所有资产
func (J_p *JMSClient) GetAllAssets() (*(models.AssetsListResult), error) {

	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/assets/assets/", map[string]string{})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询所有资产失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	result := models.AssetsListResult{}
	err = json.Unmarshal(res_p.Body(), &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("无资产: ")
	}
	return &result, nil
}

// 根据IP地址查询资产
func (J_p *JMSClient) GetAssetByIP(IP string) (*(models.AssetsListResult), error) {
	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/assets/assets/", map[string]string{
		"address": IP,
	})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询资产失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.AssetsListResult{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}
	if len(Result) == 0 {
		return nil, errors.New("未查询到资产: " + IP)
	}
	return &Result, nil
}

// 根据名称查询资产
func (J_p *JMSClient) GetAssetByName(Name string) (*(models.AssetsListResult), error) {
	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/assets/assets/", map[string]string{
		"name": Name,
	})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询资产失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.AssetsListResult{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}
	if len(Result) == 0 {
		return nil, errors.New("未查询到资产: " + Name)
	}
	return &Result, nil
}

// 根据授权关系名称查询授权关系
func (J_p *JMSClient) GetAssetPermissionByName(permissionName string) (*(models.PermissionList), error) {
	url, err := utils.ParseUrl((*J_p).Url+"/api/v1/perms/asset-permissions/", map[string]string{
		"name": permissionName,
	})
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询授权关系失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.PermissionList{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}
	if len(Result) == 0 {
		return nil, errors.New("未找到此授权关系：" + permissionName)
	}
	return &Result, nil
}

// 根据授权id获取授权详细信息
func (J_p *JMSClient) GetAssetPermissionDetailByID(permissionID string) (*(models.Permission), error) {

	url := fmt.Sprintf("%s/api/v1/perms/asset-permissions/%s/", (*J_p).Url, permissionID)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).Get(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("查询授权关系详情失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.Permission{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}

	return &Result, nil
}

// 创建空授权关系
// actions可选参数:["connect","upload","download","copy","paste"]
// protocols可选参数去UI上找，建议配["all"]
func (J_p *JMSClient) CreateEmptyPermission(name string, actions []string, protocols []string) (*(models.Permission), error) {
	url := fmt.Sprintf("%s/api/v1/perms/asset-permissions/", (*J_p).Url)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	bodyTemplete := struct {
		Name      string   `json:"name"`
		Actions   []string `json:"actions"`
		Protocols []string `json:"protocols"`
		Is_active bool     `json:"is_active"`
	}{
		Name:      name,
		Actions:   actions,
		Protocols: protocols,
		Is_active: true,
	}
	jsonBodybytes, err := json.Marshal(bodyTemplete)
	if err != nil {
		return nil, err
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).SetBody(string(jsonBodybytes)).Post(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 201 {
		return nil, errors.New("创建授权关系失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.Permission{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}

	return &Result, nil
}

// 定义授权关系配置结构体
type PermissionConfig struct {
	Name        string   `json:"name"`                  //授权关系名称
	Users       []string `json:"users,omitempty"`       //用户ID列表
	User_groups []string `json:"user_groups,omitempty"` //用户组ID列表
	Assets      []string `json:"assets,omitempty"`      //资产ID列表
	Nodes       []string `json:"nodes,omitempty"`       //节点ID列表
	Accounts    []string `json:"accounts,omitempty"`    /*"@ALL" ：所有账号
	"@SPEC"：指定账号
	"@INPUT"：手动账号
	"@USER"：同名账号
	*/
	Actions      []string `json:"actions"`                /*允许的操作，默认：all，可选值：[ "connect", "upload", "copy", "paste","download" ]*/
	Protocols    []string `json:"protocols"`              //支持协议。可选值去UI上找，建议配["all"]
	Is_active    bool     `json:"is_active,omitempty"`    //是否启用
	Date_start   string   `json:"date_start,omitempty"`   //生效时间，格式："2123-01-30T10:53:23.879Z"
	Date_expired string   `json:"date_expired,omitempty"` //失效时间，格式："2123-01-30T10:53:23.879Z"
	Comment      string   `json:"comment,omitempty"`      //备注
}

// 根据授权关系ID及配置更新授权关系
func (J_p *JMSClient) UpdatePermission(config PermissionConfig, permissionID string) (*(models.Permission), error) {
	url := fmt.Sprintf("%s/api/v1/perms/asset-permissions/%s/", (*J_p).Url, permissionID)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	bodyBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	res_p, err := (*J_p).JMSClient_p.R().SetHeaders(headers).SetBody(string(bodyBytes)).Put(url)
	if err != nil {
		return nil, err
	}
	if res_p.StatusCode() != 200 {
		return nil, errors.New("更新授权关系失败，状态码：" + strconv.Itoa(res_p.StatusCode()))
	}
	Result := models.Permission{}
	err = json.Unmarshal(res_p.Body(), &Result)
	if err != nil {
		return nil, err
	}
	return &Result, nil
}
