package controllers

import (
	"admin/lib"
	"admin/models"
	"github.com/astaxie/beego"
	"time"
)

type LoginController struct {
	beego.Controller
}

//登录页面
func (e *LoginController) Login() {
	if e.GetSession("admin_id")!=nil{
		e.Ctx.Redirect(302,"/admin/index/index")
	}
	e.TplName = "login/login.html"
}

//登录接口
func (e *LoginController) CheckLogin() {
	name := e.GetString("name")
	pwd := e.GetString("pwd") + "weizheng"
	md5 := lib.Md5(pwd)
	user, err := models.GetIsUser(name, md5)
	code := make(map[string]interface{})
	if err != nil {
		code["code"] = -1
		code["msg"] = "err"
		code["data"] = "账号或密码错误"
		e.Data["json"] = &code
	} else {
		e.SetSession("admin_id",user.Id)
		e.SetSession("admin_name",user.Name)
		e.SetSession("access_id",user.GroupId)
		//写入记录
		cruserlog:=models.CrUserLoginLog{}
		cruserlog.CrId=user.Id
		cruserlog.Time=int(time.Now().Unix())
		cruserlog.Ip=e.Ctx.Input.IP()
		models.AddCrUserLoginLog(&cruserlog)
		code["code"] = 1000
		code["msg"] = "ok"
		code["data"] = "登录成功"
		e.Data["json"] = code
	}
	e.ServeJSON()
	e.StopRun()
}

//定义权限
func (e *LoginController) ErrorAccess(){
	e.TplName="err/401.html"
}

func (e *LoginController) LoginOut()  {
	e.DelSession("admin_id")
	e.DelSession("admin_name")
	e.Redirect("/admin/login/login",302)
}
