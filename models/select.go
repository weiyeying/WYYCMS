package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func GetLinks(page,limit  int)([]orm.Params,[]orm.Params)  {
	var con []orm.Params
	var data []orm.Params
	orm:=orm.NewOrm()
	start:=page*limit
	orm.Raw("select * from cr_link limit ?,?",start,limit).Values(&data)
	orm.Raw("select count(id) as con from cr_link").Values(&con)
	return  data,con

}

func GetImg(page,limit  int)([]orm.Params,[]orm.Params)  {
	var con []orm.Params
	var data []orm.Params
	orm:=orm.NewOrm()
	start:=page*limit
	orm.Raw("select * from cr_img  limit ?,?",start,limit).Values(&data)
	orm.Raw("select count(id) as con from cr_img").Values(&con)
	return  data,con

}

func GetUser(page,limit  int)([]orm.Params,[]orm.Params)  {
	var con []orm.Params
	var data []orm.Params
	orm:=orm.NewOrm()
	start:=page*limit
	orm.Raw("select u.*,g.gname as gname from cr_user as u left join cr_group as g on u.group_id=g.id where u.status=1 limit ?,?",start,limit).Values(&data)
	orm.Raw("select count(id) as con from cr_user where status=1 ").Values(&con)
	return  data,con

}


func GetGroup(page,limit  int)([]orm.Params,[]orm.Params)  {
	var con []orm.Params
	var data []orm.Params
	orm:=orm.NewOrm()
	start:=page*limit
	orm.Raw("select * from cr_group   limit ?,?",start,limit).Values(&data)
	orm.Raw("select count(id) as con from cr_group  ").Values(&con)
	return  data,con

}

func GetAccess(page,limit  int)([]orm.Params,[]orm.Params)  {
	var con []orm.Params
	var data []orm.Params
	orm:=orm.NewOrm()
	start:=page*limit
	orm.Raw("select * from cr_access  where sata=1 order by rank desc limit ?,?",start,limit).Values(&data)
	orm.Raw("select count(id) as con from cr_access where sata=1 ").Values(&con)
	return  data,con

}
//id查询权限所有
func GetAccessAll()([]orm.Params)  {
	var data []orm.Params
	orm:=orm.NewOrm()
	orm.Raw("select * from cr_access  where sata=1 order by rank desc").Values(&data)
	return  data

}

//id查询权限所有 返回map[gid]=aid
func GetIdAccess(id int) map[string]string  {
	var data []orm.Params
	orm:=orm.NewOrm()
	maps:=make(map[string]string)
	orm.Raw("select * from cr_access_group  where aid=?",id).Values(&data)
	for _,v:=range data{
		k:=v["gid"].(string)
		vs:=v["aid"].(string)
		maps[k]=vs
	}
	return maps
}

//重组数据权限检测
func ChekcAccess(data []orm.Params,acclist map[string]string) []orm.Params {
	for _,v:=range data  {
		if _,ok:=acclist[v["id"].(string)];ok{
			v["is_check"]="1"
		}else{
			v["is_check"]="0"
		}
	}

	return data
}



//查询菜单
func GetAccessView()([]orm.Params)  {
	var data []orm.Params
	orm:=orm.NewOrm()
	orm.Raw("select * from cr_access  where sata=1 and type=1 ").Values(&data)
	return  data
}

//查询权限菜单
func GetAccessList(uid int)([]orm.Params)  {
	//orm.Debug=true
	var data []orm.Params
	orm:=orm.NewOrm()
	orm.Raw("select a.* from cr_access_group as c inner join cr_access as a on c.gid=a.id  where c.aid=? and a.sata=1 and a.type=1 order by a.rank desc ",uid).Values(&data)
	return  data
}

//查询是否有权限
func GetGroupAccess(access_id int,url string) bool {
	o:=orm.NewOrm()
	m:=CrAccess{}
	err:=o.Raw("select id from cr_access where con=? and type=2 and sata=1",url).QueryRow(&m)
	if err==nil{
	aid:=access_id
	gid:=m.Id
	ms:=CrAccessGroup{Gid:gid,Aid:aid}
	if err:=o.Read(&ms,"Gid","Aid");err==nil{
		return true
	}else{
		return false
	}
	}else{
		return false
	}

}

//查看新闻列表
func Getnews(page,limit  int,name,ntype string)([]orm.Params,[]orm.Params)  {
	orm.Debug=true
	var con []orm.Params
	var data []orm.Params
	orm:=orm.NewOrm()
	start:=page*limit
	if name!="" {
		name=" and ntitle like '%"+name+"%'"
	}
	if ntype!="" {
		ntype=" and ntype= %"+ntype+"%"
	}
	orm.Raw("select * from cr_news where nid>0 "+name+" "+ntype+" order by status desc,nid desc limit ?,?",start,limit).Values(&data)
	orm.Raw("select count(nid) as con from cr_news where nid>0 "+name+" "+ntype+" ").Values(&con)
	return  data,con

}