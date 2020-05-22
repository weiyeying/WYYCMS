package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type CrImg struct {
	Id         int    `orm:"column(id);auto"`
	Title      string `orm:"column(title);size(100);null" description:"标题"`
	GroupName  string `orm:"column(group_name);size(255);null" description:"组名称"`
	Code       string `orm:"column(code);size(50);null" description:"分组编号"`
	Img        string `orm:"column(img);size(255);null" description:"图片"`
	Url        string `orm:"column(url);size(255);null" description:"地址"`
	StartTime  int    `orm:"column(start_time);null" description:"开始时间"`
	EndTime    int    `orm:"column(end_time);null" description:"结束时间"`
	CreateTime int    `orm:"column(create_time);null" description:"创建时间"`
}

func (t *CrImg) TableName() string {
	return "cr_img"
}

func init() {
	orm.RegisterModel(new(CrImg))
}

// AddCrImg insert a new CrImg into database and returns
// last inserted Id on success.
func AddCrImg(m *CrImg) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCrImgById retrieves CrImg by Id. Returns error if
// Id doesn't exist
func GetCrImgById(id int) (v *CrImg, err error) {
	o := orm.NewOrm()
	v = &CrImg{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCrImg retrieves all CrImg matches certain condition. Returns empty list if
// no records exist
func GetAllCrImg(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CrImg))
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

	var l []CrImg
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

// UpdateCrImg updates CrImg by Id and returns error if
// the record to be updated doesn't exist
func UpdateCrImgById(m *CrImg) (err error) {
	o := orm.NewOrm()
	v := CrImg{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCrImg deletes CrImg by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCrImg(id int) (err error) {
	o := orm.NewOrm()
	v := CrImg{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CrImg{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
