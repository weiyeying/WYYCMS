package controllers

import (
	"admin/lib"
	"admin/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type AdminController struct {
	Common
}

//管理员列表-admin/userlist
func (c *AdminController) UserList() {
	page, err := c.GetInt("page")
	if err != nil {
		page = 0
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 10
	}
	data, con := models.GetUser(page, limit)
	FromTime(data, "addtime") //时间更新
	c.Data["data"] = &data
	c.Data["con"] = &con[0]
	c.Data["page"] = page

	c.TplName = "admin/user/userlist.html"
}

//删除管理员-admin/deluser
func (c *AdminController) DelUser() {
	if c.IsAjax() {
		if del, err := c.GetInt("id"); err == nil {
			o := orm.NewOrm()
			o.Raw("update cr_user set status=0 where id=?", del).Exec()
			c.RtJsonString(200, "操作成功", "")
		}
	}
	c.RtJsonString(-1, "操作失败", "")
}

//添加管理员-admin/adduser
func (c *AdminController) AddUser() {
	if c.IsAjax() {
		m := models.CrUser{}
		c.ParseForm(&m)
		m.Addtime = int(time.Now().Unix())
		m.Status = 1
		pwd := c.GetString("Password") + "weizheng"
		m.Password = lib.Md5(pwd)
		_, err := models.AddCrUser(&m)
		if err != nil {
			c.RtJsonString(-1, "添加失败", err.Error())
		} else {
			c.RtAjaxSuccess("添加成功")
		}
	}
	var m []orm.Params
	o := orm.NewOrm()
	o.Raw("select * from cr_group ").Values(&m)
	c.Data["data"] = m
	c.TplName = "admin/user/adduser.html"

}

//编辑管理员-admin/edituser
func (c *AdminController) EditUser() {
	o := orm.NewOrm()
	m := models.CrUser{}
	if c.IsAjax() {
		Yname := c.GetString("Yname")
		num := c.GetString("Num")
		GroupId := c.GetString("GroupId")
		ids := c.GetString("Id")
		_, err := o.Raw("update cr_user set num=?,yname=?,group_id=? where id=?", num, Yname, GroupId, ids).Exec()
		if err != nil {
			c.RtJsonString(-1, "编辑失败", err.Error())
		} else {
			c.RtAjaxSuccess("编辑成功")
		}
	}
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "参数不存在", err.Error())
	}
	err = o.Raw("select * from cr_user where id=?", id).QueryRow(&m)
	if err != nil {
		c.RtJsonString(-1, "参数不存在", err.Error())
	}
	var m2 []orm.Params
	o.Raw("select * from cr_group ").Values(&m2)
	c.Data["data2"] = &m2
	c.Data["data"] = &m
	c.Data["GroupId"] = strconv.Itoa(m.GroupId)
	c.TplName = "admin/user/edituser.html"

}

//角色列表-admin/grouplist
func (c *AdminController) GroupList() {
	del_status := 0
	page, err := c.GetInt("page")
	if err != nil {
		page = 0
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 10
	}
	_, err = c.GetInt("delid")
	if err == nil {
		//models.DeleteCrGroup(del)
		del_status = 1
	}

	data, con := models.GetGroup(page, limit)
	FromTime(data, "createtime") //时间更新
	c.Data["data"] = &data
	c.Data["con"] = &con[0]
	c.Data["del"] = del_status
	c.Data["page"] = page

	c.TplName = "admin/group/list.html"
}

//添加角色-admin/addgroup
func (c *AdminController) AddGroup() {
	if c.IsAjax() {
		m := models.CrGroup{}
		c.ParseForm(&m)
		m.Createtime = int(time.Now().Unix())
		_, err := models.AddCrGroup(&m)
		if err != nil {
			c.RtJsonString(-1, "添加失败", err.Error())
		} else {
			c.RtAjaxSuccess("添加成功")
		}
	}
	c.TplName = "admin/group/add.html"

}

//编辑角色-admin/editgroup
func (c *AdminController) EditGroup() {
	o := orm.NewOrm()
	m := models.CrGroup{}
	if c.IsAjax() {
		Yname := c.GetString("Gname")
		ids := c.GetString("Id")
		_, err := o.Raw("update cr_group set gname=? where id=?", Yname, ids).Exec()
		if err != nil {
			c.RtJsonString(-1, "编辑失败", err.Error())
		} else {
			c.RtAjaxSuccess("编辑成功")
		}
	}
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "参数不存在", err.Error())
	}
	err = o.Raw("select * from cr_group where id=?", id).QueryRow(&m)
	if err != nil {
		c.RtJsonString(-1, "参数不存在", err.Error())
	}
	c.Data["data"] = &m
	c.TplName = "admin/group/edit.html"

}

//设置角色权限-admin/setaccess
func (c *AdminController) SetAccess() {
	id, err := c.GetInt("gid");
	if err != nil {
		c.RtJsonString(-1, "参数错误", err.Error())
	}
	accesslist := models.GetAccessAll()
	access := models.GetIdAccess(id)
	data := models.ChekcAccess(accesslist, access)
	var os2 []orm.Params
	os := lib.Tree(data, 0, 1, &os2)
	c.Data["data"] = &os
	c.Data["aid"] = id
	c.TplName = "admin/group/setaccess.html"

}

