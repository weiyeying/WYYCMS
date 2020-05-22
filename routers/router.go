package routers

import (
	"admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
   ns:= beego.NewNamespace("admin",
   	   beego.NSAutoRouter(&controllers.LoginController{}), //自动路由
	   beego.NSAutoRouter(&controllers.IndexController{}), //自动路由
	   beego.NSAutoRouter(&controllers.AdminController{}), //自动路由
	   beego.NSAutoRouter(&controllers.NewsController{}), //自动路由
    )
   beego.AddNamespace(ns)
}
