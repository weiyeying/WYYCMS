package controllers

import (
	"admin/lib"
	"admin/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type IndexController struct {
	Common
}

//首页
func (c *IndexController) Index() {
	uid:=c.GetSession("access_id")
	if uid==nil{
		//c.Redirect("/admin/login/login",302)
	}
	uid=1
	data:=models.GetAccessList(uid.(int))
	var os []orm.Params
	os2:=lib.Tree(data,0,uid.(int),&os) //分类排序
	os3:=lib.TreeClass(*os2)
	c.Data["data"]=&os3
	c.TplName = "layu/main.html"
}

//首页
func (c *IndexController) Indexs() {
	c.TplName = "index/index.html"
}

//登录记录
func (c *IndexController) UserInfo() {
	c.TplName = "index/userinfo.html"
}

//登录记录
func (c *IndexController) UserInfoData() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	data, con, _ := models.GetAllData(page, limit)
	for k, vs := range data {
		ss := vs["time"].(string)
		int64, _ := strconv.ParseInt(ss, 10, 64)
		data[k]["time"] = beego.Date(time.Unix(int64, 0), "Y-m-d H:i:s")

	}
	c.RtJson(data, con["num"])
}

//基本配置-admin/webbase
func (c *IndexController) Webbase() {
	data, err := models.GetCrWebBaseById(1)
	if err != nil {
		c.RtJsonString(0, "", err.Error())
	}
	c.Data["data"] = &data
	c.TplName = "index/base.html"
}

//基本配置修改-admin/basepost
func (c *IndexController) BasePost() {
	m := models.CrWebBase{}
	c.ParseForm(&m)
	if c.GetString("Open") == "" {
		m.Open = 0
	} else {
		m.Open = 1
	}
	err := models.UpdateCrWebBaseById(&m)
	if err != nil {
		c.RtJsonString(0, "", err.Error())
	}
	c.RtJsonString(1, "", "操作成功")

}

//上传图片
func (c *IndexController) UploadFile() {
	f, n, err := c.GetFile("file")
	if err != nil {
		c.RtJsonString(-1, "", err.Error())

	}
	//关闭文件句柄
	defer f.Close()
	//获取文件路径
	t := time.Now()
	//文件目录创建
	path := "./static/data/" + t.Format("20060102")
	err = os.MkdirAll(path, 0777)
	if err != nil {
		c.RtJsonString(-1, "", err.Error())
	}

	//文件名拼接
	ext := filepath.Ext(n.Filename)
	rand.Seed(time.Now().UnixNano())
	name := strconv.Itoa(rand.Intn((9000)+1000)) + strconv.Itoa(int(t.Unix())) + ext
	filename := fmt.Sprintf("%s/%s", path, name)
	//保存文件
	err = c.SaveToFile("file", filename)
	if err != nil {
		c.RtJsonString(-1, "", err.Error())
	}
	//返回数据
	rename := t.Format("20060102")
	filename = fmt.Sprintf("%s/%s", rename, name)
	c.RtJsonString(0, filename, "")

}

//网站友情链接管理-admin/weblinks
func (c *IndexController) WebLinks() {
	page, err := c.GetInt("page")
	if err != nil {
		page = 0
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 10
	}
	data, con := models.GetLinks(page, limit)
	FromTime(data, "time") //时间更新
	c.Data["data"] = &data
	c.Data["con"] = &con[0]
	c.Data["page"] = page
	c.TplName = "index/links/link.html"

}

//友情链接删除-admin/dellink
func (c *IndexController) DelLink()  {
	if c.IsAjax(){
		if del, err := c.GetInt("id");err==nil{
			models.DeleteCrLink(del)
			c.RtJsonString(200, "删除成功", "")
		}
	}
	c.RtJsonString(-1, "操作失败", "操作失败")
}

