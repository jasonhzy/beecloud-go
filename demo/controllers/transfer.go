package controllers

import (
	"github.com/astaxie/beego"
	"beecloud-go/sdk"
	"reflect"
	"encoding/json"
	//"html/template"
	"strconv"
)

type TransferController struct {
	sdk.ApiController
	beego.Controller
}

//设置app id, app secret, master secret, test secret
func (this *TransferController) RegisterApp () {
	this.ApiController.RegisterApp(beego.AppConfig.String("app_id"), beego.AppConfig.String("app_secret"), beego.AppConfig.String("master_secret"), beego.AppConfig.String("test_secret"))
	this.ApiController.SetSandbox(false)
}

func (this TransferController) Print(strMsg string) {
	this.Ctx.WriteString(strMsg)
	this.StopRun()
}

func (this TransferController) Get() {
	this.RegisterApp()
	channel := this.GetString("type") //渠道

	timestamp := this.GetTimestamp()
	data := make(map[string]interface{})
	data["timestamp"] = timestamp

	transfer_no := "gotransfer" + strconv.FormatInt(timestamp, 10)
	var title string
	switch channel {
		case "WX_REDPACK" :
			title = "微信红包" //单个微信红包金额介于[1.00元，200.00元]之间
			data["transfer_no"] = transfer_no//微信要求10位数字
			data["channel_user_id"] = ""  //微信用户openid o3kKrjlROJ1qlDmFdlBQA95kvbN0
			data["channel"] = "WX_REDPACK"
			data["desc"] = "transfer test"

			redpack_info := make(map[string]interface{})
			redpack_info["send_name"] = "BeeCloud"
			redpack_info["wishing"] = "test"
			redpack_info["act_name"] = "test"
			data["redpack_info"] = redpack_info
		case "WX_TRANSFER" :
			title = "微信企业打款"
			data["transfer_no"] = transfer_no//微信企业打款为8-32位数字字母组合
			data["channel"] = "WX_TRANSFER"
			data["channel_user_id"] = ""   //微信用户openid o3kKrjlROJ1qlDmFdlBQA95kvbN0
			data["desc"] = "transfer test"
		case "ALI_TRANSFER" :
			title = "支付宝企业打款"
			data["channel"] = "ALI_TRANSFER"
			data["transfer_no"] = transfer_no

			//收款方的id 账号和 名字也需要对应
			data["channel_user_id"] = ""   //收款人账户
			data["channel_user_name"] = "" //收款人账户姓名
			data["account_name"] = "苏州比可网络科技有限公司" //注意此处需要和企业账号对应的全称
			data["desc"] = "transfer test"
		case "ALI_TRANSFERS" :
			title = "支付宝批量打款"
			data["channel"] = "ALI"
			data["desc"] = "transfer test"
			data["batch_no"] = transfer_no
			data["account_name"] = "苏州比可网络科技有限公司"

			d0 := make(map[string]interface{})
			d0["transfer_id"] = "bf693b3121864f3f969a3e1ebc5c376a"
			d0["receiver_account"] = "xxx" //收款方账户
			d0["receiver_name"] = "yyy"     //收款方账号姓名
			d0["transfer_fee"] = 1      //打款金额，单位为分
			d0["transfer_note"] = "note1"

			d1 := make(map[string]interface{})
			d1["transfer_id"] = "bf693b3121864f3f969a3e1ebc5c3768"
			d1["receiver_account"] = "xxx" //收款方账户
			d1["receiver_name"] = "yyy"     //收款方账号姓名
			d1["transfer_fee"] = 2      //打款金额，单位为分
			d1["transfer_note"] = "note2"

			var transfer_data []map[string]interface{}
			transfer_data = append(transfer_data, d0, d1)
			data["transfer_data"] = transfer_data
		case "BC_TRANSFER" :
			title = "BC企业打款"
			data["bill_no"] = transfer_no
			data["title"] = "GO DEMO测试BC企业打款"
			data["trade_source"] = "OUT_PC"
			/*
			 *  如果未能确认银行的全称信息,可通过下面的接口获取并进行确认
			 *  //P_DE:对私借记卡,P_CR:对私信用卡,C:对公账户
			 *
			 *	bank_params := make(map[string]interface{})
			 *	bank_params["type"] = "P_DE"
			 *
			 *	banks := this.get_banks(bank_params, channel)
			 *	var res map[string]interface{}
			 *	err := json.Unmarshal([]byte(banks), &res)
			 *	if err != nil {
			 *		this.Print(err.Error())
			 *	}
			 *	result_code := reflect.ValueOf(res["result_code"])
			 *	if result_code.IsValid() && result_code.Float() > 0 {
			 *		this.Print("banks result: " + reflect.ValueOf(res["err_detail"]).String())
			 *	}
			 *	fmt.Printf("%s", res["banks"])
			 */
			data["bank_fullname"] = "中国银行" //银行全称
			data["card_type"] = "DE" //银行卡类型,区分借记卡和信用卡，DE代表借记卡，CR代表信用卡，其他值为非法
			data["account_type"] = "P" //帐户类型，P代表私户，C代表公户，其他值为非法
			data["account_no"] = "6222691921993848888"   //收款方的银行卡号
			data["account_name"] = "test" //收款方的姓名或者单位名
			//选填mobile
			data["mobile"] = "" //银行绑定的手机号
			//选填optional
			optional := make(map[string]interface{})
			optional["company"] = "beecloud"
			data["optional"] = optional

			/**
			 * notify_url 选填，该参数是为打款成功之后接收返回信息配置的url,等同于在beecloud平台配置webhook，
			 * 如果两者都设置了，则优先使用notify_url。配置时请结合自己的项目谨慎配置，具体请
			 * 参考demo/webhook.php
			 */
			//data["notify_url"] = "http://beecloud.cn"
		case "CJ_TRANSFER" :
			title = "测试畅捷企业打款"
			data["bill_no"] = transfer_no
			data["title"] = "GO DEMO测试畅捷企业打款"
			/*
			 *  for bank_name, 支持的银行列表名称如下:
			 *  ICBC    中国工商银行      ABC     中国农业银行  BOC    中国银行
			 *  CCB     中国建设银行      COMM    交通银行      CMB    招商银行
			 *  CMBC    中国民生银行      CEB     中国光大银行  CIB   兴业银行
			 *  PSBC    中国邮政储蓄银行   GDB    广发银行      SPDB   上海浦东发展银行
			 *  SPDB    浦发银行          HXB    华夏银行
			 */
			data["bank_name"] = "中国银行" //银行全称
			data["card_type"] = "DEBIT" //银行卡类型,区分借记卡和信用卡，DEBIT代表借记卡，CREDIT代表信用卡，其他值为非法
			data["card_attribute"] = "B" //帐户类型，B代表公户，C代表私户，其他值为非法
			data["bank_account_no"] = "6222691921993848888"   //收款方的银行卡号
			data["bank_branch"] = "中国银行独墅湖支行"   //收款方的银行卡号
			data["account_name"] = "test" //收款方的姓名或者单位名
			data["province"] = "江苏省" //银行所在省份
			data["city"] = "苏州市" //银行所在市
			//选填optional
			optional := make(map[string]interface{})
			optional["company"] = "beecloud"
			data["optional"] = optional
		case "JD_TRANSFER" :
			title = "BC京东代付"
			data["bill_no"] = transfer_no
			data["channel"] = "JD_TRANSFER"
			data["title"] = "GO DEMO测试京东代付"
			data["trade_source"] = "OUT_PC"
			/*
			 *  如果未能确认银行的全称信息,可通过下面的接口获取并进行确认
			 *  //P_DE:对私借记卡,P_CR:对私信用卡,C:对公账户
			 *
			 *	bank_params := make(map[string]interface{})
			 *	bank_params["type"] = "P_DE"
			 *
			 *	banks := this.get_banks(bank_params, channel)
			 *	var res map[string]interface{}
			 *	err := json.Unmarshal([]byte(banks), &res)
			 *	if err != nil {
			 *		this.Print(err.Error())
			 *	}
			 *	result_code := reflect.ValueOf(res["result_code"])
			 *	if result_code.IsValid() && result_code.Float() > 0 {
			 *		this.Print("banks result: " + reflect.ValueOf(res["err_detail"]).String())
			 *	}
			 *	fmt.Printf("%s", res["banks"])
			 */
			data["bank_fullname"] = "中国银行" //银行全称
			data["card_type"] = "DE" //银行卡类型,区分借记卡和信用卡，DE代表借记卡，CR代表信用卡，其他值为非法
			data["account_type"] = "P" //帐户类型，P代表私户，C代表公户，其他值为非法
			data["account_no"] = "6226621808888888"   //收款方的银行卡号
			data["account_name"] = "test" //收款方的姓名或者单位名
			//选填mobile
			data["mobile"] = "" //银行绑定的手机号
			//选填optional
			optional := make(map[string]interface{})
			optional["company"] = "beecloud"
			data["optional"] = optional
		default :
			this.Print("No this channel" + channel)
	}

	this.RegisterApp()
	var result []byte
	var transferErr error
	switch channel {
		case "ALI_TRANSFERS":
			result, transferErr = this.Transfers(data)
		case "BC_TRANSFER", "JD_TRANSFER":
			result, transferErr = this.Bc_transfer(data)
		case "CJ_TRANSFER":
			result, transferErr = this.Cj_transfer(data)
		default :
			result, transferErr = this.Transfer(data)
	}
	if transferErr != nil {
		this.Print(transferErr.Error())
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
		this.Print("transfer result: " + reflect.ValueOf(res["err_detail"]).String())
	}

	url := reflect.ValueOf(res["url"])
	if url.IsValid() && url.String() != "" {
		this.Redirect(url.String(), 302)
	}else{
		this.Print(title + "打款成功")
	}
}

