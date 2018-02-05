package sdk

import (
	//"github.com/astaxie/beego"
	"bytes"
	"net/http"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"time"
	"strconv"
	"encoding/json"
	"net/url"
	"errors"
)

type BaseController struct {
	//beego.Controller
	app_id string
	app_secret string
	master_secret string
	test_secret string
	flag bool
}

func (this *BaseController) RegisterApp(app_id string, app_secret string, master_secret string, test_secret string){
	this.app_id = app_id
	this.app_secret = app_secret
	this.master_secret = master_secret
	this.test_secret = test_secret
}

func (this *BaseController) SetSandbox(flag bool){
	this.flag = flag
}

func (this *BaseController) GetSandbox() bool{
	return this.flag
}

func (this BaseController) GetTimestamp() int64{
	return time.Now().UnixNano() / 1000000 //毫秒
}

func (this BaseController) In_array(strMsg string, list []string) bool {
	for _, v  := range list  {
		if strMsg == v {
			return true
		}
	}
	return false
}

func (this BaseController) GetApiUrl(url string) string {
	//return beego.AppConfig.String("api_url") + "/" + beego.AppConfig.String("api_version") + "/" + url
	return api_url + "/" + api_version + "/" + url
}

func (this BaseController) GetSign(app_id string, timestamp int64, secret string) string{
	signStr := app_id + strconv.FormatInt(timestamp, 10)+ secret

	sign := md5.New()
	sign.Write([]byte(signStr))
	return hex.EncodeToString(sign.Sum(nil))
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func (this BaseController) Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}
	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}

func (this BaseController) Error(strMsg string, err error) error {
	return errors.New(strMsg + err.Error())
}

//go sdk version
func (this BaseController) GetSdkVersion (data map[string]interface{}) {
	version := make(map[string]string)
	version["sdk_version"] = sdk_version //beego.AppConfig.String("sdk_version")
	data["bc_analysis"] = version
}

func (this BaseController) GetCommonParams(data map[string]interface{}, secret_type int){
	var secret string
	switch secret_type {
		case 0:
			secret = this.app_secret
		case 1:
			secret = this.master_secret
		case 2:
			secret = this.test_secret
	}
	data["app_id"] = this.app_id
	if data["timestamp"] == nil {
		data["timestamp"] = this.GetTimestamp()
	}
	data["app_sign"] = this.GetSign(this.app_id, data["timestamp"].(int64), secret)
}

func (this BaseController) http_post(request_url string, data map[string]interface{}) ([]byte, error){
	jsonObj, err := json.Marshal(data)
	if err != nil {
		return nil, this.Error("json err: ", err)
	}
	info := bytes.NewBuffer([]byte(jsonObj))
	url := this.GetApiUrl(request_url)
	bc := &http.Client{
		Timeout: 30 * time.Second, //设置超时时间30s
	}
	res, err := bc.Post(url,"application/json;charset=utf-8", info)
	if err != nil {
		return nil, this.Error("post: ", err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, this.Error("post result err: ", err)
	}
	return result, nil
}

func (this BaseController) http_get(request_url string, data map[string]interface{}) ([]byte, error){
	jsonObj, err := json.Marshal(data)
	if err != nil {
		return nil, this.Error("json err: ", err)
	}

	u := url.Values{}
	u.Set("para", string(jsonObj))
	url_prefix := this.GetApiUrl(request_url)
	url := url_prefix + "?" + u.Encode()

	bc := &http.Client{
		Timeout: 30 * time.Second, //设置超时时间30s
	}
	res, err := bc.Get(url)
	if err != nil {
		return nil, this.Error("get: ", err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, this.Error("get result err: ", err)
	}
	return result, nil
}





