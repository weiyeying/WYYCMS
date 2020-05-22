package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type CrNews struct {
	Id          int    `orm:"column(nid);auto"`
	Ntitle      string `orm:"column(ntitle);size(255)"`
	Nkey        string `orm:"column(nkey);size(255);null" description:"seo关键字"`
	Nmiaosu     string `orm:"column(nmiaosu);size(255)" description:"描述"`
	Nclick      int    `orm:"column(nclick)"`
	Nyuanchuang int8   `orm:"column(nyuanchuang)" description:"0是技术文章    1是转载    2是日志"`
	Naddtime    int    `orm:"column(naddtime)" description:"添加时间"`
	Nimg        string `orm:"column(nimg);size(255);null" description:"文章图片"`
	Cid         int    `orm:"column(cid)" description:"分类"`
	Author      string `orm:"column(author);size(50);null" description:"作者"`
	Source      string `orm:"column(source);size(50);null" description:"来源"`
	Status      int8   `orm:"column(status);null" description:"状态1 上线  0下线"`
	Ntype       int8   `orm:"column(ntype);null" description:"类型 1技术文章 2散文"`
}

func (t *CrNews) TableName() string {
	return "cr_news"
}

func init() {
	orm.RegisterModel(new(CrNews))
}

// AddCrNews insert a new CrNews into database and returns
// last inserted Id on success.
func AddCrNews(m *CrNews) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCrNewsById retrieves CrNews by Id. Returns error if
// Id doesn't exist
func GetCrNewsById(id int) (v *CrNews, err error) {
	o := orm.NewOrm()
	v = &CrNews{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCrNews retrieves all CrNews matches certain condition. Returns empty list if
// no records exist
func GetAllCrNews(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CrNews))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []CrNews
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateCrNews updates CrNews by Id and returns error if
// the record to be updated doesn't exist
func UpdateCrNewsById(m *CrNews) (err error) {
	o := orm.NewOrm()
	v := CrNews{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCrNews deletes CrNews by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCrNews(id int) (err error) {
	o := orm.NewOrm()
	v := CrNews{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CrNews{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
