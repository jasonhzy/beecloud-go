package routers

import (
	"github.com/astaxie/beego"
	"beecloud-go/demo/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/bill", &controllers.PayController{}, "get:ToPay") //支付
	beego.Router("/wxpay/demo/pay", &controllers.PayController{}, "get:ToPay") //支付
	beego.Router("/bill/query", &controllers.PayController{}, "get:BillQuery") //订单查询
	beego.Router("/bill/id", &controllers.PayController{}, "get:BillID") //支付订单查询(指定id)
	beego.Router("/bill/status", &controllers.PayController{}, "get:BillStatus") //订单状态查询

	beego.Router("/refund", &controllers.PayController{}, "get:ToRefund") //退款
	beego.Router("/refund/query", &controllers.PayController{}, "get:RefundQuery") //退款查询
	beego.Router("/refund/id", &controllers.PayController{}, "get:RefundID") //退款订单查询(指定id)
	beego.Router("/refund/status", &controllers.PayController{}, "get:RefundStatus") //退款状态更新

	beego.Router("/pay/confirm", &controllers.PayController{}, "*:PayConfirm") //认证支付
	beego.Router("/card/charge", &controllers.PayController{}, "*:CardCharge") //签约支付

	beego.Router("/transfer", &controllers.TransferController{} ) //微信、支付宝、BC企业打款等
	beego.Router("/transfer/t1", &controllers.TransferController{}, "*:BcT1Transfer" ) //BC T1代付
	beego.Router("/transfer/gateway", &controllers.TransferController{}, "*:TransferTobank" ) //打款到银行卡

	beego.Router("/notify", &controllers.WebhookController{}) //webhook
}
