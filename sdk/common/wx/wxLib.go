package wx

import (
	"math/rand"
	"fmt"
	"sort"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"
	"io/ioutil"
	"errors"
	"bytes"
	"regexp"
)

type WxLibController struct{

}

//产生随机字符串，不长于32位
func (this WxLibController) CreateNoncestr( length int) (nonceStr string) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return
}

//格式化参数，签名过程需要使用
func (this WxLibController) FormatParams(paramsMap map[string]string) string{
	//STEP 1, 对key进行升序排序.
	var sorted_keys []string
	for k, _ := range paramsMap {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var paramsStr []string
	for _, k := range sorted_keys {
		v := fmt.Sprintf("%v", strings.TrimSpace(paramsMap[k]))
		if v != "" {
			paramsStr = append(paramsStr, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return strings.Join(paramsStr, "&")
}

//生成签名
func (this WxLibController) GetSign (paramsMap map[string]string, wxKey string ) string {
	//STEP 1：按字典序排序参数
	paramsStr := this.FormatParams(paramsMap)
	//STEP 2：在string后加入KEY
	signStr := paramsStr + "&key=" + wxKey
	//STEP 3：MD5加密
	sign := md5.New()
	sign.Write([]byte(signStr))
	//STEP 3：所有字符转为大写
	return strings.ToUpper(hex.EncodeToString(sign.Sum(nil)))
}

//xml结构
func (this WxLibController) ParamsToXml(data map[string]string) string {
	buf := bytes.NewBufferString("<xml>")
	for k, v := range data {
		str := "<![CDATA[%s]]>"
		flag , _ := regexp.MatchString("^\\d+\\.?\\d*$", v)
		if flag {
			str = "%s"
		}
		buf.WriteString(fmt.Sprintf("<%s>" + str + "</%s>", k, v, k))
	}
	buf.WriteString("</xml>")
	return buf.String()
}


func (this WxLibController) Error(strMsg string, err error) error {
	if err == nil {
		return errors.New(strMsg)
	}else {
		return errors.New(strMsg + ": " + err.Error())
	}
}

func (this WxLibController) http_post(url string, xml string) ([]byte, error){
	bc := &http.Client{
		Timeout: 30 * time.Second, //设置超时时间30s
	}
	res, err := bc.Post(url,"text/xml:charset=UTF-8", strings.NewReader(xml))
	if err != nil {
		return nil, this.Error("post", err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, this.Error("post result err", err)
	}
	return result, nil
}

func (this WxLibController) http_get(url string) ([]byte, error){
	bc := &http.Client{
		Timeout: 30 * time.Second, //设置超时时间30s
	}
	res, err := bc.Get(url)
	if err != nil {
		return nil, this.Error("get", err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, this.Error("get result err", err)
	}
	return result, nil
}