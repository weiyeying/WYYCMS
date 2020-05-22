package main

import (
	"admin/controllers"
	"admin/models"
	_ "admin/routers"
	"github.com/astaxie/beego"

)

func init() {
	beego.ErrorController(&controllers.ErrorController{})
	//定义模板后缀
	beego.AddTemplateExt(".html")
	//数据库连接
	models.Init()
}

func main() {
	beego.Run()
}
