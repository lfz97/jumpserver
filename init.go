package jumpserver

import (
	"github.com/go-resty/resty/v2"
	"github.com/lfz97/jumpserver/functions"
	"github.com/lfz97/jumpserver/mylogger"
	"net/http"
)

// 通过JMS&PAM APIKEY 获取认证client
func Init(Url string, JMSApiID string, JMSApiSecret string, PAMApiID string, PAMApiSecret string, logfilePath string) (*functions.JMSClient, error) {
	JMSClient_p := resty.New().SetDebug(true)
	PAMClient_p := resty.New().SetDebug(true)

	//通过SetPreRequestHook在请求前注册一个函数为请求签名
	JMSClient_p.SetPreRequestHook(func(client_p *resty.Client, request_p *http.Request) error {

		//jumpserver要求最少date头必须参与签名
		functions.Sign(request_p, JMSApiID, JMSApiSecret, []string{"(request-target)", "date"})

		return nil
	})
	PAMClient_p.SetPreRequestHook(func(client_p *resty.Client, request_p *http.Request) error {

		//jumpserver要求最少date头必须参与签名
		functions.Sign(request_p, PAMApiID, PAMApiSecret, []string{"(request-target)", "date"})

		return nil
	})
	logger_p, err := mylogger.LoggerInit(logfilePath)
	if err != nil {
		panic("初始化日志失败：" + err.Error())
	}

	//将设置好的resty client对象指针们放进结构体
	NewClient_p := &functions.JMSClient{
		Url:         Url,
		JMSClient_p: JMSClient_p,
		PAMClient_p: PAMClient_p,
		Logger_p:    logger_p,
	}
	return NewClient_p, nil
}
