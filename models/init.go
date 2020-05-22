package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)


//数据库统一注册驱动
func Init() {
	mysql_user:=beego.AppConfig.String("mysql_user")
	mysql_password:=beego.AppConfig.String("mysql_password")
	mysql_host:=beego.AppConfig.String("mysql_host")
	mysql_dbname:=beego.AppConfig.String("mysql_dbname")
	dsn:=mysql_user+":"+mysql_password+"@tcp("+mysql_host+")/"+mysql_dbname+"?charset=utf8&loc=Local&parseTime=true"
	orm.RegisterDataBase("default","mysql" , dsn)
}
