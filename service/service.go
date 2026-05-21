package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lfz97/jumpserver/functions"
	"github.com/lfz97/jumpserver/logic"
	"github.com/lfz97/jumpserver/models"
	"github.com/lfz97/jumpserver/models/serviceModel"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 申请新授权，校验用户和节点，自动创建授权关系并授权到用户
func RequestNewPermission(JMSClient_p *functions.JMSClient, requestUsers []string, assetNodeNames []string) error {
	(*JMSClient_p).Logger_p.Println("收到新的授权请求，用户列表：" + strings.Join(requestUsers, ",") + "，节点列表：" + strings.Join(assetNodeNames, ","))
	//检查用户是否存在，获取存在的用户表和不存在的用户列表
	ExistingUsers, NotExistingUsers, err := logic.CheckUserExists(JMSClient_p, requestUsers)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("检查用户存在性失败，错误信息：" + err.Error())
		return err
	}
	JMSClient_p.Logger_p.Println("以下用户不存在，无法添加到授权关系中:")
	for _, userName := range NotExistingUsers {
		(*JMSClient_p).Logger_p.Println("用户名称：" + userName)
	}
	if len(ExistingUsers) == 0 {
		(*JMSClient_p).Logger_p.Println("没有存在的用户，无法创建授权关系")
		return errors.New("没有存在的用户，无法创建授权关系")
	}

	//检查节点是否存在，获取存在的节点全名列表和不存在的节点全名列表
	ExistingNodes, NotExistingFullNameNodes, err := logic.CheckNodeExists(JMSClient_p, assetNodeNames)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("检查节点存在性失败，错误信息：" + err.Error())
		return err
	}

	//打印不存在的节点列表
	JMSClient_p.Logger_p.Println("以下节点不存在，无法授权:")
	for _, nodeName := range NotExistingFullNameNodes {
		(*JMSClient_p).Logger_p.Println("节点名称：" + nodeName)
	}

	//根据存在的节点列表，检查与node同名授权关系是否存在，获取存在的授权关系表和不存在的节点全名列表
	ExistingNamesakePermissions, NotExistingNamesakePermissionNodes, err := logic.CheckNamesakePermissions(JMSClient_p, ExistingNodes)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("检查同名授权关系存在性失败，错误信息：" + err.Error())
		return err
	}

	//为不存在的同名授权关系节点创建授权关系模板
	FailedCreatedPermissions := map[string]string{}
	for NodeName, NodeID := range NotExistingNamesakePermissionNodes {
		PermissionID, err := logic.CreatePermissionTemplate(JMSClient_p, NodeName, NodeID)
		if err != nil {
			//创建失败的，添加到失败列表中
			FailedCreatedPermissions[NodeName] = NodeID
			(*JMSClient_p).Logger_p.Println("创建授权关系模板失败，节点名称：" + NodeName + "，节点ID：" + NodeID + "，错误信息：" + err.Error())
		} else {
			//创建成功的，添加到已存在授权关系表中
			ExistingNamesakePermissions[NodeName] = PermissionID
		}
	}

	//将用户和授权M x N 全交叉添加到授权关系中
	for PermissionName, PermissionID := range ExistingNamesakePermissions {
		for username, userid := range ExistingUsers {
			err := logic.InsertUserToPermission(JMSClient_p, PermissionID, userid)
			if err != nil {
				(*JMSClient_p).Logger_p.Println("添加用户到授权关系失败，授权关系名称：" + PermissionName + "，用户名称：" + username + "，错误信息：" + err.Error())
				continue
			}
		}
	}

	return nil
}

