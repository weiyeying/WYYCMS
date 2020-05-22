package models

import "github.com/astaxie/beego/orm"

//文章分类查询

func GetNewsGateory() []*CrNewscategory{
	var o []*CrNewscategory
	orm := orm.NewOrm()
	orm.Raw("select * from cr_newscategory order by sorts desc").QueryRows(&o)
	tree := Newstree(o)
	return tree
}


func NewsSortUserScore(fc []*CrNewscategory) []*CrNewscategory{
	for i:=0;i<len(fc)-1;i++{
		for j:= i+1 ;j<len(fc);j++{
			if fc[i].Sorts < fc[j].Sorts{
				fc[i],fc[j] = fc[j],fc[i]
			}
		}
	}
	return fc
}


func Newstree(list []*CrNewscategory)  []*CrNewscategory {
	data := NewsbuildData(list)
	result := NewsmakeTreeCore(0, data)
	//body, err := json.Marshal(result)
	//if err != nil {
	//	fmt.Println(err)
	//}
	return result
}

func NewsbuildData(list []*CrNewscategory) map[int]map[int]*CrNewscategory {
	var data map[int]map[int]*CrNewscategory = make(map[int]map[int]*CrNewscategory)
	for _, v := range list {
		id := v.Id
		fid := v.Pid
		if _, ok := data[fid]; !ok {
			data[fid] = make(map[int]*CrNewscategory)
		}
		data[fid][id] = v
	}
	return data
}

func NewsmakeTreeCore(index int, data map[int]map[int]*CrNewscategory) []*CrNewscategory {
	tmp := make([]*CrNewscategory, 0)
	for id, item := range data[index] {
		if data[id] != nil {
			item.Children = NewsmakeTreeCore(id, data)
		}
		tmp = append(tmp, item)
	}
	return tmp
}