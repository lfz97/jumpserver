package functions

import (
	"github.com/go-resty/resty/v2"
	"gopkg.in/twindagger/httpsig.v1"
	"log"
	"net/http"
)

type JMSClient struct {
	Url         string
	JMSClient_p *resty.Client
	PAMClient_p *resty.Client
	Logger_p    *log.Logger
}

// 为请求做签名，需设置参与签名的参数（APIID、APIKEY、参与的Header元素）
func Sign(req_p *http.Request, ID string, Secret string, headers []string) error {

	//设定sign参数获取sign对象
	Signer_p, err := httpsig.NewRequestSigner(ID, Secret, "hmac-sha256")
	if err != nil {
		return err
	}

	//对请求签名，插入Authorization头
	err = Signer_p.SignRequest(req_p, headers, nil)
	if err != nil {
		return err
	}

	return nil
}
