package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.test"
	c.Data["Email"] = "liida_1@163.com"
	c.TplNames = "index.tpl"
}
