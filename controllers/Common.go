package controllers

import (
	"admin/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type Common struct {
	beego.Controller
}

func (e *Common) Prepare()  {
	if e.GetSession("admin_id")==nil{
		//e.Redirect("/admin/login/login",302)
	}
	aid:=e.GetSession("access_id") //权限组id
	aid=1
	url:=e.Ctx.Request.URL.Path
	ck:=models.GetGroupAccess(aid.(int),url)
	fmt.Println(url)
	fmt.Println(ck)
	if url!="/admin/index/index" && ck!=true{
		if e.IsAjax(){
			maps:=make(map[string]interface{})
			maps["code"]=-1
			maps["data"]="操作失败【没有权限】"
			maps["msg"]="操作失败【没有权限】"
			e.Data["json"]=maps
			e.ServeJSON()
			e.StopRun()
		}else{
			e.Redirect("/admin/login/erroraccess",302)
		}

		}
	e.Data["admin_name"]=e.GetSession("admin_name")
}



//返回json
func (e *Common) RtJson(data interface{},con interface{})  {
	maps:=make(map[string]interface{})
	maps["code"]=0
	maps["data"]=data
	maps["count"]=con
	e.Data["json"]=maps
	e.ServeJSON()
	e.StopRun()
}

//返回error json
func (e *Common) RtJsonString(code int,data string,err string)  {
	maps:=make(map[string]interface{})
	maps["code"]=code
	maps["data"]=data
	maps["msg"]=err
	e.Data["json"]=maps
	e.ServeJSON()
	e.StopRun()
}
//返回print json
func (e *Common) RtJsonWrite(code int,data ...interface{})  {
	maps:=make(map[string]interface{})
	maps["code"]=code
	if len(data)>0{
		maps["data0"]=data[0]
	}
	if len(data)>1{
		maps["data1"]=data[1]
	}
	e.Data["json"]=maps
	e.ServeJSON()
	e.StopRun()
}

//时间格式化
func FromTime( data []orm.Params,times string )  []orm.Params{
		if times==""{
			times="time"
		}
		for  k,v:=range data {
			ss:=v[times].(string)
			int64, _ := strconv.ParseInt(ss, 10, 64)
			data[k][times] = beego.Date(time.Unix(int64, 0), "Y-m-d H:i:s")
	}
		return data
}

//返回操作成功json
func (c *Common) RtAjaxSuccess(data interface{})  {
	maps:=make(map[string]interface{})
	maps["code"]=200
	maps["msg"]="操作成功"
	maps["data"]=data
	c.Data["json"]=maps
	c.ServeJSON()
	c.StopRun()
}



//验证ajax请求
func (c *Common) AjaxCheck()  {
	if !c.IsAjax(){
		maps:=make(map[string]interface{})
		maps["code"]=-1
		maps["err"]="请求错误"
		c.Data["json"]=maps
		c.ServeJSON()
		c.StopRun()
	}
}