//BeePay自动打款-打款到银行卡
func (this TransferController) TransferTobank() {
	this.RegisterApp()

	data := make(map[string]interface{})
	//data["app_id"] = beego.AppConfig.String("app_id")
	data["withdraw_amount"] = 1 //必须是正整数，单位为分
	//商户订单号, 8到32位数字和/或字母组合，请自行确保在商户系统中唯一，同一订单号不可重复提交，否则会造成订单重复

	timestamp := this.GetTimestamp()
	data["bill_no"] = "gotransfer" + strconv.FormatInt(timestamp, 10)

	data["transfer_type"] = "1" //"1"代表对私打款，"2"代表对公打款
	data["bank_name"] = "中国银行" //银行全称, 不能写银行的缩写
	data["bank_account_no"] = "622269192199384xxxx" //收款方的银行卡号
	data["bank_account_name"] = "刘" //收款方的姓名或者单位名
	data["bank_code"] = "BOC" //银行的标准编码
	data["note"] = "测试"   //用户付款原因
	//选填optional
	optional := make(map[string]interface{})
	optional["company"] = "beecloud"
	data["optional"] = optional
	//选填notify_url，商户可通过此参数设定回调地址，此地址会覆盖用户在控制台设置的回调地址。必须以http://或https://开头
	//$data["notify_url"] = ""

	/*
	 * 对关键参数的签名，签名方式为MD5（32位小写字符）, 编码格式为UTF-8
	 * 验签规则即：app_id + bill_no + withdraw_amount + bank_account_no + master_secret的MD5生成的签名
	 * 其中master_secret为用户创建Beecloud App时获取的参数。
	 */
	//app_id := beego.AppConfig.String("app_id")
	//bill_no := reflect.ValueOf(data["bill_no"]).String()
	//withdraw_amount := reflect.ValueOf(data["withdraw_amount"]).String()
	//bank_account_no := reflect.ValueOf(data["bank_account_no"]).String()
	//
	//signStr := app_id + bill_no + withdraw_amount + bank_account_no + beego.AppConfig.String("master_secret")
	//sign := md5.New()
	//sign.Write([]byte(signStr))
	//data["signature"] = hex.EncodeToString(sign.Sum(nil))

	/*
	 *  返回结果:json格式，错误码（错误详细信息 参考err_detail字段)，如下所示：
	 *  result_code result_msg      含义
	 *   0	        OK	            调用成功
	 *   1	        APP_INVALID	    根据app_id找不到对应的APP或者app_sign不正确
	 *   4	        MISS_PARAM	    缺少必填参数
	 *   5	        PARAM_INVALID	参数不合法
	 *   14	        RUNTIME_ERROR	运行时错误
	 *   15	        NETWORK_ERROR	网络异常错误
	 */
	result, transferErr := this.Gateway_transfer(data)
	if transferErr != nil {
		this.Print(transferErr.Error())
	}
	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		//fmt.Printf("%s", res)
		this.Print("transfer to bank result: " + reflect.ValueOf(res["err_detail"]).String())
	}
	this.Print("打款成功, 打款记录唯一标识: " + reflect.ValueOf(res["id"]).String())
}

