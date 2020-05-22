package controllers

import (
	"github.com/astaxie/beego"
)


type ErrorController struct {
	beego.Controller
}
//定义错误404
func (e *ErrorController) Error404()  {
		e.TplName="err/404.html"
}



func main() {
	
}
