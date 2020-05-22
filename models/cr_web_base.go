package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type CrWebBase struct {
	Id      int    `orm:"column(id);auto"`
	Title   string `orm:"column(title);size(255);null" description:"标题"`
	Key     string `orm:"column(key);size(255);null" description:"关键词"`
	Strkey  string `orm:"column(strkey);size(255);null" description:"描述"`
	Name    string `orm:"column(name);size(20);null" description:"联系人"`
	Email   string `orm:"column(email);size(50);null" description:"邮箱"`
	Phone   string `orm:"column(phone);size(20);null" description:"电话"`
	Img     string `orm:"column(img);size(255);null"`
	Banquan string `orm:"column(banquan);size(255);null"`
	Beian   string `orm:"column(beian);size(255);null"`
	Open    int8   `orm:"column(open);null" description:"网站开启 1开启 0关闭"`
	Deskey  string `orm:"column(deskey);size(255);null" description:"网站介绍"`
}

func (t *CrWebBase) TableName() string {
	return "cr_web_base"
}

func init() {
	orm.RegisterModel(new(CrWebBase))
}

// AddCrWebBase insert a new CrWebBase into database and returns
// last inserted Id on success.
func AddCrWebBase(m *CrWebBase) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCrWebBaseById retrieves CrWebBase by Id. Returns error if
// Id doesn't exist
func GetCrWebBaseById(id int) (v *CrWebBase, err error) {
	o := orm.NewOrm()
	v = &CrWebBase{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCrWebBase retrieves all CrWebBase matches certain condition. Returns empty list if
// no records exist
func GetAllCrWebBase(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CrWebBase))
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

	var l []CrWebBase
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

// UpdateCrWebBase updates CrWebBase by Id and returns error if
// the record to be updated doesn't exist
func UpdateCrWebBaseById(m *CrWebBase) (err error) {
	o := orm.NewOrm()
	v := CrWebBase{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCrWebBase deletes CrWebBase by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCrWebBase(id int) (err error) {
	o := orm.NewOrm()
	v := CrWebBase{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CrWebBase{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