// 申请ROOT权限，校验用户现有授权清单，自动创建授权关系并授权到用户
func RequestRootPermission(JMSClient_p *functions.JMSClient, requestUser string, assetNodeNames []string, assetIPlist []string, Days int) error {

	(*JMSClient_p).Logger_p.Println("收到新的ROOT授权请求，用户：" + requestUser + "，节点列表：" + strings.Join(assetNodeNames, ",") + "，IP列表：" + strings.Join(assetIPlist, ","))

	//检查用户
	userList_p, err := JMSClient_p.GetUserByName(requestUser)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("查询用户失败，用户名称：" + requestUser + "，错误信息：" + err.Error())
		return err
	}

	//如果查出多个用户，无法确定唯一用户，返回错误
	if len(*userList_p) > 1 {
		(*JMSClient_p).Logger_p.Println("查询到多个用户，无法确定唯一用户，用户名称：" + requestUser + ",查询结果：\n")
		for _, user := range *userList_p {
			(*JMSClient_p).Logger_p.Println("用户ID：" + user.ID + "，用户名：" + user.Username + "，邮箱：" + user.Email)
		}
		return errors.New("查询到多个用户，无法确定唯一用户，用户名称：" + requestUser)
	}

	//检查授权天数必须大于0
	if Days <= 0 {
		(*JMSClient_p).Logger_p.Println("授权天数必须大于0，输入的天数：" + strconv.Itoa(Days))
		return errors.New("授权天数必须大于0")
	}

	//获取结束时间
	utcEnd := time.Now().UTC().Add(time.Duration(Days) * 24 * time.Hour).Format(time.RFC3339)

	//检查申请的节点是否存在，获取存在的节点全名列表和不存在的节点全名列表
	ExistingNodes, NotExistingFullNameNodes, err := logic.CheckNodeExists(JMSClient_p, assetNodeNames)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("检查节点存在性失败，错误信息：" + err.Error())
		return err
	}

	//记录不存在的节点列表
	JMSClient_p.Logger_p.Println("以下节点不存在，无法授权:")
	for _, nodeName := range NotExistingFullNameNodes {
		(*JMSClient_p).Logger_p.Println("节点名称：" + nodeName)
	}

	//展开存在的节点列表，获取节点下的资产列表，同时根据输入的IP列表获取资产列表，合并去重后得到最终的目标资产列表
	ExpandedAssets := models.AssetsListResult{}
	for _, id := range ExistingNodes {
		//根据节点id获取节点中的资产列表，不包含子节点资产
		AssetList_p, err := JMSClient_p.GetAssetsByNodeID(id, false)
		if err != nil {
			(*JMSClient_p).Logger_p.Println("查询节点资产失败，节点ID：" + id + "，错误信息：" + err.Error())
			continue
		}
		ExpandedAssets = append(ExpandedAssets, (*AssetList_p)...)
	}
	for _, ip := range assetIPlist {
		//根据IP获取资产列表
		AssetList_p, err := JMSClient_p.GetAssetByIP(ip)
		if err != nil {
			(*JMSClient_p).Logger_p.Println("查询IP资产失败，IP地址：" + ip + "，错误信息：" + err.Error())
			continue
		}
		ExpandedAssets = append(ExpandedAssets, (*AssetList_p)...)
	}

	//如果最终目标资产列表为空，说明没有任何资产可以授权，返回错误
	if len(ExpandedAssets) == 0 {
		(*JMSClient_p).Logger_p.Println("未查询到任何目标资产，无法授权")
		return errors.New("未查询到任何目标资产，无法授权")

	}

	//获取用户已授权的资产列表
	AssetList_p, err := JMSClient_p.GetUserAssetsListByUid((*userList_p)[0].ID)
	if err != nil {
		return errors.New("查询用户资产失败，用户ID：" + (*userList_p)[0].ID + "，错误信息：" + err.Error())
	}

	//对比用户已授权资产和目标资产，获取需要授权的资产列表
	AssetCanBeAuthorized := models.AssetsListResult{}    //可以授权ROOT的资产列表（存在于用户的授权关系中）
	AssetCanNotBeAuthorized := models.AssetsListResult{} //不可以授权ROOT的资产列表（不存在于用户的授权关系中）
	for _, asset := range ExpandedAssets {
		for i := 0; i < len(*AssetList_p); i++ {
			if asset.ID == (*AssetList_p)[i].ID {
				AssetCanBeAuthorized = append(AssetCanBeAuthorized, asset)
				break
			}
			//如果循环到最后一个已授权资产仍未匹配到，则说明目标资产与用户授权资产无交叉，无法授权ROOT，记录到不可授权列表中
			if i == len(*AssetList_p)-1 {
				AssetCanNotBeAuthorized = append(AssetCanNotBeAuthorized, asset)
			}
		}

	}

	//打印不可授权资产列表
	(*JMSClient_p).Logger_p.Println("以下资产不在用户的授权范围，禁止申请root授权：")
	for _, asset := range AssetCanNotBeAuthorized {
		(*JMSClient_p).Logger_p.Println("资产ID："+asset.ID+"，资产名称："+asset.Name, ", 资产IP："+asset.Address)
	}

	//如果目标资产与用户授权资产无任何交叉，那么没有任何资产可以授权，返回错误
	if len(AssetCanBeAuthorized) == 0 {
		(*JMSClient_p).Logger_p.Println("没有任何资产在用户的授权范围内，无法申请root授权")
		return errors.New("没有任何资产在用户的授权范围内，无法申请root授权")
	}

	//将可以授权ROOT的资产列表解析为资产ID列表
	AssetCanBeAuthorizedIds := []string{}
	for _, asset := range AssetCanBeAuthorized {
		AssetCanBeAuthorizedIds = append(AssetCanBeAuthorizedIds, asset.ID)
	}

	//创建授权关系模板
	RandomId := uuid.New().String()
	PermissionName := "AUTO_ROOT_" + requestUser + "_" + RandomId
	PermissionID, err := logic.CreatePermissionTemplate(JMSClient_p, PermissionName, "")
	if err != nil {
		(*JMSClient_p).Logger_p.Println("创建授权关系模板失败，授权关系名称：" + "AUTO_ROOT_" + requestUser + "_" + RandomId + "，错误信息：" + err.Error())
		return err
	}

	//更新授权关系
	_, err = JMSClient_p.UpdatePermission(functions.PermissionConfig{
		Name:   PermissionName,
		Users:  []string{(*userList_p)[0].ID},
		Assets: AssetCanBeAuthorizedIds,
		Accounts: []string{
			"@SPEC",
			"root",
		},
		Actions:      []string{"connect", "upload", "copy", "paste"},
		Protocols:    []string{"all"},
		Date_expired: utcEnd,
	}, PermissionID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("更新授权关系失败，授权关系名称：" + PermissionName + "，错误信息：" + err.Error())
		return err
	}

	return nil
}