//友情链接状态修改-admin/uplinkstatus
func (c *IndexController) UpLinkeStatus() {
	id, err := c.GetInt("id")
	status, err2 := c.GetInt("status")
	if err != nil || err2 != nil {
		c.RtJsonString(0, "", "修改失败")
	}

	o := orm.NewOrm()
	v := models.CrLink{Id: id}
	if err = o.Read(&v); err == nil {
		v.Status = int8(status)
		_, err := o.Update(&v, "status")
		if err != nil {
			c.RtJsonString(0, "", "修改失败")
		} else {
			c.RtJsonString(1, "", "修改成功")
		}
	} else {
		c.RtJsonString(0, "", "修改失败")
	}

}

//友情链接添加-admin/addweblink
func (c *IndexController) AddWebLink() {
	name := c.GetString("Name")
	if name != "" {
		m := models.CrLink{}
		c.ParseForm(&m)
		m.Time = int(time.Now().Unix())
		m.Status = 1
		_, err := models.AddCrLink(&m)
		if err != nil {
			c.RtJsonString(0, "", err.Error())
		} else {
			c.RtJsonString(1, "", "添加成功")
		}
	}
	c.TplName = "index/links/addlink.html"
}

//友情链接编辑-admin/editweblink
func (c *IndexController) EditWebLink() {
	name := c.GetString("Name")
	id, err := c.GetInt("id")
	if name != "" && err != nil {
		m := models.CrLink{}
		c.ParseForm(&m)
		m.Time = int(time.Now().Unix())
		err := models.UpdateCrLinkById(&m)
		if err != nil {
			c.RtJsonString(0, "", err.Error())
		} else {
			c.RtJsonString(1, "", "编辑成功")
		}
	}
	data, err := models.GetCrLinkById(id)
	c.Data["data"] = &data
	c.TplName = "index/links/editlink.html"
}

//网站留言-index/webreplace
func (c *IndexController) WebReplace() {
	c.TplName = "index/replace.html"
}

//网站留言获取
func (c *IndexController) ReplaceData() {
	page, err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 10
	}
	rname := c.GetString("rnames")
	rdata := c.GetString("rdata")
	os1, os2 := models.GetRelape(page, limit, rdata, rname)
	os1 = FromTime(os1, "rtime")
	c.RtJson(os1, os2["con"])

}

//导航管理-index/gategory
func (c *IndexController) Gategory() {
	o := models.GetGateory()
	sortcat := models.SortUserScore(o)
	c.Data["list"] = sortcat
	c.TplName = "index/category/category.html"
}

//导航修改排序-index/upcatrgorysort
func (c *IndexController) UpCatrgorySort() {
	c.AjaxCheck() //检测请求
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "id参数错误", err.Error())
	}
	val, err := c.GetInt("val")
	if err != nil {
		c.RtJsonString(-1, "val参数错误", err.Error())
	}
	o := orm.NewOrm()
	m := models.CrCategory{Id: id}
	if err = o.Read(&m); err == nil {
		m.Sorts = val
		num, _ := o.Update(&m, "sorts")
		c.RtAjaxSuccess("更新数量" + strconv.Itoa(int(num)))
	} else {
		c.RtJsonString(-1, "操作失败", err.Error())
	}

}

//导航删除-index/delcatrgory
func (c *IndexController) DelCatrgory() {
	c.AjaxCheck()
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "id参数错误", err.Error())
	}
	o := orm.NewOrm()
	m := models.CrCategory{Id: id}
	//验证是否有子
	var s []orm.Params
	cnum, _ := o.Raw("select id from cr_category where pid=? and pid!=0 ", id).Values(&s)
	if cnum > 0 {
		c.RtJsonString(-1, "删除失败请先删除子节点", "")
	}
	if err = o.Read(&m); err == nil {
		var num int64
		if num, err = o.Delete(&models.CrCategory{Id: id}); err == nil {
			c.RtAjaxSuccess("删除数量" + strconv.Itoa(int(num)))
		} else {
			c.RtJsonString(-1, "删除失败", err.Error())
		}
	} else {
		c.RtJsonString(-1, "id删除失败", err.Error())
	}
}

