package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"reflect"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"fmt"
)

type WebhookController struct {
	beego.Controller
}

func (this WebhookController) Print(strMsg string) {
	this.Ctx.WriteString(strMsg)
	this.StopRun()
}
/**
 * http类型为 Application/json, 非XMLHttpRequest的application/x-www-form-urlencoded, $_POST方式是不能获取到的，
 * APP ID和Master Secret可以在https://beecloud.cn平台登录后获取
 *
 * 备注：secret是一个非常重要的数据，请您必须小心谨慎的确保此数据保存在足够安全的地方。
 *      您从BeeCloud官方获得此数据的同时，即表明您保证不会被用于非法用途和不会在没有得到您授权的情况下被盗用，
 *      一旦因此数据保管不善而导致的经济损失及法律责任，均由您独自承担。
 */
func (this WebhookController) Post() {
	var data map[string]interface{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		print(err.Error())
	}
	//打印所有字段,仅供调试使用
	fmt.Printf("%s", data)

	// webhook字段文档: https://beecloud.cn/doc/?sdk=php#3-4或者https://github.com/beecloud/beecloud-webhook
	app_id := beego.AppConfig.String("app_id")
	master_secret := beego.AppConfig.String("master_secret")
	transaction_id := reflect.ValueOf(data["transaction_id"]).String()
	transaction_type := reflect.ValueOf(data["transaction_type"]).String()
	channel_type := reflect.ValueOf(data["channel_type"]).String()
	sub_channel := reflect.ValueOf(data["sub_channel_type"]).String()

	var transaction_fee float64 = 0
	if data["transaction_fee"] != nil {
		transaction_fee = reflect.ValueOf(data["transaction_fee"]).Float()
	}
	var status bool = false
	if data["trade_success"] != nil {
		status = reflect.ValueOf(data["trade_success"]).Bool()
	}

	//第一步:验证签名
	//签名算法：app_id + transaction_id + transaction_type + channel_type + transaction_fee + master_secret
	signStr := app_id + transaction_id + transaction_type + channel_type + strconv.FormatFloat(transaction_fee, 'g', 1, 64) + master_secret
	sign := md5.New()
	sign.Write([]byte(signStr))
	signRes := hex.EncodeToString(sign.Sum(nil))

	// 签名不正确
	if signRes != reflect.ValueOf(data["signature"]).String() {
		this.Print("签名错误")
	}

	//第二步:过滤重复的Webhook
	//客户需要根据订单号进行判重，忽略已经处理过的订单号对应的Webhook
	//if (transaction_id对应的订单号已经处理完毕){
	//  TODO...
	//}
	//

	//第三步:验证订单金额与购买的产品实际金额是否一致
	//也就是验证调用Webhook返回的transaction_fee订单金额是否与客户服务端内部的数据库查询得到对应的产品的金额是否相同
	//if (transaction_fee != 客户服务端查询得到的实际产品金额) {
	//  TODO...
	//}

	//第四步:处理业务逻辑和返回
	switch transaction_type {
		//推送支付的结果
		case "PAY":
			if status { //支付状态是否变为支付成功,true代表成功
				//TODO...
			}
			//message_detail 参考文档
			//channel_type 微信/支付宝/银联/快钱/京东/百度/易宝/PAYPAL/BC
			sub_channel := reflect.ValueOf(data["sub_channel_type"]).String()
			switch channel_type {
				case "WX":
				case "ALI":
				case "UN":
				case "KUAIQIAN":
				case "JD":
				case "BD":
				case "YEE":
				case "PAYPAL":
				case "BC":
					//BC订阅收费
					if sub_channel == "BC_SUBSCRIPTION" {

					}
					//BC代扣
					if sub_channel == "BC_CARD_CHARGE" {

					}
			}
		//退款的结果
		case "REFUND":
			if sub_channel == "BC_TRANSFER" {
				//message_detail中包含打款相关的详细信息
				//TODO...
			}
		/*
		 * 推送企业打款结果的
		 * transaction_id就是企业打款的交易单号, 对应支付请求的bill_no, transaction_type为TRANSFER, sub_channel_type为BC_TRANSFER
		 * message_detail中包含打款相关的详细信息
		 */
		case "TRANSFER":
			if status { //企业打款状态是否为成功,true代表成功
				//TODO...
			}
	    /*
		 * 推送订阅结果的
		 * transaction_id就是创建订阅时返回的订阅id，transaction_type为SUBSCRIPTION，sub_channel_type为BC_SUBSCRIPTION，
		 * message_detail中包含用户相关的注册信息，其中的card_id注意留存。
		 * 该id由{bank_name、card_no、id_name、id_no、mobile}共同决定，可以直接用于发起订阅
		 */
		case "SUBSCRIPTION":
			if status { //创建的订阅状态是否为成功,true代表成功
				//TODO...
			}
			if sub_channel == "BC_SUBSCRIPTION" {
				//message_detail中包含签约相关的详细信息，包括card_id
				//TODO...
			}
		/*
		 * 推送代扣签约结果的
		 * transaction_id就是代扣签约返回的id，transaction_type为SIGN，sub_channel_type为BC_CARD_CHARGE，
		 * message_detail中包含签约相关的详细信息，其中的card_id注意留存。
		 */
		case "SIGN":
			if status { //创建的代扣签约状态是否为成功,true代表成功
				//TODO...
			}
			if sub_channel  == "BC_CARD_CHARGE" {
				//message_detail中包含签约相关的详细信息，包括card_id
				//TODO...
			}
	}

	//业务逻辑处理完毕后，返回success。表示正确接收并处理了本次Webhook，其他返回都代表需继续重传本次的Webhook请求
	this.Print("success")
}