func (this TransferController) BcT1Transfer (){
	//设置app id, app secret, master secret, test secret
	this.RegisterApp()

	data := make(map[string]interface{})
	//data["app_id"] = beego.AppConfig.String("app_id")
	data["total_fee"] = 1 //必须是正整数，单位为分
	//商户订单号, 8到32位数字和/或字母组合，请自行确保在商户系统中唯一，同一订单号不可重复提交，否则会造成订单重复
	timestamp := this.GetTimestamp()
	data["bill_no"] = "gotransfer" + strconv.FormatInt(timestamp, 10)

	data["is_personal"] = "0" //"1"代表对私打款，"0"代表对公打款
	data["bank_account_no"] = "622269192199384xxxx" //收款方的银行卡号
	data["bank_account_name"] = "刘xx" //收款方的姓名或者单位名
	/*
	 * 对关键参数的签名，签名方式为MD5（32位小写字符）, 编码格式为UTF-8
	 * 验签规则即：app_id + bill_no + total_fee + bank_account_no + master_secret的MD5生成的签名
	 * 其中master_secret为用户创建Beecloud App时获取的参数。
	 */
	//app_id := beego.AppConfig.String("app_id")
	//bill_no := reflect.ValueOf(data["bill_no"]).String()
	//total_fee := reflect.ValueOf(data["total_fee"]).String()
	//bank_account_no := reflect.ValueOf(data["bank_account_no"]).String()
	//
	//signStr := app_id + bill_no + total_fee + bank_account_no + beego.AppConfig.String("master_secret")
	//sign := md5.New()
	//sign.Write([]byte(signStr))
	//data["signature"] = hex.EncodeToString(sign.Sum(nil))

	//获取银行列表
	//banks := this.get_banks(make(map[string]interface{}), "T1_EXPRESS_TRANSFER")
	//var res map[string]interface{}
	//err := json.Unmarshal([]byte(banks), &res)
	//if err != nil {
	//	this.Print(err.Error())
	//}
	//result_code := reflect.ValueOf(res["result_code"])
	//if result_code.IsValid() && result_code.Float() > 0 {
	//	this.Print("banks result: " + reflect.ValueOf(res["err_detail"]).String())
	//}
	//fmt.Printf("%s", res["banks"])


	data["bank_name"] = "中国工商银行" //银行全称, 不能写银行的缩写
	//选填optional
	optional := make(map[string]interface{})
	optional["company"] = "beecloud"
	data["optional"] = optional

	/*
	 *  返回结果:json格式，错误码（错误详细信息 参考err_detail字段)，如下所示：
	 *  result_code result_msg      含义
	 *   0	        OK	            调用成功
	 *   1	        APP_INVALID	    根据app_id找不到对应的APP或者app_sign不正确
	 *   4	        MISS_PARAM	    缺少必填参数
	 *   5	        PARAM_INVALID	参数不合法
	 *   14	        RUNTIME_ERROR	运行时错误
	 *   15	        NETWORK_ERROR	网络异常错误
	 */
	result, transferErr := this.Bct1_transfer(data)
	if transferErr != nil {
		this.Print(transferErr.Error())
	}

	var res map[string]interface{}
	err := json.Unmarshal([]byte(result), &res)
	if err != nil {
		this.Print(err.Error())
	}
	result_code := reflect.ValueOf(res["result_code"])
	if result_code.IsValid() && result_code.Float() > 0 {
		this.Print("transfer t1 result: " + reflect.ValueOf(res["err_detail"]).String())
	}
	this.Print("打款成功, 打款记录唯一标识: " + reflect.ValueOf(res["id"]).String())
}