//设置角色权限方法
func (c *AdminController) UpSetAccess() {
	if c.IsAjax() {
		id, err := c.GetInt("id");
		if err != nil {
			c.RtJsonString(-1, "参数错误", err.Error())
		}
		aid, err := c.GetInt("aid");
		if err != nil {
			c.RtJsonString(-1, "参数错误", err.Error())
		}
		status, err := c.GetInt("status");
		if err != nil {
			c.RtJsonString(-1, "参数错误", err.Error())
		}
		fmt.Println(id, status, aid)
		if status == 1 {
			m := models.CrAccessGroup{}
			m.Aid = aid
			m.Gid = id
			m.Createtime = int(time.Now().Unix())
			if _, err = models.AddCrAccessGroup(&m); err == nil {
				c.RtJsonString(200, "操作成功","")
			}
		}
		if status == 0 {
			o := orm.NewOrm()
			if _, err = o.Raw("delete from cr_access_group where aid=? and gid=?", aid, id).Exec(); err == nil {
				c.RtJsonString(200, "操作成功", "")
			}
		}
	}
	c.RtJsonString(-1, "参数错误", "")
}


//权限列表-admin/accesslist
func (c *AdminController) AccessList() {
	del_status := 0
	page, err := c.GetInt("page")
	if err != nil {
		page = 0
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 100000
	}
	del, err := c.GetInt("delid")
	if err == nil {
		o := orm.NewOrm()
		o.Raw("update cr_access set status=0 where id=?", del).Exec()
		del_status = 1
	}
	var os2 []orm.Params
	data, con := models.GetAccess(page, limit)
	os := lib.Tree(data, 0, 1, &os2)
	c.Data["data"] = &os
	c.Data["con"] = &con[0]
	c.Data["del"] = del_status
	c.Data["page"] = page

	c.TplName = "admin/access/list.html"
}

//添加权限-admin/addaccess
func (c *AdminController) AddAccess() {
	if c.IsAjax() {
		m := models.CrAccess{}
		c.ParseForm(&m)
		_, err := models.AddCrAccess(&m)
		if err != nil {
			c.RtJsonString(-1, "添加失败", err.Error())
		} else {
			c.RtAjaxSuccess("添加成功")
		}
	}
	var os2 []orm.Params
	data := models.GetAccessView()
	os := lib.Tree(data, 0, 1, &os2)
	c.Data["data"] = &os

	c.TplName = "admin/access/add.html"

}

//删除权限-admin/delaccess
func (c *AdminController) DelAccess() {
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "参数错误", err.Error())
	}
	o := orm.NewOrm()
	m := models.CrAccess{}
	err = o.Raw("select id from cr_access where pid=? and sata=1", id).QueryRow(&m)
	if err == nil {
		c.RtJsonString(-1, "请先删除子节点", "")
	}
	m.Id = id
	m.Sata = 0
	v := models.CrAccess{Id: id}
	if err = o.Read(&v); err == nil {
		if _, err = o.Update(&m, "sata"); err == nil {
			c.RtJsonString(200, "操作成功", "")
		}
	}
	c.RtJsonString(-1, "操作失败", "")
}

//权限排序-admin/soryaccess
func (c *AdminController) SortAccess() {
	orm.Debug = true
	id, err1 := c.GetInt("id")
	rank, err := c.GetInt("rank")
	if err != nil || err1 != nil {
		c.RtJsonString(-1, "参数错误", err.Error())
	}
	o := orm.NewOrm()
	m := models.CrAccess{}
	m.Id = id
	m.Rank = rank
	v := models.CrAccess{Id: id}
	if err = o.Read(&v); err == nil {
		if _, err = o.Update(&m, "rank"); err == nil {
			c.RtJsonString(200, "操作成功", "")
		}
	}
	c.RtJsonString(-1, "操作失败", "")
}

//编辑权限-admin/editaccess
func (c *AdminController) EditAccess() {
	o := orm.NewOrm()
	m := models.CrAccess{}
	if c.IsAjax() {
		id, err := c.GetInt("Id");
		if err != nil {
			c.RtJsonString(-1, "参数错误", err.Error())
		}
		v := models.CrAccess{Id: id}
		if err := o.Read(&v); err == nil {
			c.ParseForm(&m)
			if _, err = o.Update(&m, "pid", "accessname", "con", "type"); err == nil {
				c.RtJsonString(200, "操作成功", "")
			}
		}
		c.RtJsonString(-1, "编辑失败", err.Error())

	}
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "参数不存在", err.Error())
	}
	err = o.Raw("select * from cr_access where id=?", id).QueryRow(&m)
	if err != nil {
		c.RtJsonString(-1, "参数不存在", err.Error())
	}
	var os2 []orm.Params
	data := models.GetAccessView()
	os := lib.Tree(data, 0, 1, &os2)
	c.Data["data2"] = &m
	c.Data["data"] = &os
	c.Data["pid"] = strconv.Itoa(m.Pid)
	c.TplName = "admin/access/edit.html"

}