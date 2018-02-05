package sdk

import (
	//"github.com/astaxie/beego"
	"reflect"
	"crypto/md5"
	"encoding/hex"
)

type ApiController struct {
	BaseController
}

/*********************************订单支付／查询／退款**************************************/

func (this ApiController) Bill(data map[string]interface{}) ([]byte, error){
	var request_url string
	var secret_type int = 0
	if this.GetSandbox() {
		request_url = uri_test_bill //beego.AppConfig.String("test::uri_bill")
		secret_type = 2
	}else{
		request_url = uri_bill//beego.AppConfig.String("bc::uri_bill")
	}
	this.GetCommonParams(data, secret_type)
	this.GetSdkVersion(data)
	return this.http_post(request_url, data)
}

func (this ApiController) Bills(data map[string]interface{}) ([]byte, error){
	var request_url string
	var secret_type int = 0
	if this.GetSandbox() {
		request_url = uri_test_bills //beego.AppConfig.String("test::uri_bills")
		secret_type = 2
	}else{
		request_url = uri_bills//beego.AppConfig.String("bc::uri_bills")
	}
	this.GetCommonParams(data, secret_type)

	return this.http_get(request_url, data)
}

func (this ApiController) Bills_count(data map[string]interface{}) ([]byte, error){
	var request_url string
	var secret_type int = 0
	if this.GetSandbox() {
		request_url = uri_test_bills_count //beego.AppConfig.String("test::uri_bills_count")
		secret_type = 2
	}else{
		request_url = uri_bills_count //beego.AppConfig.String("bc::uri_bills_count")
	}
	this.GetCommonParams(data, secret_type)

	return this.http_get(request_url, data)
}

func (this ApiController)  Bill_query_byid(bill_id string) ([]byte, error){
	var request_url string
	var secret_type int = 0
	if this.GetSandbox() {
		request_url = uri_test_bill //beego.AppConfig.String("test::uri_bill")
		secret_type = 2
	}else{
		request_url = uri_bill //beego.AppConfig.String("bc::uri_bill")
	}
	data := make(map[string]interface{})
	this.GetCommonParams(data, secret_type)

	return this.http_get(request_url + "/" + bill_id, data)
}

func  (this ApiController) Refund(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_refund"), data)
	return this.http_post(uri_refund, data)
}

func (this ApiController) Refunds(data map[string]interface{}) ([]byte, error){
	this.GetCommonParams(data, 0)
	//return this.http_get(beego.AppConfig.String("uri_refunds"), data)
	return this.http_get(uri_refunds, data)
}

func (this ApiController) Refunds_count(data map[string]interface{}) ([]byte, error){
	this.GetCommonParams(data, 0)
	//return this.http_get(beego.AppConfig.String("uri_refunds_count"), data)
	return this.http_get(uri_refunds_count, data)
}

func (this ApiController)  Refund_query_byid(refund_id string) ([]byte, error){
	data := make(map[string]interface{})
	this.GetCommonParams(data, 0)
	//return this.http_get(beego.AppConfig.String("uri_refund") + "/" + refund_id, data)
	return this.http_get(uri_refund + "/" + refund_id, data)
}

func (this ApiController) Refund_status(data map[string]interface{}) ([]byte, error){
	this.GetCommonParams(data, 0)
	//return this.http_get(beego.AppConfig.String("uri_refund_status"), data)
	return this.http_get(uri_refund_status, data)
}

func (this ApiController) Offline_bill(data map[string]interface{}) ([]byte, error){
	this.GetCommonParams(data, 0)
	this.GetSdkVersion(data)
	//return this.http_post(beego.AppConfig.String("uri_offline_bill"), data)
	return this.http_post(uri_offline_bill, data)
}

func (this ApiController) Offline_refund (data map[string]interface{}) ([]byte, error){
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_offline_refund"), data)
	return this.http_post(uri_offline_refund, data)
}

func (this ApiController) Offline_bill_status(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_offline_bill_status"), data)
	return this.http_post(uri_offline_bill_status, data)
}

func (this ApiController) International_bill(data map[string]interface{}) ([]byte, error){
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_international_bill"), data)
	return this.http_post(uri_international_bill, data)
}

/**
 * @desc: 认证支付－确认支付
 *
 * @param $data
 *   token 渠道返回的token
 *   bc_bill_id  BeeCloud生成的唯一支付记录id
 *   verify_code 短信验证码
 *
 * @return json
 * @author: jason
 * @since: 2016-09-01
 */
func (this ApiController) Pay_confirm(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_pay_confirm"), data)
	return this.http_post(uri_pay_confirm, data)
}

/**
 * @desc: 签约API
 *
 * @param $data
 *   mobile 手机号
 *   bank  银行名称
 *   id_no 身份证号
 *   name   姓名
 *   card_no 银行卡号(借记卡,不支持信用卡)
 *   sms_id  获取验证码接口返回验证码记录的唯一标识
 *   sms_code 手机端接收到验证码
 *
 * @return json
 * @author: jason
 * @since: 2016-09-01
 */
func (this ApiController) Card_charge_sign(data map[string]interface{}) ([]byte, error)  {
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_card_charge_sign"), data)
	return this.http_post(uri_card_charge_sign, data)
}

/***************************微信、支付宝、BC企业打款等************************************/

//单笔打款 - 支付宝/微信红包
func (this ApiController) Transfer(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 1)
	//return this.http_post(beego.AppConfig.String("uri_transfer"), data)
	return this.http_post(uri_transfer, data)
}

