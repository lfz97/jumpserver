package utils

import (
	"net/url"
)

// 拼接参数到url
func ParseUrl(baseUrl string, querys map[string]string) (string, error) {

	//基于baseUrl获取url对象
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	//拼接查询参数
	queryParams := url.Values{}
	for k, v := range querys {
		queryParams.Add(k, v)
	}

	//将查询参数更新到url对象中
	u.RawQuery = queryParams.Encode()

	//返回url的字符串形式
	return u.String(), nil
}
