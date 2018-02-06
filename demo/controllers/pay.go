package controllers

import (
	"github.com/astaxie/beego"
	"beecloud-go/sdk"
	wxClass "beecloud-go/sdk/common/wx"
	"reflect"
	"encoding/json"
	//"html/template"
	"strconv"
	"time"
	"strings"
)

type PayController struct {
	beego.Controller
	sdk.ApiController
}

func (this *PayController) RegisterApp () {
	this.BaseController.RegisterApp(beego.AppConfig.String("app_id"), beego.AppConfig.String("app_secret"), beego.AppConfig.String("master_secret"), beego.AppConfig.String("test_secret"));
	this.BaseController.SetSandbox(false);
}

func (this PayController) Print(strMsg string) {
	this.Ctx.WriteString(strMsg)
	this.StopRun()
}

func (this PayController) printJson (jsonObj map[string]interface{}){
	this.Data["json"] = jsonObj
	this.ServeJSON()
	return
}

func (this PayController) GetClientIp() string {
	ip := this.Ctx.Request.Header.Get("Remote_addr")
	if ip == "" {
		ip = this.Ctx.Request.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		ip = this.Substr(ip, 0, strings.Index(ip, ":"))
	}
	return ip
}

func (this *PayController) ToPay() {
	this.RegisterApp()
	channel := this.GetString("type") //渠道
	timestamp := this.GetTimestamp()
	bill := make(map[string]interface{})
	bill["timestamp"] = timestamp
	bill["bill_no"] = "godemo" + strconv.FormatInt(timestamp, 10)
	//total_fee(int 类型) 单位分
	bill["total_fee"] = 1
	//title UTF8编码格式，32个字节内，最长支持16个汉字
	bill["title"] = "GO " + channel  + "支付测试"
	//渠道类型:ALI_WEB 或 ALI_QRCODE 或 UN_WEB或JD_WAP或JD_WEB, BC_GATEWAY为京东、BC_WX_WAP、BC_ALI_WEB渠道时为必填, BC_ALI_WAP不支持此参数
	bill["return_url"] = "https://beecloud.cn"
	//选填 optional, 附加数据, eg: {"key1”:“value1”,“key2”:“value2”}
	//用户自定义的参数，将会在webhook通知中原样返回，该字段主要用于商户携带订单的自定义数据
	optional := make(map[string]interface{})
	optional["company"] = "beecloud";
	bill["optional"] = optional
	//选填 订单失效时间bill_timeout
	//必须为非零正整数，单位为秒，建议最短失效时间间隔必须大于360秒
	//京东(JD*)不支持该参数。
	//bill["bill_timeout"] = 360

	/**
	 * notify_url 选填，该参数是为接收支付之后返回信息的,仅适用于线上支付方式, 等同于在beecloud平台配置webhook，
	 * 如果两者都设置了，则优先使用notify_url。配置时请结合自己的项目谨慎配置，具体请
	 * 参考demo/webhook.php
	 *
	 */
	//bill["notify_url"] = "https://beecloud.cn"

	/**
	 * buyer_id选填
	 * 商户为其用户分配的ID.可以是email、手机号、随机字符串等；最长64位；在商户自己系统内必须保证唯一。
	 */
	//bill["buyer_id"] = "xxxx"

	/*
	 * coupon_id string 选填 卡券id
	 * 传入卡券id，下单时会自动扣除优惠金额再发起支付
	 */
	//bill["coupon_id"] = "xxxxxx"

	/**
	 * analysis选填, 分析数据
	 * 用于统计分析的数据，将会在控制台的统计分析报表中展示，用户自愿上传。包括以下基本字段：
	 *      os_name(系统名称，如"iOS"，"Android") os_version(系统版本，如"5.1") model(手机型号，如"iPhone 6")
	 *      app_name(应用名称) app_version(应用版本号) device_id(设备ID) category(类别，用户可自定义，如游戏分发渠道，门店ID等)
	 *      browser_name(浏览器名称) browser_version(浏览器版本)
	 * 下单产品保存格式：
	 *      product固定的key，
	 *      name 产品名，eg: T恤
	 *      count 产品件数, eg: 1
	 *      price 单价（单位分）,eg : 200
	 *  {"product":[{"name" : "xxx", "count":1, "price" : 111}, {"name" : "yyy", "count":2, "price" : 222}]}
	 */
	//type Product struct{
	//	Name string	`json:"name"`
	//	Count int	`json:"count"`
	//	Price int	`json:"price"`
	//}
	//type Prods struct {
	//	ProdItem []Product
	//}
	//var m Prods
	//m.ProdItem = append(m.ProdItem, Product{Name: "xxx", Count: 1, Price: 111})
	//m.ProdItem = append(m.ProdItem, Product{Name: "yyy", Count: 2, Price: 222})
	//
	//detail := make(map[string]interface{})
	//detail["product"] = m.ProdItem
	//detail["key"] = "value"
	//bill["analysis"] = detail

	switch channel {
		case "ALI_WEB" : //支付宝即时到账
			bill["channel"] = "ALI_WEB"
		case "ALI_WAP" : //支付宝移动网页
			bill["channel"] = "ALI_WAP"
			//非必填参数,boolean型,是否使用APP支付,true使用,否则不使用
			//bill["use_app"] = false
		case "ALI_QRCODE" : //支付宝扫码支付
			bill["channel"] = "ALI_QRCODE"
			//qr_pay_mode必填 二维码类型含义
			//0： 订单码-简约前置模式, 对应 iframe 宽度不能小于 600px, 高度不能小于 300px
			//1： 订单码-前置模式, 对应 iframe 宽度不能小于 300px, 高度不能小于 600px
			//3： 订单码-迷你前置模式, 对应 iframe 宽度不能小于 75px, 高度不能小于 75px
			bill["qr_pay_mode"] = "0"
		case "ALI_SCAN" : //支付宝刷卡
			bill["channel"] = "ALI_SCAN"
			bill["auth_code"] = "28886955594xxxxxxxx"
		case "ALI_OFFLINE_QRCODE" : //支付宝线下扫码
			bill["channel"] = "ALI_OFFLINE_QRCODE"
		case "BD_WEB" : //百度网页支付
			bill["channel"] = "BD_WEB"
		case "BD_WAP" : //百度移动网页
			bill["channel"] = "BD_WAP"
		case "JD_B2B" : //京东B2B
			bill["channel"] = "JD_B2B"
			/*
			 * bank_code(int 类型) for channel JD_B2B
			9102    中国工商银行      9107    招商银行
			9103    中国农业银行      9108    光大银行
			9104    交通银行          9109    中国银行
			9105    中国建设银行		9110 	 平安银行
			*/
			bill["bank_code"] = 9102
		case "JD_WEB" : //京东网页
			bill["channel"] = "JD_WEB"
		case "JD_WAP" : //京东移动网页
			bill["channel"] = "JD_WAP"
		case "UN_WEB" ://银联网页
			bill["channel"] = "UN_WEB"
		case "UN_WAP" : //银联移动网页, 由于银联做了适配,需在移动端打开,PC端仍显示网页支付
			bill["channel"] = "UN_WAP"
		case "WX_NATIVE": //微信扫码
			bill["channel"] = "WX_NATIVE"
		case "WX_JSAPI": //微信公众号
			bill["channel"] = "WX_JSAPI"
			//获取openid
			this.JsApiPay(bill, channel)
		case "WX_WAP": //微信H5网页, 请在手机浏览器内测试
			bill["channel"] = "WX_WAP"
			//需要参数终端ip，格式如下：
			analysis := make(map[string]interface{})
			analysis["ip"] = this.GetClientIp()
			bill["analysis"] = analysis
		case "WX_SCAN" :
			bill["channel"] = "WX_SCAN";
			bill["auth_code"] = "13022657110xxxxxxxx"
		case "YEE_WEB" : //易宝网页
			bill["channel"] = "YEE_WEB"
		case "YEE_WAP" : //易宝移动网页
			bill["channel"] = "YEE_WAP"
			bill["identity_id"] = "xxxxxxxxxxxxxx"
		case "YEE_NOBANKCARD": //易宝点卡支付
			//total_fee(订单金额)必须和充值卡面额相同，否则会造成金额丢失(渠道方决定)
			bill["total_fee"] = 10
			bill["channel"] = "YEE_NOBANKCARD"
			//点卡卡号，每种卡的要求不一样
			bill["cardno"] = "622662180018xxxx"
			//点卡密码，简称卡密
			bill["cardpwd"] = "xxxxxxxxxxxxxx"
			/*
			 * frqid 点卡类型编码
			 * 骏网一卡通(JUNNET),盛大卡(SNDACARD),神州行(SZX),征途卡(ZHENGTU),Q币卡(QQCARD),联通卡(UNICOM),
			 * 久游卡(JIUYOU),易充卡(YICHONGCARD),网易卡(NETEASE),完美卡(WANMEI),搜狐卡(SOHU),电信卡(TELECOM),
			 * 纵游一卡通(ZONGYOU),天下一卡通(TIANXIA),天宏一卡通(TIANHONG),32 一卡通(THIRTYTWOCARD)
			 */
			bill["frqid"] = "SZX"
		case "KUAIQIAN_WEB" : //快钱移动网页
			bill["channel"] = "KUAIQIAN_WEB"
		case "KUAIQIAN_WAP" : //快钱移动网页
			bill["channel"] = "KUAIQIAN_WEB"
		case "PAYPAL_PAYPAL" : //Paypal网页
			bill["channel"] = "PAYPAL_PAYPAL"
			/*
			 * currency参数的对照表, 请参考:
			 * https://github.com/beecloud/beecloud-rest-api/tree/master/international
			 */
			bill["currency"] = "USD"
		case "PAYPAL_CREDITCARD" : //Paypal信用卡
			bill["channel"] = "PAYPAL_CREDITCARD"
			/*
			 * currency参数的对照表, 请参考:
			 * https://github.com/beecloud/beecloud-rest-api/tree/master/international
			 */
			bill["currency"] = "USD"

			card_info := make(map[string]interface{})
			card_info["card_number"] = ""
			card_info["expire_month"] = 1
			card_info["expire_year"] = 2016
			card_info["cvv"] = 0
			card_info["first_name"] = ""
			card_info["last_name"] = ""
			card_info["card_type"] = "visa"
			bill["credit_card_info"] = card_info
		case "PAYPAL_SAVED_CREDITCARD" : //Paypal快捷
			bill["channel"] = "PAYPAL_SAVED_CREDITCARD"
			/*
			 * currency参数的对照表, 请参考:
			 * https://github.com/beecloud/beecloud-rest-api/tree/master/international
			 */
			bill["currency"] = "USD"
			bill["credit_card_id"] = ""
		case "BC_GATEWAY" : //BC网关支付
			bill["channel"] = "BC_GATEWAY"
			/*
			 * card_type(string 类型) for channel BC_GATEWAY
			 * 卡类型: 1代表信用卡；2代表借记卡
			*/
			bill["card_type"] = "1"
			bill["bank"] = "交通银行"
		case "BC_EXPRESS" : //快捷支付
			bill["channel"] = "BC_EXPRESS"
			//银行卡卡号, (选填，注意：可能必填，根据信息提示进行调整)
			//bill["card_no"] = "622662183243xxxx"
		case "BC_NATIVE" : //BC微信扫码
			bill["channel"] = "BC_NATIVE"
		case "BC_WX_WAP" : //BC微信H5支付, 请在手机浏览器内测试
			bill["channel"] = "BC_WX_WAP"
			//需要参数终端ip，格式如下：
			analysis := make(map[string]interface{})
			analysis["ip"] = this.GetClientIp()
			bill["analysis"] = analysis
		case "BC_WX_SCAN" : //BC微信刷卡
			bill["channel"] = "BC_WX_SCAN"
			bill["auth_code"] = "13022657110xxxxxxxx"
		case "BC_WX_JSAPI": //微信公众号
			bill["channel"] = "BC_WX_JSAPI"
			//获取openid
			this.JsApiPay(bill, channel)
		case "BC_ALI_QRCODE" : //BC支付宝线下扫码
			bill["channel"] = "BC_ALI_QRCODE"
		case "BC_ALI_SCAN" : //BC支付宝刷卡
			bill["channel"] = "BC_ALI_SCAN"
			bill["auth_code"] = "28886955594xxxxxxxx"
		case "BC_ALI_WAP" : //请在手机浏览器内测试
			bill["channel"] = "BC_ALI_WAP"
		case "BC_QQ_NATIVE" :
			bill["channel"] = "BC_QQ_NATIVE"
		case "BC_JD_QRCODE" : //BC京东扫码
			bill["channel"] = "BC_JD_QRCODE"
		default:
			this.Print("No this type.")
	}


	//get openid
	if this.In_array(channel, []string{"WX_JSAPI", "BC_WX_JSAPI"}) {
		this.JsApiPay(bill, channel)
	}

	var result []byte
	var payErr error
	if this.In_array(channel, []string{"PAYPAL_PAYPAL", "PAYPAL_CREDITCARD", "PAYPAL_SAVED_CREDITCARD"}) {
		result, payErr = this.International_bill(bill)
	}else if this.In_array(channel, []string{"ALI_OFFLINE_QRCODE", "ALI_SCAN", "BC_ALI_SCAN", "WX_SCAN", "BC_WX_SCAN"}) {
		result, payErr = this.Offline_bill(bill)
	}else {
		result, payErr = this.Bill(bill)
	}
	if(payErr != nil){
		this.Print(payErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	//fmt.Printf("%s", res)
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		//fmt.Printf("%s", res)
		this.Print("pay result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	url := reflect.ValueOf(res["url"])
	html := reflect.ValueOf(res["html"])
	code_url := reflect.ValueOf(res["code_url"])
	credit_card_id := reflect.ValueOf(res["credit_card_id"])
	id := reflect.ValueOf(res["id"])

	if this.In_array(channel, []string{"WX_JSAPI", "BC_WX_JSAPI"}) { //微信公众号支付
		appId := reflect.ValueOf(res["app_id"])
		strPackage := reflect.ValueOf(res["package"])
		signType := reflect.ValueOf(res["sign_type"])
		paySign := reflect.ValueOf(res["pay_sign"])
		nonceStr := reflect.ValueOf(res["nonce_str"])
		timeStamp := reflect.ValueOf(res["timestamp"])

		if !appId.IsValid() || !strPackage.IsValid() || !signType.IsValid() || !paySign.IsValid()  || !timeStamp.IsValid() {
			this.Print("wx pay params invalid")
		}
		if appId.String() == "" || strPackage.String() == "" || signType.String() == "" || paySign.String() == ""  || timeStamp.Float() == 0 {
			this.Print("wx pay params is empty")
		}

		jsapiParameters := make(map[string]interface{})
		jsapiParameters["appId"] = appId.String()
		jsapiParameters["timeStamp"] = timeStamp.Float()
		jsapiParameters["nonceStr"] = nonceStr.String()
		jsapiParameters["package"] = strPackage.String()
		jsapiParameters["signType"] = signType.String()
		jsapiParameters["paySign"] = paySign.String()

		this.Data["jsapi"] = jsapiParameters
		this.Data["channel"] = "JSAPI"
		this.TplName = "pay.tpl"
	}else if url.IsValid() && url.String() != "" {
		this.Redirect(url.String(), 302)
	}else if code_url.IsValid() && code_url.String() != "" {
		if channel == "WX_NATIVE" || channel == "BC_NATIVE" || channel == "BC_ALI_QRCODE" || channel == "BC_QQ_NATIVE" || channel == "BC_JD_QRCODE" {
			this.Data["title"] = channel + "支付"
			this.Data["channel"] = channel
			this.Data["id"] = reflect.ValueOf(res["id"])
			this.Data["bill_no"] = reflect.ValueOf(bill["bill_no"])
			this.Data["code_url"] = code_url.String()
			this.TplName = "qrcode.tpl"
		}else{
			this.Redirect(code_url.String(), 302)
		}
	}else if html.IsValid() && html.String() != "" {
		//this.Data["channel"] = channel + "支付"
		//this.Data["content"] = template.HTML(html.String())
		//this.TplName = "pay.tpl"
		this.Print("<html><head><meta charset=\"UTF-8\"></head><body>" + html.String() + "</body></html>")
	}else if credit_card_id.IsValid() && credit_card_id.String() != "" {
		this.Print("信用卡id(PAYPAL_CREDITCARD): " + credit_card_id.String())
	}else if id.IsValid() && id.String() != "" {
		this.Print("支付成功:" + id.String())
	}
}

func (this *PayController) JsApiPay(bill map[string]interface{}, channel string) {
	//定义的三种方式
	//wx := &wxClass.WxController{}
	//wx := new(wxClass.WxController)
	wx := wxClass.WxController{}
	wx.RegisterWx()

	//获取openid
	var openid string
	var openidErr error

	wxCode := this.GetString("code")
	if wxCode == "" {
		this.Redirect(wx.CreateOauthUrlForCode("http://test3.beecloud.cn/wxpay/demo?type=" + channel), 302)
	}else{
		wx.SetCode(wxCode)
		openid, openidErr = wx.GetOpenid()
		if openidErr != nil {
			this.Print(openidErr.Error())
		}
	}
	//openid = "ofEy7uK2JsSOXVpHHYErRPtrdVWg"
	bill["openid"] = openid

	//wx.SetParameter("openid", openid)
	//wx.SetParameter("out_trade_no", reflect.ValueOf(bill["bill_no"]).String())
	//wx.SetParameter("total_fee", strconv.FormatInt(reflect.ValueOf(bill["total_fee"]).Int(), 10))
	//wx.SetParameter("trade_type", "JSAPI")
	//wx.SetParameter("body", reflect.ValueOf(bill["title"]).String())
	//wx.SetParameter("notify_url", "http://test3.beecloud.cn/notify")
	//wx.SetParameter("spbill_create_ip", this.GetClientIp())

	//prepay_id, err := wx.GetPrepayId()
	//if err != nil {
	//	this.Print(err.Error())
	//}
	//wx.SetPrepayId(prepay_id)

	//this.Data["jsapi"] = wx.GetJsapiParameters()
	//this.Data["channel"] = "JSAPI"
	//this.TplName = "pay.tpl"
}

func (this PayController) BillStatus() {
	this.RegisterApp()
	bill_no := this.GetString("bill_no")
	channel := this.GetString("channel")

	conditions := make(map[string]interface{})
	conditions["bill_no"] = bill_no
	conditions["channel"] = channel

	var result []byte
	var payErr error
	var res map[string]interface{}
	jsonObj := make(map[string]interface{})
	jsonObj["result_code"] = 1
	switch channel {
		case "BC_QQ_NATIVE", "BC_JD_QRCODE", "BC_NATIVE", "BC_ALI_QRCODE":
			result, payErr = this.Bills(conditions)
			if payErr != nil {
				jsonObj["err_detail"] = payErr.Error()
				this.printJson(jsonObj)
			}

			err := json.Unmarshal([]byte(result), &res)
			if err != nil {
				jsonObj["err_detail"] = err.Error()
				this.printJson(jsonObj)
			}
			result_code := reflect.ValueOf(res["result_code"])
			if result_code.IsValid() && result_code.Float() > 0 {
				jsonObj["err_detail"] = "bill status query result: " + reflect.ValueOf(res["err_detail"]).String()
				this.printJson(jsonObj)
			}

			var pay_result bool = false
			bills := res["bills"].([]interface{})
			if len(bills) > 0 {
				for _, bill := range bills  {
					info := bill.(map[string]interface{})
					if info["spay_result"].(bool) {
						pay_result = true
						break
					}
				}
			}
			jsonObj["result_code"] = 0
			jsonObj["pay_result"] = pay_result
		case "ALI_OFFLINE_QRCODE", "ALI_SCAN", "WX_SCAN", "WX_NATIVE":
			result, payErr = this.Offline_bill_status(conditions)
			if payErr != nil {
				jsonObj["err_detail"] = payErr.Error()
				this.printJson(jsonObj)
			}
			err := json.Unmarshal([]byte(result), &res)
			if err != nil {
				jsonObj["err_detail"] = err.Error()
				this.printJson(jsonObj)
			}
			result_code := reflect.ValueOf(res["result_code"])
			if result_code.IsValid() && result_code.Float() > 0 {
				jsonObj["err_detail"] = "bill status query result: " + reflect.ValueOf(res["err_detail"]).String()
				this.printJson(jsonObj)
			}
			jsonObj["result_code"] = 0
			jsonObj["pay_result"] = reflect.ValueOf(res["pay_result"]).Bool()
		default:
			this.Print("No this channel")
	}
	this.printJson(jsonObj)
}

/**
 *  根据具体的条件进行查询，可查询的字段以文档为准
 *  参考文档：https://beecloud.cn/doc/#4-4
 *
 */
func (this *PayController) BillQuery() {
	channel := this.GetString("type")
	var title string
	switch channel {
		case "ALI" :
			title = "支付宝"
		case "BD" :
			title = "百度"
		case "JD" :
			title = "京东"
		case "WX" :
			title = "微信"
		case "UN" :
			title = "银联"
		case "YEE" :
			title = "易宝"
		case "KUAIQIAN" :
		 	title = "快钱"
		case "PAYPAL" :
			title = "PAYPAL"
		case "BC" :
			title = "BC网关/快捷支付"
		default:
			this.Print("No this channel")
	}

	this.RegisterApp()

	conditions := make(map[string]interface{})
	conditions["channel"] = channel
	conditions["spay_result"] = true //只列出了支付成功的订单
	conditions["limit"] = 10

	var res map[string]interface{}
	result, queryErr := this.Bills(conditions)
	if queryErr != nil {
		this.Print(queryErr.Error())
	}

	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("bill query result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	result1, queryErr1 := this.Bills_count(conditions)
	if queryErr1 != nil {
		this.Print(queryErr1.Error())
	}

	var res1 map[string]interface{}
	err1 := json.Unmarshal([]byte(result1), &res1)
	if err1 != nil {
		this.Print(err1.Error())
	}
	result_code1 := reflect.ValueOf(res1["result_code"])
	if result_code1.IsValid() && result_code1.Float() > 0 {
		this.Print("bill count result: " + reflect.ValueOf(res1["err_detail"]).String())
	}

	this.Data["title"] = title + "支付"
	this.Data["data"] = res["bills"].([]interface{})
	this.Data["count"] = res1["count"].(float64)
	this.Data["type"] = "bill"
	this.TplName = "bills.tpl"
}

func (this *PayController) BillID() {
	bill_id := this.GetString("id")
	if bill_id == "" {
		this.Print("请输入订单唯一标识ID")
	}

	this.RegisterApp()

	result, queryErr := this.Bill_query_byid(bill_id)
	if queryErr != nil {
		this.Print(queryErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("bill query by id result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	//var data [1]map[string]interface{}
	//data[0] = res["pay"].(map[string]interface{})

	var data []map[string]interface{}
	data = append(data, res["pay"].(map[string]interface{}))

	this.Data["data"] = data
	this.Data["type"] = "bill"
	this.TplName = "bills.tpl"
}

func (this PayController) ToRefund(){
	str_channel := this.GetString("type") //渠道
	str_refund_fee := this.GetString("refund_fee")

	timestamp := this.GetTimestamp()
	data := make(map[string]interface{})
	data["timestamp"] = timestamp
	data["bill_no"] = this.GetString("bill_no")

	//refund_no退款单号,为(预)退款使用的, 格式为:退款日期(8位) + 流水号(3~24 位)。
	//请自行确保在商户系统中唯一，且退款日期必须是发起退款的当天日期, 同一退款单号不可重复提交，否则会造成退款单重复。
	//流水号可以接受数字或英文字符，建议使用数字，但不可接受“000”
	t := time.Now()
	t.Format("20060102")
	data["refund_no"] = t.Format("20060102") + strconv.FormatInt(timestamp, 10)

	refund_fee, err := strconv.ParseInt(str_refund_fee, 10, 64)
	if(err != nil){
		this.Print(err.Error())
	}

	data["refund_fee"] = refund_fee

	//选填 optional
	optional := make(map[string]interface{})
	optional["company"] = "beecloud"
	data["optional"] = optional

	//refund_account(类型Integer),适用于WX_NATIVE, WX_JSAPI, WX_SCAN, WX_APP
	//退款资金来源 1:可用余额退款 0:未结算资金退款（默认使用未结算资金退款）
	//data["refund_account"] = 1

	/**
	 * notify_url 选填，该参数是为退款成功之后接收返回信息配置的url,等同于在beecloud平台配置webhook，
	 * 如果两者都设置了，则优先使用notify_url。配置时请结合自己的项目谨慎配置，具体请
	 * 参考demo/webhook.php
	 */
	//data["notify_url"] = 'http://beecloud.cn'

	var title string
	switch str_channel {
		case "ALI" :
			title = "支付宝"
			data["channel"] = "ALI"
		case "BD" :
			title = "百度"
			data["channel"] = "BD"
		case "JD" :
			title = "京东"
			data["channel"] = "JD"
		case "WX" :
			title = "微信"
			data["channel"] = "WX"
		case "UN" :
			title = "银联"
			data["channel"] = "UN"
		case "YEE" :
			title = "易宝"
			data["channel"] = "YEE"
		case "KUAIQIAN" :
			data["channel"] = "KUAIQIAN"
			title = "快钱"
		case "BC" :
			title = "BC支付"
			data["channel"] = "BC"
		case "PAYPAL" :
			title = "PAYPAL"
			data["channel"] = "PAYPAL"
		case "ALI_OFFLINE_QRCODE", "ALI_SCAN":
			title = str_channel + "线下退款"
			data["channel"] = "ALI"
		case "WX_SCAN", "WX_NATIVE": //非服务商WX_NATIVE 可通过/rest/refund/ 或 /rest/offline/refund/进行退款
			title = str_channel + "线下退款"
			data["channel"] = "WX"
		case "BC_WX_SCAN", "BC_ALI_SCAN", "BC_ALI_QRCODE", "BC_NATIVE":
			title = str_channel + "线下退款"
			data["channel"] = "BC"
		default :
			this.Print("No this channel.")
	}

	this.RegisterApp()

	var result []byte
	var refundErr error
	switch str_channel {
		case "ALI_OFFLINE_QRCODE", "ALI_SCAN", "WX_SCAN", "WX_NATIVE", "BC_NATIVE", "BC_WX_SCAN", "BC_ALI_SCAN", "BC_ALI_QRCODE":
			result, refundErr = this.Offline_refund(data)
		default:
			result, refundErr = this.Refund(data)
	}
	if refundErr != nil {
		this.Print(refundErr.Error())
	}

	var res map[string]interface{}
	jsonerr := json.Unmarshal([]byte(result), &res)
	if jsonerr != nil {
		this.Print(jsonerr.Error())
	}

	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("bill refund result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	//当channel为ALI_APP、ALI_WEB、ALI_QRCODE，并且不是预退款
	if str_channel == "ALI" {
		this.Redirect(reflect.ValueOf(res["url"]).String(), 302)
	}

	this.Print(title + "退款成功, 退款表记录唯一标识ID:" + reflect.ValueOf(res["id"]).String() )
}

func (this *PayController) RefundQuery() {
	channel := this.GetString("type")
	var title string
	switch channel {
		case "ALI" :
			title = "支付宝"
		case "BD" :
			title = "百度"
		case "JD" :
			title = "京东"
		case "WX" :
			title = "微信"
		case "UN" :
			title = "银联"
		case "YEE" :
			title = "易宝"
		case "KUAIQIAN" :
			title = "快钱"
		case "PAYPAL" :
			title = "PAYPAL"
		case "BC" :
			title = "BC网关/快捷"
		default:
			this.Print("No this channel")
	}

	this.RegisterApp()

	conditions := make(map[string]interface{})
	conditions["channel"] = channel
	conditions["limit"] = 10

	result, refundErr := this.Refunds(conditions)
	if refundErr != nil {
		this.Print(refundErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("bill refund query result: " + reflect.ValueOf(res["err_detail"]).String())
	}


	result1, refundErr1 := this.Refunds_count(conditions)
	if refundErr1 != nil {
		this.Print(refundErr1.Error())
	}

	var res1 map[string]interface{}
	err1 := json.Unmarshal([]byte(result1), &res1)
	if err1 != nil {
		this.Print(err1.Error())
	}
	result_code1 := reflect.ValueOf(res1["result_code"])
	if result_code1.IsValid() && result_code1.Float() > 0 {
		this.Print("bill refund count result: " + reflect.ValueOf(res1["err_detail"]).String())
	}

	this.Data["title"] = title + "支付"
	this.Data["data"] = res["refunds"].([]interface{})
	this.Data["count"] = res1["count"].(float64)
	this.Data["type"] = "refund"
	this.TplName = "bills.tpl"
}

func (this *PayController) RefundID() {
	refund_id := this.GetString("id")
	if refund_id == "" {
		this.Print("请输入退款单唯一标识ID")
	}

	this.RegisterApp()

	result, refundErr := this.Refund_query_byid(refund_id)
	if refundErr != nil {
		this.Print(refundErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("refund query by id result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	//var data [1]map[string]interface{}
	//data[0] = res["refund"].(map[string]interface{})

	var data []map[string]interface{}
	data = append(data, res["refund"].(map[string]interface{}))

	this.Data["data"] = data
	this.Data["type"] = "refund"
	this.TplName = "bills.tpl"
}

func (this PayController) RefundStatus() {
	this.RegisterApp()
	refund_no := this.GetString("refund_no")
	channel := this.GetString("channel")

	conditions := make(map[string]interface{})
	conditions["refund_no"] = refund_no
	conditions["channel"] = channel

	var result []byte
	var refundErr error
	var res map[string]interface{}
	switch channel {
		case "WX", "YEE", "KUAIQIAN", "BD":
			result, refundErr = this.Refund_status(conditions)
			if refundErr != nil {
				this.Print(refundErr.Error())
			}

			err := json.Unmarshal([]byte(result), &res)
			if err != nil {
				this.Print(err.Error())
			}
			result_code := reflect.ValueOf(res["result_code"])
			if result_code.IsValid() && result_code.Float() > 0 {
				this.Print("refund status query result: " + reflect.ValueOf(res["err_detail"]).String())
			}

			this.printJson(res)
		default:
			this.Print("No this channel")
	}
}

func (this PayController) PayConfirm ()  {
	this.RegisterApp()

	timestamp := this.GetTimestamp()
	bill := make(map[string]interface{})
	bill["timestamp"] = timestamp
	bill["bill_no"] = "godemo" + strconv.FormatInt(timestamp, 10)

	//total_fee(int 类型) 单位分
	bill["total_fee"] = 1
	//title UTF8编码格式，32个字节内，最长支持16个汉字
	bill["title"] = "godemo认证支付测试"
	//支付渠道
	bill["channel"] = "BC_EXPRESS"

	/**
	 * buyer_id选填
	 * 商户为其用户分配的ID.可以是email、手机号、随机字符串等；最长64位；在商户自己系统内必须保证唯一。
	 */
	//bill["buyer_id"] = "xxx"

	//第一次发起支付时，在optional中传入phone_no(手机号)，card_no（银行卡号），id_no（身份证号），customer_name（银行卡持有者姓名）等四个要素，
	//第一次发起支付成功后，可以传入token（第一次发起支付时返回的授权码）一个要素即可
	optional := make(map[string]interface{})
	optional["id_no"] = "21302619870917xxxx"
	optional["customer_name"] = "xxx"
	optional["id_no"] = "21302619870917xxxx"
	optional["card_no"] = "622622180408xxxx"
	optional["phone_no"] = "1596214xxxx"
	//optional["token"] = "235BAFEF6039440C7045D1A5E972xxxx"
	bill["optional"] = optional

	//第一步: 获取验证码；得到参数token, 支付记录的唯一标识id以及手机上收到的短信验证码
	result, confirmErr := this.Bill(bill)
	if confirmErr != nil {
		this.Print(confirmErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("refund status query result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	//第二步: 确认支付；传入token、支付记录的id、短信验证码
	verify := make(map[string]interface{})
	verify["bc_bill_id"] = reflect.ValueOf(res["id"]).String() //BeeCloud生成的唯一支付记录id
	verify["token"] = reflect.ValueOf(res["token"]).String()		//渠道返回的token
	verify["verify_code"] = "yyyy" //短信验证码

	result1, confirmErr1 := this.Pay_confirm(verify)
	if confirmErr1 != nil {
		this.Print(confirmErr1.Error())
	}

	var res1 map[string]interface{}
	err1 := json.Unmarshal([]byte(result1), &res1)
	if err1 != nil {
		this.Print(err1.Error())
	}
	result_code1 := reflect.ValueOf(res["result_code"])
	if result_code1.IsValid() && result_code1.Float() > 0 {
		this.Print("pay confirm result: " + reflect.ValueOf(res1["err_detail"]).String())
	}
	this.Print("支付成功")
}

func (this PayController) CardCharge ()  {
	this.RegisterApp()

	timestamp := this.GetTimestamp()


	//第一步: 获取验证码；,得到参数sms_id, sms_code: 请查看手机收到的验证码
	sms := make(map[string]interface{})
	sms["hone"] = "159621xxxxx"
	sms["timestamp"] = timestamp
	sms_result, smsErr := this.Sms(sms)
	if smsErr != nil {
		this.Print(smsErr.Error())
	}

	var sms_res map[string]interface{}
	sms_err := json.Unmarshal([]byte(sms_result), &sms_res)
	if sms_err != nil {
		this.Print(sms_err.Error())
	}
	sms_result_code := reflect.ValueOf(sms_res["result_code"])
	if sms_result_code.IsValid() && sms_result_code.Float() > 0 {
		this.Print("card charge sign result: " + reflect.ValueOf(sms_res["err_detail"]).String())
	}

	//sms_id := reflect.ValueOf(sms_res["sms_id"]).String()

	/*
	* 第二步: 签约API, 配置webhook,签约成功之后, 获取到card_id(注意保存)
	* 具体参数含义如下:
	*   mobile 手机号
	*   bank  银行名称
	*   id_no 身份证号
	*   name   姓名
	*   card_no 银行卡号(借记卡,不支持信用卡)
	*   sms_id  获取验证码接口返回验证码记录的唯一标识
	*   sms_code 手机端接收到验证码
	*/
	sign := make(map[string]interface{})
	sign["timestamp"] = timestamp
	sign["mobile"] = "1596214xxxxx"
	sign["bank"] = "中国银行"
	sign["id_no"] = "413026xxxxxxxxxxx"
	sign["name"] = "jason"
	sign["card_no"] = "6226xxxxxxxxxxxxxxxx"
	sign["sms_id"] = "d4fb7cdd-13ff-4c6c-ac57-df5aee717988"
	sign["sms_code"] = "374932"

	sign_result, signErr := this.Card_charge_sign(sign)
	if signErr != nil {
		this.Print(signErr.Error())
	}

	var sign_res map[string]interface{}
	err1 := json.Unmarshal([]byte(sign_result), &sign_res)
	if err1 != nil {
		this.Print(err1.Error())
	}
	sign_result_code := reflect.ValueOf(sign_res["result_code"])
	if sign_result_code.IsValid() && sign_result_code.Float() > 0 {
		this.Print("card charge sign result: " + reflect.ValueOf(sign_res["err_detail"]).String())
	}
	//card_id := reflect.ValueOf(res["card_id"]).String()

	//第三步: 通过第二步取到card_id,进行支付
	bill := make(map[string]interface{})

	bill["card_id"] = "21e58f1c-de22-4979-a95c-3dfxxxxxx" //第二步签约获取
	bill["timestamp"] = timestamp
	bill["bill_no"] = "godemo" + strconv.FormatInt(timestamp, 10)

	//total_fee(int 类型) 单位分,  最小金额150分
	bill["total_fee"] = 1
	//title UTF8编码格式，32个字节内，最长支持16个汉字
	bill["title"] = "godemo支付测试"
	//渠道
	bill["channel"] = "BC_CARD_CHARGE"
	//渠道类型:ALI_WEB 或 ALI_QRCODE 或 UN_WEB或JD_WAP或JD_WEB时为必填
	bill["return_url"] = "https://beecloud.cn"

	optional := make(map[string]interface{})
	optional["company"] = "beecloud"
	bill["optional"] = optional

	result, payErr := this.Bill(bill)
	if payErr != nil {
		this.Print(payErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("pay result: " + reflect.ValueOf(res["err_detail"]).String())
	}
	this.Print("支付成功")
}