//添加导航-index/adcategory
func (c *IndexController) AddCategory() {
	o := orm.NewOrm()
	if c.IsAjax() {
		s := models.CrCategory{}
		c.ParseForm(&s)
		_, err := o.Insert(&s)
		if err != nil {
			c.RtJsonString(-1, "添加失败", err.Error())

		} else {
			c.RtAjaxSuccess("添加成功")
		}

	}
	var m []orm.Params
	o.Raw("select id,category from cr_category where pid=0 order by sorts desc").Values(&m)
	c.Data["data"] = m
	c.TplName = "index/category/addcategory.html"

}

//编辑导航-index/editcategory
func (c *IndexController) EditCategory() {
	o := orm.NewOrm()
	s := models.CrCategory{}

	if c.IsAjax() {
		c.ParseForm(&s)
		err := models.UpdateCrCategoryById(&s)
		if err != nil {
			c.RtJsonString(-1, "更新失败", err.Error())

		} else {
			c.RtAjaxSuccess("添加成功")
		}

	}
	var m []orm.Params
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "编辑失败", err.Error())
	}
	o.Raw("select id,category from cr_category where pid=0 order by sorts desc").Values(&m)
	o.Raw("select * from cr_category where id=?", id).QueryRow(&s)
	c.Data["data"] = m  //主导航
	c.Data["data2"] = s //导航数据
	c.Data["pid"] = strconv.Itoa(s.Pid)
	c.TplName = "index/category/editcategory.html"

}

//图片管理列表-index/imglist
func (c *IndexController) ImgList() {
	page, err := c.GetInt("page")
	if err != nil {
		page = 0
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 10
	}

	data, con := models.GetImg(page, limit)
	FromTime(data, "start_time") //时间更新
	FromTime(data, "end_time") //时间更新
	c.Data["data"] = &data
	c.Data["con"] = &con[0]
	c.Data["page"] = page

	c.TplName = "index/img/imglist.html"
}

//图片删除-index/delimg
func (c *IndexController) DelImg(){
	if c.IsAjax(){
		if del, err := c.GetInt("id");err==nil{
			models.DeleteCrImg(del)
			c.RtJsonString(200, "删除成功", "")
		}
	}
	c.RtJsonString(-1, "操作失败", "操作失败")
}

//图片添加-index/addimg
func (c *IndexController) AddImg() {
	if c.IsAjax() {
		m := models.CrImg{}
		c.ParseForm(&m)
		starttime := lib.String_Unix(c.GetString("Start_time"))
		endtime := lib.String_Unix(c.GetString("End_time"))
		m.StartTime = starttime
		m.EndTime = endtime
		_, err := models.AddCrImg(&m)
		if err != nil {
			c.RtJsonString(-1, "添加失败", err.Error())
		} else {
			c.RtAjaxSuccess("添加成功")
		}
	}

	c.TplName = "index/img/addimg.html"

}

//图片编辑-index/editimg
func (c *IndexController) EditImg() {
	if c.IsAjax() {
		m := models.CrImg{}
		c.ParseForm(&m)
		starttime := lib.String_Unix(c.GetString("Start_time"))
		endtime := lib.String_Unix(c.GetString("End_time"))
		m.StartTime = starttime
		m.EndTime = endtime
		err := models.UpdateCrImgById(&m)
		if err != nil {
			c.RtJsonString(-1, "编辑失败", err.Error())
		} else {
			c.RtAjaxSuccess("编辑成功")
		}
	}
	id,err:=c.GetInt("id")
	if err!=nil{
		c.RtJsonString(-1,"参数不存在",err.Error())
	}
	o:=orm.NewOrm()
	m:=models.CrImg{}
	err=o.Raw("select * from cr_img where id=?",id).QueryRow(&m)
	if err!=nil{
		c.RtJsonString(-1,"参数不存在",err.Error())
	}

	c.Data["data"]=&m
	c.Data["starttime"]=lib.Unix_time(m.StartTime)
	c.Data["endtime"]=lib.Unix_time(m.EndTime)
	c.TplName = "index/img/editimg.html"

}
