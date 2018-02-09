## BeeCloud Golang SDK (Open Source)

![license](https://img.shields.io/badge/license-MIT-brightgreen.svg) ![version](https://img.shields.io/badge/version-v1.0.0-blue.svg)

## 简介

本项目的官方GitHub地址是 [https://github.com/jasonhzy/beecloud-go](https://github.com/jasonhzy/beecloud-go)，本SDK 是基于 [BeeCloud RESTful API](https://github.com/beecloud/beecloud-rest-api)开发的Golang SDK，目前支持以下功能：
- 微信支付、支付宝支付、银联在线支付、百度钱包支付、京东支付、PayPal等多种支付方式
- 支付/退款订单总数的查询
- 支付订单和退款订单的查询
- 根据ID(支付/退款订单唯一标识)查询订单记录、退款记录

## 依赖

1、安装或者升级 Beego 和 Bee 的开发工具:

    $ go get -u github.com/astaxie/beego
    $ go get -u github.com/beego/bee

2、配置环境变量$GOPATH

    # 如果您还没添加 $GOPATH 变量
    $ echo 'export GOPATH="$HOME/go"' >> ~/.profile #或者 ~/.zshrc, 您所使用的sh对应的配置文件
    # 如果您已经添加了 $GOPATH 变量
    $ echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.profile #或者 ~/.zshrc, 您所使用的sh对应的配置文件

3、运行项目
    以/Users/jason/go($HOME:/Users/jason, $GOPATH:/Users/jason/go)为例， 安装Beego和Bee之后，会有bin、pkg、src三个目录。复制项目beecloud-go到/Users/jason/go/src下：

    $ cd /Users/jason/go/src/beecloud-go
    $ bee run beecloud-go # go run main.go

## 流程

下图为整个支付的流程:
![Flow](http://7xavqo.com1.z0.glb.clouddn.com/img-beecloud%20sdk.png)

步骤①：**（从网页服务器端）发送订单信息**
步骤②：**收到BeeCloud返回的渠道支付地址（比如支付宝的收银台）**
>*特别注意：
微信扫码返回的内容不是和支付宝一样的一个有二维码的页面，而是直接给出了微信的二维码的内容，需要用户自己将微信二维码输出成二维码的图片展示出来*

步骤③：**将支付地址展示给用户进行支付**
步骤④：**用户支付完成后通过一开始发送的订单信息中的return_url来返回商户页面**
>*特别注意：
微信没有自己的页面，二维码展示在商户自己的页面上，所以没有return url的概念，需要商户自行使用一些方法去完成这个支付完成后的跳转（比如后台轮询查支付结果）*

此时商户的网页需要做相应界面展示的更新（比如告诉用户"支付成功"或"支付失败")。**不允许**使用同步回调的结果来作为最终的支付结果，因为同步回调有极大的可能性出现丢失的情况（即用户支付完没有执行return url，直接关掉了网站等等），最终支付结果应以下面的异步回调为准。

步骤⑤：**（在商户后端服务端）处理异步回调结果（[Webhook](https://beecloud.cn/doc/?index=webhook)）**

付款完成之后，根据客户在BeeCloud后台的设置，BeeCloud会向客户服务端发送一个Webhook请求，里面包括了数字签名，订单号，订单金额等一系列信息。客户需要在服务端依据规则要验证**数字签名是否正确，购买的产品与订单金额是否匹配，这两个验证缺一不可**。验证结束后即可开始走支付完成后的逻辑。

## 初始化

1. 注册开发者: BeeCloud平台[注册账号](http://beecloud.cn/register/)
2. 创建应用: 使用注册的账号登陆,在控制台中创建应用,点击**"+添加应用"**创建新应用,具体可参考[快速开始](https://beecloud.cn/apply/)
3. 获取参数: 在新创建的应用中即可获取APP ID,APP Secret,Master Secret,Test Secret
4. 在代码中调用方法registerApp(请注意各个参数一一对应):

```
/* registerApp fun four params
 * @param(first) app_id beecloud平台的APP ID
 * @param(second) app_secret  beecloud平台的APP SECRET
 * @param(third) master_secret  beecloud平台的MASTER SECRET
 * @param(fouth) test_secret  beecloud平台的TEST SECRET, for sandbox
 */
import (
    "beecloud-go/sdk"
)

type PayController struct {
	sdk.ApiController
}

func (this *PayController) RegisterApp () {
	this.BaseController.RegisterApp("app_id", "app_secret", "master_secret", "test_secret");
	this.BaseController.SetSandbox(false);
	//this.BaseController.SetSandbox(true); //开启测试模式
}
```

5. LIVE模式和TEST模式

在代码中调用方法SetSandbox, 即:
- SetSandbox(false)或者不调用此方法, 即LIVE模式
- SetSandbox(true), 即TEST模式, 仅提供下单和支付订单查询的Sandbox模式

## 发起支付订单

    func (this *PayController) ToPay() {
        this.RegisterApp()
        bill := make(map[string]interface{})
        //bill["bill_no"] = "godemo1234567890"
        //total_fee(int 类型) 单位分
        bill["total_fee"] = 1
        //title UTF8编码格式，32个字节内，最长支持16个汉字
        bill["title"] = "xxxx"
        bill["channel"] = "yyyy"
        ...
        //发起支付方法 TODO...
    }

### 国际支付

国际支付目前主要是PayPal支付方式，主要提供了三种支付渠道类型：

- 当channel参数为PAYPAL_PAYPAL，即PayPal立即支付，接口返回url，用户跳转至此url，登陆paypal便可完成支付
- 当channel参数为PAYPAL_CREDITCARD，即信用卡支付，直接支付成功，接口返回的信用卡ID，此ID在快捷支付时需要
- 当channel参数为PAYPAL_SAVED_CREDITCARD，即存储的信用卡id支付，直接支付成功

支付调用的方法：

    result, err = this.International_bill(bill) //bill为支付所需参数

注：具体的请求参数和返回参数，请参考[国际支付REST API](https://github.com/beecloud/beecloud-rest-api/tree/master/international)

### 国内支付

国内支付适用于支付宝支付、京东支付、微信支付、百度钱包支付、银联在线支付等多种支付方式，选择不同的支付渠道类型，请求参数和返回参数也不尽相同，开发过程中要特别留意，严格按照提供的[线上支付REST API](https://github.com/beecloud/beecloud-rest-api/tree/master/online)文档进行开发。

支付调用方法：

    result, err = this.Bill(bill) //bill为支付所需参数

## 支付订单查询

用户提供相应的请求参数，并调用api中的bills方法，即可进行查询，主要包括单个（根据ID查询，即订单的唯一标识，而不是bill\_no）和多条记录的查询

多条记录的查询调用方法：

    func (this *PayController) BillQuery() {
        this.RegisterApp()
        conditions := make(map[string]interface{})
        conditions["channel"] = channel
        conditions["spay_result"] = true
        ...

        var res map[string]interface{}
        result, err := this.Bills(conditions) //根据条件查询的方法
    }

根据ID查询调用的方法：

    func (this *PayController) BillID() {
        this.RegisterApp()

        bill_id := "godemo1234567890"
        result, err := this.Bill_query_byid(bill_id) //根据ID查询的方法
    }

注：具体的请求参数和返回参数，请参考[线上支付REST API](https://github.com/beecloud/beecloud-rest-api/tree/master/online) **【5.订单查询】【11.支付订单查询(指定ID)】** 部分

## 订单总数查询

该接口主要用于对订单总数的统计，其中我们可以对其中的一段时间（即设置start\_time、end\_time）订单的统计，也可以只统计成功支付的订单（即设置spay\_result为true即可）

调用的方法：

    func (this *PayController) BillCount() {
        this.RegisterApp()
        conditions := make(map[string]interface{})
        conditions["channel"] = channel
        conditions["spay_result"] = true
        ...

        result1, err := this.Bills_count(conditions)
    }

注：具体的请求参数和返回参数，请参考[线上支付REST
API](https://github.com/beecloud/beecloud-rest-api/tree/master/online) **【6. 订单总数查询】** 部分

## 发起退款

退款接口仅支持对已经支付成功的订单进行退款，对于同一笔订单，仅能退款成功一次（同一个退款请求，第一次退款申请被驳回，可进行第二次退款申请）。
退款金额refund\_fee必须小于或者等于原始支付订单的total\_fee，如果是小于，则表示部分退款

退款接口包含预退款功能，当need\_approval值为true时，该接口开启预退款功能，预退款仅创建退款记录，并不真正发起退款，需后续调用审核接口，或者通过BeeCloud控制台的预退款界面，审核同意或者否决，才真正发起退款或者拒绝预退款。

退款调用方法：

    func (this PayController) ToRefund(){
        this.RegisterApp()
        data := make(map[string]interface{})
        data["bill_no"] = "godemo1234567890"
        data["refund_fee"] = 1
        ...
        this.Refund(data) //退款方法
    }

注：具体的请求参数和返回参数，请参考[线上支付REST
API](https://github.com/beecloud/beecloud-rest-api/tree/master/online) **【3. 退款】** 部分

## 退款订单查询

用户提供相应的请求参数，并调用api中的refunds方法，即可进行查询，主要包括单个（根据ID查询，即退款订单的唯一标识，而不是refund\_no）和多条记录的查询

多条记录查询调用的方法：

    func (this *PayController) RefundQuery() {
    	this.RegisterApp()

    	conditions := make(map[string]interface{})
    	conditions["channel"] = "xxxx"
    	conditions["limit"] = 10
    	...

    	result, err := this.Refunds(conditions)
    }

根据ID查询调用的方法：

    func (this *PayController) RefundQuery() {
        this.RegisterApp()

        refund_id := "godemo1234567890"
        result, err := this.Refund_query_byid(refund_id)
    }

注：具体的请求参数和返回参数，请参考[线上支付REST
API](https://github.com/beecloud/beecloud-rest-api/tree/master/online) **【7. 退款查询】【10.退款订单查询(指定ID)】** 部分


## 退款总数查询

该接口主要用于对退款订单总数的统计，其中我们可以对其中的一段时间（即设置start\_time、end\_time）退款订单的统计，也可以统计是否是预退款的退款订单（即设置need\_approval为true即可）

调用的方法：

    func (this *PayController) RefundQuery() {
        this.RegisterApp()

        conditions := make(map[string]interface{})
        conditions["channel"] = "xxxx"
        ...

        result, err := this.Refunds_count(conditions)
    }

注：具体的请求参数和返回参数，请参考[线上支付REST
API](https://github.com/beecloud/beecloud-rest-api/tree/master/online) **【8. 退款总数查询】** 部分

## 退款状态更新

退款状态更新接口提供查询退款状态以更新退款状态的功能，用于对退款后不发送回调的渠道（WX、YEE、KUAIQIAN、BD）退款后的状态更新。

调用方法：

    func (this PayController) RefundStatus() {
    	this.RegisterApp()

    	conditions := make(map[string]interface{})
    	conditions["refund_no"] = "godemo1234567890"
    	conditions["channel"] = "xxxx"

    	result, err = this.Refund_status(conditions)
    }

注：具体的请求参数和返回参数，请参考[线上支付REST
API](https://github.com/beecloud/beecloud-rest-api/tree/master/online) **【9. 退款状态更新】** 部分

## BeeCloud企业打款 - 打款到银行卡

调用方法：

    func (this TransferController) TransferTobank() {
        this.RegisterApp()

        data := make(map[string]interface{})
        data["bill_no"] = "gotransfer1234567890"
        data["bank_name"] = "中国银行" //银行全称, 不能写银行的缩写
        data["bank_account_no"] = "622269192199384xxxx" //收款方的银行卡号
        data["bank_account_name"] = "刘" //收款方的姓名或者单位名
        data["bank_code"] = "BOC" //银行的标准编码
        ...

        result, err := this.Gateway_transfer(data)
    }

注：具体的请求参数和返回参数，请参考[企业打款REST API](https://github.com/beecloud/beecloud-rest-api/tree/master/transfer) **【BeeCloud企业打款 - 打款到银行卡】** 部分

## 微信企业打款/微信红包

主要包括微信红包、微信企业打款

打款调用方法：

    func (this *) {
        this.RegisterApp()

        timestamp := this.GetTimestamp()
        data := make(map[string]interface{})
        data["timestamp"] = timestamp
        data["transfer_no"] = "gotransfer1234567890"
        data["channel"] = "xxxx"

        result, err = this.Transfer(data)
    }

注：具体的请求参数和返回参数，请参考[企业打款REST API](https://github.com/beecloud/beecloud-rest-api/tree/master/transfer) **【微信企业打款/微信红包】** 部分