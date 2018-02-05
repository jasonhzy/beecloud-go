package sdk

const (
	//beecloud version
	api_version = "2"
	sdk_version = "go_1.0.0"

	//beecloud request url
	api_url = "https://api.beecloud.cn"

	//支付、支付订单查询(指定id)
	uri_bill = "rest/bill"
	//订单查询
	uri_bills = "rest/bills"
	//订单总数查询
	uri_bills_count = "rest/bills/count"

	uri_test_bill = "rest/sandbox/bill"
	uri_test_bills = "rest/sandbox/bills"
	uri_test_bills_count = "rest/sandbox/bills/count"

	//确认支付
	uri_pay_confirm = "rest/bill/confirm"

	//退款预退款批量审核退款订单查询(指定id)
	uri_refund = "rest/refund"
	//退款查询
	uri_refunds = "rest/refunds"
	//退款总数查询
	uri_refunds_count = "rest/refunds/count"
	//退款状态更新
	uri_refund_status = "rest/refund/status"

	//获取银行列表
	uri_bc_gateway_banks = "rest/bc_gateway/banks"

	//单笔打款 - 支付宝/微信
	uri_transfer = "rest/transfer"
	//批量打款 - 支付宝
	uri_transfers = "rest/transfers"
	//bc企业打款 - 支持银行
	uri_bc_transfer_banks = "rest/bc_transfer/banks"
	//代付 - 银行卡
	uri_bc_transfer = "rest/bc_transfer"
	//畅捷代付
	uri_cj_transfer = "rest/cj_transfer"
	//京东代付
	uri_jd_transfer = "rest/bc_user_transfer"
	//beepay自动打款 - 打款到银行卡
	uri_gateway_transfer = "rest/gateway/bc_transfer"

	//线下支付-撤销订单
	uri_offline_bill = "rest/offline/bill"
	//线下订单状态查询
	uri_offline_bill_status = "rest/offline/bill/status"
	//线下退款
	uri_offline_refund = "rest/offline/refund"

	//international
	uri_international_bill = "rest/international/bill"
	uri_international_refund = "rest/international/refund"

	//发送验证码
	uri_sms = "sms"

	//代扣api
	uri_card_charge_sign = "sign"

	//t1代付银行列表接口
	uri_t1_express_transfer_banks = "rest/t1express/transfer/banks"
	//代付接口
	uri_t1_express_transfer = "rest/t1express/transfer"

	//auth
	uri_auth = "auth"

	//subscription
	uri_subscription = "subscription"
	uri_subscription_plan = "plan"
	uri_subscription_banks = "subscription_banks"

	//user system
	//单个用户注册接口
	uri_usersys_user = "rest/user"
	//批量用户导入接口／查询接口
	uri_usersys_multi_users = "rest/users"
	//历史数据补全接口（批量）
	uri_usersys_history_bills = "rest/history_bills"

	//coupon
	//发放卡券, 优惠券根据id或其他条件查询
	uri_coupon = "rest/coupon"
	//根据优惠券模板id或其他条件查询
	uri_coupon_temp = "rest/coupon/template"
)








