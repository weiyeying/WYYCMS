package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//用户获取
func GetIsUser(name, pwd string) (CrUser, error) {
	o := orm.NewOrm()
	var user CrUser
	err := o.Raw("select * from cr_user where name=? and password=? and status=1", name, pwd).QueryRow(&user)
	return user, err
}

//网站留言获取
func GetRelape(page int, limit int, rdata, rname string) ([]orm.Params, orm.Params) {
	var os1 []orm.Params
	var os2 []orm.Params
	pages := (page - 1) * limit
	orm := orm.NewOrm()
	where := " where rid>0 "
	if rname != "" {
		where += " and rname like '%" + rname + "%' "
	}
	if rdata != "" {
		where += " and rcon like '%" + rdata + "%'"
	}
	orm.Raw("select * from cr_rplace  "+where+" limit ?,? ", pages, limit).Values(&os1)
	orm.Raw("select count(rid) as con from cr_rplace  " + where + "").Values(&os2)
	return os1, os2[0]
}

//栏目查询

func GetGateory() []*CrCategory{
	var o []*CrCategory
	orm := orm.NewOrm()
	orm.Raw("select * from cr_category order by sorts desc").QueryRows(&o)
	tree := tree(o)
	return tree
}





func SortUserScore(fc []*CrCategory) []*CrCategory{
	for i:=0;i<len(fc)-1;i++{
		for j:= i+1 ;j<len(fc);j++{
			if fc[i].Sorts < fc[j].Sorts{
				fc[i],fc[j] = fc[j],fc[i]
			}
		}
	}
	return fc
}


func tree(list []*CrCategory)  []*CrCategory {
	data := buildData(list)
	result := makeTreeCore(0, data)
	//body, err := json.Marshal(result)
	//if err != nil {
	//	fmt.Println(err)
	//}
	return result
}

func buildData(list []*CrCategory) map[int]map[int]*CrCategory {
	var data map[int]map[int]*CrCategory = make(map[int]map[int]*CrCategory)
	for _, v := range list {
		id := v.Id
		fid := v.Pid
		if _, ok := data[fid]; !ok {
			data[fid] = make(map[int]*CrCategory)
		}
		data[fid][id] = v
	}
	return data
}

func makeTreeCore(index int, data map[int]map[int]*CrCategory) []*CrCategory {
	tmp := make([]*CrCategory, 0)
	for id, item := range data[index] {
		if data[id] != nil {
			item.Children = makeTreeCore(id, data)
		}
		tmp = append(tmp, item)
	}
	return tmp
}




