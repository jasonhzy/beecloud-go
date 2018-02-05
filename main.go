package main

import (
	_ "beecloud-go/routers"
	"github.com/astaxie/beego"
	"time"
	"net/http"
	"text/template"
)

//时间转换
func convertT(num float64) (out string){
	tm := time.Unix(int64(num / 1000), 0)
	out = tm.Format("2006-01-02 03:04:05")
	return
}

func changeStatus(channel string) (out bool){
	switch channel {
		case "WX", "YEE", "KUAIQIAN", "BD":
			out = true
		default:
			out = false
	}
	return
}

func page_not_found(rw http.ResponseWriter, r *http.Request){
	t,_:= template.New("404.tpl").ParseFiles(beego.BConfig.WebConfig.ViewsPath+"/404.tpl")
	data :=make(map[string]interface{})
	data["content"] = "page not found"
	t.Execute(rw, data)
}

func main() {
	//404页面
	beego.ErrorHandler("404", page_not_found)

	beego.AddFuncMap("changeStatus", changeStatus)
	beego.AddFuncMap("convertT", convertT)
	beego.Run()
}



