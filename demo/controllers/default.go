package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "https://beecloud.cn"
	c.Data["Email"] = "jasonhzy@beecloud.cn"
	c.TplName = "index.tpl"
}