//批量打款 - 支付宝
func (this ApiController) Transfers(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 1)
	//return this.http_post(beego.AppConfig.String("uri_transfers"), data)
	return this.http_post(uri_transfers, data)
}

//BC企业打款 - 支持bank
func (this ApiController) Bc_transfer_banks(data map[string]interface{}) ([]byte, error) {
	//return this.http_get(beego.AppConfig.String("uri_bc_transfer_banks"), data)
	return this.http_get(uri_bc_transfer_banks, data)
	
}

//BC企业打款 - 银行卡
func (this ApiController) Bc_transfer(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 1)

	var request_url string
	request_url = uri_bc_transfer// beego.AppConfig.String("uri_bc_transfer")
	channel := reflect.ValueOf(data["channel"]).String()
	if channel == "JD_TRANSFER" {
		request_url = uri_jd_transfer //beego.AppConfig.String("uri_jd_transfer")
	}
	return this.http_post(request_url, data)
}

//畅捷企业打款
func (this ApiController) Cj_transfer(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 1)
	//return this.http_post(beego.AppConfig.String("uri_cj_transfer"), data)
	return this.http_post(uri_cj_transfer, data)
}

//BeePay自动打款 - 打款到银行卡
func (this ApiController) Gateway_transfer(data map[string]interface{}) ([]byte, error) {
	app_id := reflect.ValueOf(data["app_id"]).String()
	if(app_id == ""){
		data["app_id"] = app_id //beego.AppConfig.String("app_id")
	}

	/*
	 * 对关键参数的签名，签名方式为MD5（32位小写字符）, 编码格式为UTF-8
	 * 验签规则即：app_id + bill_no + withdraw_amount + bank_account_no + master_secret的MD5生成的签名
	 * 其中master_secret为用户创建Beecloud App时获取的参数。
	 */
	sign := reflect.ValueOf(data["signature"]).String()
	if sign == "" {
		app_id := reflect.ValueOf(data["app_id"]).String()
		bill_no := reflect.ValueOf(data["bill_no"]).String()
		withdraw_amount := reflect.ValueOf(data["withdraw_amount"]).String()
		bank_account_no := reflect.ValueOf(data["bank_account_no"]).String()

		signStr := app_id + bill_no + withdraw_amount + bank_account_no + this.master_secret
		sign := md5.New()
		sign.Write([]byte(signStr))
		data["signature"] = hex.EncodeToString(sign.Sum(nil))
	}
	//return this.http_post(beego.AppConfig.String("uri_gateway_transfer"), data)
	return this.http_post(uri_gateway_transfer, data)
}

//T1代付接口
func (this ApiController) Bct1_transfer( data map[string]interface{}) ([]byte, error) {
	app_id := reflect.ValueOf(data["app_id"]).String()
	if(app_id == ""){
		data["app_id"] = app_id // beego.AppConfig.String("app_id")
	}
	/*
	 * 对关键参数的签名，签名方式为MD5（32位小写字符）, 编码格式为UTF-8
	 * 验签规则即：app_id + bill_no + total_fee + bank_account_no的MD5生成的签名
	 */
	sign := reflect.ValueOf(data["signature"]).String()
	if sign == "" {
		app_id := reflect.ValueOf(data["app_id"]).String()
		bill_no := reflect.ValueOf(data["bill_no"]).String()
		total_fee := reflect.ValueOf(data["total_fee"]).String()
		bank_account_no := reflect.ValueOf(data["bank_account_no"]).String()

		signStr := app_id + bill_no + total_fee + bank_account_no + this.master_secret
		sign := md5.New()
		sign.Write([]byte(signStr))
		data["signature"] = hex.EncodeToString(sign.Sum(nil))
	}
	//return this.http_post(beego.AppConfig.String("uri_t1_express_transfer"), data)
	return this.http_post(uri_t1_express_transfer, data)
}


func (this ApiController) Get_banks(data map[string]interface{}, strType string) ([]byte, error) {
	this.GetCommonParams(data, 0)

	switch strType {
		case "BC_GATEWAY":
			//return this.http_get(beego.AppConfig.String("uri_bc_gateway_banks"), data)
			return this.http_get(uri_bc_gateway_banks, data)
		case "T1_EXPRESS_TRANSFER":
			//return this.http_get(beego.AppConfig.String("uri_t1_express_transfer_banks"), data)
			return this.http_get(uri_t1_express_transfer_banks, data)
		default: //BC企业打款 - 支持bank
			//return this.http_get(beego.AppConfig.String("uri_bc_transfer_banks"), data)
			return this.http_get(uri_bc_transfer_banks, data)
	}
}

/******************************获取手机验证码************************************/
/*
 * @desc 发送短信验证码,返回验证码记录的唯一标识,并且手机端接收到验证码,二者供创建subscription使用
 * @param array $data, 主要包含以下四个参数:
 *  app_id string APP ID
 *  timestamp long 时间戳
 *  app_sign string 签名验证
 *  phone string 手机号
 * @return json:
 * 	result_code string
 *  result_msg string
 *  err_detail string
 *  sms_id string
 */
func (this ApiController) Sms(data map[string]interface{}) ([]byte, error) {
	this.GetCommonParams(data, 0)
	//return this.http_post(beego.AppConfig.String("uri_sms"), data)
	return this.http_post(uri_sms, data)
}



