package wx

import (
	"net/url"
	"encoding/json"
	"reflect"
	"strconv"
	"time"
	"encoding/xml"
)

const (
	oauth_code_url = "https://open.weixin.qq.com/connect/oauth2/authorize"
	oauth_token_url = "https://api.weixin.qq.com/sns/oauth2/access_token"
	pay_url = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type WxController struct{
	WxLibController
	/********************基本参数****************************/
	//微信公众号身份的唯一标识。审核通过后，在微信发送的邮件中查看
	//或者微信小程序的appid，可在小程序后台设置－开发设置内查看
	appid string
	//受理商id，身份标识
	mchid string
	//商户支付密钥key。审核通过后，在微信发送的邮件中查看
	key string
	//jsapi接口中获取openid，审核后在公众平台开启开发模式后可查看
	//或者微信小程序的appsecret，可在小程序后台设置－开发设置内查看
	appsecret string
	//获取access_token过程中的跳转uri，通过跳转将code传入jsapi支付页面

	//jsapi参数
	parameters map[string]string
	//code码，用以获取openid
	code string
	//使用统一支付接口得到的预支付id
	prepay_id string
}

func (this *WxController) RegisterWx() {
	this.appid = appid
	this.appsecret = appsecret
	this.mchid = mchid
	this.key = key
}

//设置请求参数
func (this *WxController) SetParameter(key string, value string){
	if(this.parameters == nil){
		this.parameters = make(map[string]string)
	}
	this.parameters[key] = value
}

func (this *WxController) SetCode(code string){
	this.code = code
}

//设置prepay_id
func (this *WxController) SetPrepayId (prepayId string) {
	this.prepay_id = prepayId
}

func (this WxController) CreateOauthUrlForCode(redirectUrl string) string {
	urlObj := make(map[string]string)
	urlObj["appid"] = this.appid
	urlObj["redirect_uri"] = url.QueryEscape(redirectUrl)
	urlObj["response_type"] = "code"
	urlObj["scope"] = "snsapi_base"
	urlObj["state"] = "STATE#wechat_redirect"
	codeStr := this.FormatParams(urlObj)
	return oauth_code_url + "?" + codeStr
}

//生成可以获得openid的url
func (this WxController) CreateOauthUrlForOpenid() string {
	urlObj := make(map[string]string)
	urlObj["appid"] = this.appid
	urlObj["secret"] = this.appsecret
	urlObj["code"] = this.code
	urlObj["grant_type"] = "authorization_code"
	openidStr := this.FormatParams(urlObj)
	return oauth_token_url + "?" + openidStr
}

func (this WxController) GetOpenid() (string, error) {
	url := this.CreateOauthUrlForOpenid()
	result, err := this.http_get(url)
	if err != nil {
		return "", this.Error("get openid error", err)
	}
	var res map[string]interface{}
	jsonErr := json.Unmarshal([]byte(result), &res)
	if jsonErr != nil {
		return "", this.Error("get openid json error", jsonErr)
	}

	errcode := reflect.ValueOf(res["errcode"])
	if errcode.IsValid() && errcode.Float() > 0 {
		errmsg := reflect.ValueOf(res["errmsg"]).String()
		return "", this.Error("get openid result error: " + strconv.FormatFloat(errcode.Float(), 'f', 0, 64) + "-" + errmsg, nil)
	}
	return reflect.ValueOf(res["openid"]).String(), nil
}

func (this WxController) CreateXml() (string, error) {
	//检测必填参数
	if this.parameters["out_trade_no"] == ""  {
		return "", this.Error("缺少统一支付接口必填参数out_trade_no(商户订单号)！", nil)
	}else if this.parameters["body"] == "" {
		return "", this.Error("缺少统一支付接口必填参数body(商品描述)！", nil)
	}else if  this.parameters["total_fee"] == "" {
		return "", this.Error("缺少统一支付接口必填参数total_fee(交易金额)！", nil)
	}else if  this.parameters["notify_url"] == ""  {
		return "", this.Error("缺少统一支付接口必填参数notify_url(异步接收微信支付结果通知的回调地址)！", nil)
	}else if  this.parameters["trade_type"] == ""  {
		return "", this.Error("缺少统一支付接口必填参数trade_type(交易类型)！", nil)
	}else if  this.parameters["spbill_create_ip"] == "" {
		return "", this.Error("缺少统一支付接口必填参数spbill_create_ip(终端ip)", nil)
	}else if  this.parameters["trade_type"] == "JSAPI" && this.parameters["openid"] == "" {
		return "", this.Error("统一支付接口中，缺少必填参数openid！trade_type为JSAPI时，openid为必填参数！", nil)
	}else if  this.parameters["trade_type"] == "NATIVE" && this.parameters["product_id"] == "" {
		return "", this.Error("统一支付接口中，缺少必填参数product_id！trade_type为JSAPI时，product_id为必填参数！", nil)
	}

	this.parameters["appid"] = this.appid//公众账号ID
	this.parameters["mch_id"] = this.mchid//商户号
	this.parameters["nonce_str"] = this.CreateNoncestr(32) //随机字符串
	this.parameters["sign"] = this.GetSign(this.parameters, this.key)//签名
	return  this.ParamsToXml(this.parameters), nil
}


type XmlResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Result_code string `xml:"result_code"`
	Err_code string `xml:"err_code"`
	Err_code_des string `xml:"err_code_des"`
	Prepay_id   string `xml:"prepay_id"`
}

func (this WxController) GetPrepayId() (string, error) {
	strXml, err := this.CreateXml()
	if err != nil {
		return "", this.Error("get prepay_id error", err)
	}
	result, err := this.http_post(pay_url, strXml)
	if err != nil {
		return "", this.Error("get prepay_id error", err)
	}

	res := new(XmlResp)
	xmlErr :=  xml.Unmarshal(result, &res)
	if xmlErr != nil {
		return "", this.Error("get prepay_id xml error", xmlErr)
	}
	if res.Return_code != "SUCCESS" {
		return "", this.Error("get prepay_id result error: " + res.Return_msg, nil)
	}

	if res.Result_code != "SUCCESS" {
		return "", this.Error("get prepay_id result error: " + res.Err_code + "-" + res.Err_code_des, nil)
	}
	if res.Prepay_id == "" {
		return "", this.Error("get prepay_id result error: not get prepay_id", nil )
	}
	return res.Prepay_id, nil
}

func (this WxController) GetJsapiParameters() map[string]string {
	jsApiObj := make(map[string]string)
	jsApiObj["appId"] = this.appid
	jsApiObj["timeStamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	jsApiObj["nonceStr"] = this.CreateNoncestr(32)
	jsApiObj["package"] = "prepay_id=" + this.prepay_id
	jsApiObj["signType"] = "MD5"
	jsApiObj["paySign"] = this.GetSign(jsApiObj,this.key)
	return jsApiObj
}