// 检出指定资产的指定账号的密码
func CheckoutPassword(JMSClient_p *functions.JMSClient, requestUser string, assetNames []string, accounts []string) ([]serviceModel.SecretInfo, error) {
	//检查用户
	userList_p, err := JMSClient_p.GetUserByName(requestUser)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("查询用户失败，用户名称：" + requestUser + "，错误信息：" + err.Error())
		return nil, err
	}

	//如果查出多个用户，无法确定唯一用户，返回错误
	if len(*userList_p) > 1 {
		(*JMSClient_p).Logger_p.Println("查询到多个用户，无法确定唯一用户，用户名称：" + requestUser + ",查询结果：\n")
		for _, user := range *userList_p {
			(*JMSClient_p).Logger_p.Println("用户ID：" + user.ID + "，用户名：" + user.Username + "，邮箱：" + user.Email)
		}
		return nil, errors.New("查询到多个用户，无法确定唯一用户，用户名称：" + requestUser)
	}

	//检查资产是否存在，获取存在的资产对象列表和不存在的资产名称列表
	ExistingAssets := models.AssetsListResult{}
	NotExistingAssetNames := []string{}
	for _, assetNames := range assetNames {
		AssetsListResult, err := JMSClient_p.GetAssetByName(assetNames)
		if err != nil {
			(*JMSClient_p).Logger_p.Println("查询资产失败，资产名称：" + assetNames + "，错误信息：" + err.Error())
			NotExistingAssetNames = append(NotExistingAssetNames, assetNames)
			continue
		}
		ExistingAssets = append(ExistingAssets, (*AssetsListResult)...)
	}

	//获取用户已授权的资产列表
	AssetList_p, err := JMSClient_p.GetUserAssetsListByUid((*userList_p)[0].ID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("查询用户资产失败，用户ID：" + (*userList_p)[0].ID + "，错误信息：" + err.Error())
		return nil, errors.New("查询用户资产失败，用户ID：" + (*userList_p)[0].ID + "，错误信息：" + err.Error())
	}

	//对比用户已授权资产和目标资产，获取需要授权的资产列表
	PwdCanBeCheckout := models.AssetsListResult{}    //可以取密码的资产列表（存在于用户的授权关系中）
	PwdCanNotBeCheckout := models.AssetsListResult{} //不可取密码的资产列表（不存在于用户的授权关系中）
	for _, asset := range ExistingAssets {
		for i := 0; i < len(*AssetList_p); i++ {
			if asset.ID == (*AssetList_p)[i].ID {
				PwdCanBeCheckout = append(PwdCanBeCheckout, asset)
				break
			}
			//如果循环到最后一个已授权资产仍未匹配到，则说明目标资产与用户授权资产无交叉，无法授权ROOT，记录到不可授权列表中
			if i == len(*AssetList_p)-1 {
				PwdCanNotBeCheckout = append(PwdCanNotBeCheckout, asset)
			}
		}

	}
	(*JMSClient_p).Logger_p.Println("以下资产不在用户的授权范围，禁止申请密码：")
	for _, asset := range PwdCanNotBeCheckout {
		(*JMSClient_p).Logger_p.Println("资产ID："+asset.ID+"，资产名称："+asset.Name, ", 资产IP："+asset.Address)
	}
	//获取密码
	SecretInfoList := []serviceModel.SecretInfo{}
	for _, asset := range PwdCanBeCheckout {

		SecretInfo := serviceModel.SecretInfo{
			AssetName:    asset.Name,
			AssetID:      asset.ID,
			AssetAddress: asset.Address,
		}
		for _, account := range accounts {
			secret := serviceModel.Secret{
				Account: account,
			}
			Secret_p, err := JMSClient_p.GetSecret(asset.ID, account)
			if err != nil {
				(*JMSClient_p).Logger_p.Println("查询密码失败，资产ID：" + asset.ID + "，账号：" + account + "，错误信息：" + err.Error())
				continue
			}
			secret.Password = (*Secret_p).Secret
			SecretInfo.Secrets = append(SecretInfo.Secrets, secret)
		}
		SecretInfoList = append(SecretInfoList, SecretInfo)
	}
	return SecretInfoList, nil
}

