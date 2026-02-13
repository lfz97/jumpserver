package logic

import (
	"github.com/lfz97/jumpserver/functions"
	"strings"
)

// 检查节点是否存在，返回存在的节点全名列表和不存在的节点全名列表
func CheckNodeExists(JMSClient_p *(functions.JMSClient), assetNodeNames []string) (map[string]string, []string, error) {

	//存在的节点全名表,key：节点全名，value：节点ID
	ExistingNodes := map[string]string{}

	//不存在的节点全名列表
	NotExistingFullNameNodes := []string{}
	for _, name := range assetNodeNames {

		//查询节点，如果没有error，说明节点存在
		AssetNodeList_p, err := JMSClient_p.GetAssetNodeByFullValue(name)
		if err != nil {
			//添加到不存在节点
			NotExistingFullNameNodes = append(NotExistingFullNameNodes, name)
			(*JMSClient_p).Logger_p.Println("节点不存在，跳过：" + name)

			continue
		}
		//添加节点全名与对应的节点ID到表中
		ExistingNodes[name] = (*AssetNodeList_p)[0].ID
		(*JMSClient_p).Logger_p.Println("节点存在，添加到待检查列表，名称：" + name + "，ID：" + (*AssetNodeList_p)[0].ID)
	}
	return ExistingNodes, NotExistingFullNameNodes, nil
}

// 检查与node同名授权关系是否存在，返回存在的授权关系表和不存在的节点全名列表
func CheckNamesakePermissions(JMSClient_p *(functions.JMSClient), Nodes map[string]string) (map[string]string, map[string]string, error) {
	//存在的授权关系表，key：授权关系名称，value：授权关系ID
	ExistingNamesakePermissions := map[string]string{}
	//不存在的同名授权关系节点列表，key：节点名称，value：节点ID
	NotExistingNamesakePermissionNodes := map[string]string{}

	for nodeFullName, nodeID := range Nodes {

		//查询是否存在与节点名称相同的授权关系
		Permissions_p, err := JMSClient_p.GetAssetPermissionByName(nodeFullName)
		if err != nil {

			//如果错误中包含“未找到此授权关系”，则说明授权关系不存在，记录到不存在列表中
			if strings.Contains(err.Error(), "未找到此授权关系") {
				(*JMSClient_p).Logger_p.Println("节点对应授权关系不存在，节点名称：" + nodeFullName + "，节点ID：" + nodeID)
				NotExistingNamesakePermissionNodes[nodeFullName] = nodeID
				continue

			}

			//其他错误，打印日志并跳过
			(*JMSClient_p).Logger_p.Println("查询节点对应授权关系失败，跳过：" + nodeFullName + "，错误信息：" + err.Error())
			continue
		}

		//如果查询到的授权关系名称和节点名称完全相同，则保存授权关系ID
		if (*Permissions_p)[0].Name == nodeFullName {
			ExistingNamesakePermissions[nodeFullName] = (*Permissions_p)[0].ID

			//如果查询到的授权关系名称和节点名称不完全匹配（比如模糊查询部分匹配），则重新创建一个空的授权关系
		} else {
			(*JMSClient_p).Logger_p.Println("查询到的授权关系名称和节点名称不完全匹配，期望的名称：" + nodeFullName + "，实际名称：" + (*Permissions_p)[0].Name)
			NotExistingNamesakePermissionNodes[nodeFullName] = nodeID
		}
	}
	return ExistingNamesakePermissions, NotExistingNamesakePermissionNodes, nil
}

// 创建节点对应的授权关系模板
func CreatePermissionTemplate(JMSClient_p *(functions.JMSClient), PermissionName string, NodeID string) (string, error) {

	//创建空的授权关系
	Permission_p, err := JMSClient_p.CreateEmptyPermission(PermissionName, []string{"connect", "upload", "copy", "paste"}, []string{"all"})
	if err != nil {
		(*JMSClient_p).Logger_p.Println("创建节点对应授权关系失败：" + PermissionName + "，错误信息：" + err.Error())
		return "", err

	}
	ND := []string{NodeID}
	if NodeID == "" {
		ND = nil
	}
	//在新创建的空授权关系中添加与授权关系名称相同的资产节点，并配置默认账号
	_, err = JMSClient_p.UpdatePermission(functions.PermissionConfig{
		Name:      PermissionName,
		Nodes:     ND,
		Accounts:  []string{"@SPEC", "@INPUT", "@USER", "dominos-user"},
		Actions:   []string{"connect", "upload", "copy", "paste"},
		Protocols: []string{"all"},
	}, (*Permission_p).ID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("更新节点对应授权关系失败：" + PermissionName + "，错误信息：" + err.Error())
		return "", err
	}
	return (*Permission_p).ID, nil
}

// 检查用户是否存在，返回存在的用户表和不存在的用户列表
func CheckUserExists(JMSClient_p *(functions.JMSClient), users []string) (map[string]string, []string, error) {
	//存在的用户表,key：用户名，value：用户ID
	ExistingUsers := map[string]string{}
	NotExistingUsers := []string{}
	for _, username := range users {
		UserList_p, err := JMSClient_p.GetUserByName(username)
		if err != nil {
			(*JMSClient_p).Logger_p.Println("查询用户失败，跳过：" + username + "，错误信息：" + err.Error())
			NotExistingUsers = append(NotExistingUsers, username)
		} else {
			ExistingUsers[username] = (*UserList_p)[0].ID
		}
	}
	return ExistingUsers, NotExistingUsers, nil
}

// 将用户添加到指定的授权关系中
func InsertUserToPermission(JMSClient_p *(functions.JMSClient), PermissionID string, UserID string) error {
	PermissionDetail_p, err := JMSClient_p.GetAssetPermissionDetailByID(PermissionID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("获取授权关系详情失败，跳过：" + PermissionID + "，错误信息：" + err.Error())
		return err
	}
	//获取授权关系详情中的当前各项信息
	PermissionName := (*PermissionDetail_p).Name
	Accounts := (*PermissionDetail_p).Accounts
	Protocols := (*PermissionDetail_p).Protocols
	//获取当前已有的用户ID列表
	UserIDs := []string{}
	for _, userinfo := range (*PermissionDetail_p).Users {
		UserIDs = append(UserIDs, userinfo.ID)
	}
	//将申请的用户ID添加到列表中
	UserIDs = append(UserIDs, UserID)
	//获取当前已有的节点ID列表
	NodeIDs := []string{}
	for _, nodeInfo := range (*PermissionDetail_p).Nodes {
		NodeIDs = append(NodeIDs, nodeInfo.ID)
	}

	//获取当前已有的Action列表
	Actions := []string{}
	for _, actionInfo := range (*PermissionDetail_p).Actions {
		Actions = append(Actions, actionInfo.Value)
	}
	//更新授权关系将申请的用户都添加进去
	_, err = JMSClient_p.UpdatePermission(functions.PermissionConfig{
		Name:      PermissionName,
		Nodes:     NodeIDs,
		Accounts:  Accounts,
		Actions:   Actions,
		Users:     UserIDs,
		Protocols: Protocols,
	}, PermissionID)
	if err != nil {
		(*JMSClient_p).Logger_p.Println("更新授权关系失败，跳过：" + PermissionName + "，错误信息：" + err.Error())
		return err
	}
	(*JMSClient_p).Logger_p.Println("更新授权关系成功：" + PermissionName + "，ID：" + PermissionID)
	return nil

}
