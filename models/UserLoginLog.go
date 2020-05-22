package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//获取所有数据分页获取
func GetAllData(pages,limit int) ([]orm.Params,orm.Params, error) {
	start_num := limit *(pages-1)
	var count []orm.Params
	var users []orm.Params
	orm := orm.NewOrm()
	 orm.Raw("select u.id,c.name,u.ip,u.time from cr_user_login_log as u left join cr_user as c on u.cr_id=c.id  order by u.id desc  limit  ?,?",start_num,limit).Values(&users)
	 orm.Raw("select count(id) as num from cr_user_login_log").Values(&count)
 	return users,count[0],nil
}