// ChangeSecretItem 改密任务项
type ChangeSecretItem struct {
	AssetID string
	Account string
}

// 改密：并发执行，轮询确认密码更新后自动删除改密计划
func ChangeSecret(JMSClient_p *functions.JMSClient, items []ChangeSecretItem) {
	(*JMSClient_p).Logger_p.Println("收到改密请求，共 " + strconv.Itoa(len(items)) + " 个账号")

	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(item ChangeSecretItem) {
			defer wg.Done()
			changeSingleSecret(JMSClient_p, item)
		}(item)
	}

	// 监听 goroutine：全部完成后退出
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	<-done
	(*JMSClient_p).Logger_p.Println("改密请求全部处理完毕")
}

// 单个账号改密流程
func changeSingleSecret(JMSClient_p *functions.JMSClient, item ChangeSecretItem) {
	prefix := item.Account + "@" + item.AssetID

	// 0. 记录当前密码时间
	oldDateUpdated := ""
	accountList, err := JMSClient_p.GetSpecifiedAccount(item.AssetID, item.Account)
	if err != nil {
		(*JMSClient_p).Logger_p.Println(prefix + " 查询账号失败：" + err.Error())
		return
	}
	if len(accountList.Results) > 0 {
		oldDateUpdated = accountList.Results[0].DateUpdated
		(*JMSClient_p).Logger_p.Println(prefix + " 当前密码更新时间：" + oldDateUpdated)
	} else {
		(*JMSClient_p).Logger_p.Println(prefix + " 未查询到账号，跳过")
		return
	}

	// 1. 创建改密计划
	randomId := uuid.New().String()
	req := models.CreateChangeSecretAutomationReq{
		Name:           "AUTO_CHANGE_SECRET_" + randomId,
		Accounts:       []string{item.Account},
		Assets:         []string{item.AssetID},
		IsActive:       true,
		IsPeriodic:     false,
		SecretStrategy: "random",
		SecretType:     "password",
		PasswordRules:  map[string]int{"length": 30},
	}

	automation, err := JMSClient_p.CreateChangeSecretAutomation(req)
	if err != nil {
		(*JMSClient_p).Logger_p.Println(prefix + " 创建改密计划失败：" + err.Error())
		return
	}
	(*JMSClient_p).Logger_p.Println(prefix + " 改密计划创建成功，ID：" + automation.ID)

	// 2. 执行改密计划
	execResult, err := JMSClient_p.ExecuteChangeSecret(automation.ID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println(prefix + " 执行改密计划失败：" + err.Error())
		// 执行失败也要清理
		JMSClient_p.DeleteChangeSecretAutomation(automation.ID)
		return
	}
	(*JMSClient_p).Logger_p.Println(prefix + " 改密计划执行成功，任务ID：" + execResult.Task)

	// 3. 轮询密码是否更新（最多10次 ，每次1秒）
	success := false
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		accountList, err := JMSClient_p.GetSpecifiedAccount(item.AssetID, item.Account)
		if err != nil {
			(*JMSClient_p).Logger_p.Println(prefix + " 第" + strconv.Itoa(i+1) + "次查询密码状态失败：" + err.Error())
			continue
		}
		if len(accountList.Results) > 0 && accountList.Results[0].DateUpdated != oldDateUpdated {
			(*JMSClient_p).Logger_p.Println(prefix + " 密码已更新，新时间：" + accountList.Results[0].DateUpdated)
			success = true
			break
		}
	}
	if !success {
		(*JMSClient_p).Logger_p.Println(prefix + " 轮询超时，密码未确认更新")
	}

	// 4. 删除改密计划
	err = JMSClient_p.DeleteChangeSecretAutomation(automation.ID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println(prefix + " 删除改密计划失败：" + err.Error())
	} else {
		(*JMSClient_p).Logger_p.Println(prefix + " 改密计划已删除，ID：" + automation.ID)
	}
}
