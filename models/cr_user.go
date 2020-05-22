package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type CrUser struct {
	Id       int    `orm:"column(id);auto"`
	Name     string `orm:"column(name);size(11);null" description:"后台用户名"`
	Password string `orm:"column(password);size(255);null" description:"密码"`
	GroupId  int    `orm:"column(group_id);null" description:"组id"`
	Num      string `orm:"column(num);size(200);null" description:"员工编号"`
	Yname    string `orm:"column(yname);size(200);null" description:"员工姓名"`
	Addtime  int    `orm:"column(addtime);null" description:"入职时间"`
	Bre      int    `orm:"column(bre);null" description:"入职时间"`
	Status   int8   `orm:"column(status);null" description:"状态1 启用 0删除"`
}

func (t *CrUser) TableName() string {
	return "cr_user"
}

func init() {
	orm.RegisterModel(new(CrUser))
}

// AddCrUser insert a new CrUser into database and returns
// last inserted Id on success.
func AddCrUser(m *CrUser) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCrUserById retrieves CrUser by Id. Returns error if
// Id doesn't exist
func GetCrUserById(id int) (v *CrUser, err error) {
	o := orm.NewOrm()
	v = &CrUser{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCrUser retrieves all CrUser matches certain condition. Returns empty list if
// no records exist
func GetAllCrUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CrUser))
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

	var l []CrUser
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

// UpdateCrUser updates CrUser by Id and returns error if
// the record to be updated doesn't exist
func UpdateCrUserById(m *CrUser) (err error) {
	o := orm.NewOrm()
	v := CrUser{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCrUser deletes CrUser by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCrUser(id int) (err error) {
	o := orm.NewOrm()
	v := CrUser{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CrUser{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
