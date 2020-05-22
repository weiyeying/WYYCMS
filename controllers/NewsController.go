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

type NewsController struct {
	Common
}

//文章列表
func (c *NewsController) NewsList() {
	page, err := c.GetInt("page")
	if err != nil {
		page = 0
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 8
	}
	name:=c.GetString("name")
	ntype:=c.GetString("ntype")
fmt.Println(name)
	data, con := models.Getnews(page, limit,name,ntype)
	FromTime(data, "naddtime") //时间更新
	c.Data["data"] = &data
	c.Data["con"] = &con[0]
	c.Data["page"] = page
	c.Data["name"] = name
	c.Data["ntype"] = ntype

	c.TplName = "news/news/newslist.html"
}

//添加文章
func (c *NewsController) AddNews() {
	if c.IsAjax() {
		o := models.CrNews{}
		o.Naddtime = int(time.Now().Unix())
		c.ParseForm(&o)
		id, err := models.AddCrNews(&o)
		if err == nil {
			o2 := models.CrNewsData{}
			o2.Nid = int(id)
			o2.Body = c.GetString("bz")
			_, err = models.AddCrNewsData(&o2)
			if err == nil {
				c.RtJsonString(1, "添加成功", "添加成功")
			} else {
				c.RtJsonString(-1, "添加失败", err.Error())
			}
		} else {
			c.RtJsonString(-1, "添加失败", err.Error())
		}
	}
	var m []orm.Params
	var d []orm.Params
	orm := orm.NewOrm()
	orm.Raw("select * from cr_newscategory order by sorts desc").Values(&m)
	data := lib.Tree(m, 0, 1, &d)

	c.Data["data"] = &data
	c.TplName = "news/news/addnews.html"
}

//删除文章
func (c *NewsController) DelNews() {
	if c.IsAjax() {
		id, err := c.GetInt("id")
		if err != nil {
			c.RtJsonString(-1, "删除失败", err.Error())
		}
		orm := orm.NewOrm()
		if _, err = orm.Delete(&models.CrNews{Id: id}); err != nil {
			c.RtJsonString(-1, "删除失败", err.Error())
		}
		if _, err = orm.Delete(&models.CrNewsData{Nid: id},"Nid"); err != nil {
			c.RtJsonString(-1, "删除失败", err.Error())
		}
		c.RtJsonString(1, "删除成功", "")

	}
}

//更新文章
func (c *NewsController) UpdateNews() {
	if c.IsAjax() {
		or := orm.NewOrm()
		or.Begin()
		newsid := c.GetString("Id")
		newsdataid := c.GetString("Bid")

		newsid2, _ := strconv.Atoi(newsid)
		o := models.CrNews{Id: newsid2}
		o.Naddtime = int(time.Now().Unix())
		c.ParseForm(&o)
		if c.GetString("open") == "on" {
			o.Status = 1
		} else {
			o.Status = 0
		}
		_, err := or.Update(&o, "ntitle", "nkey", "nmiaosu", "nyuanchuang", "naddtime", "nimg", "author", "source", "status","ntype")

		newsdataid2, _ := strconv.Atoi(newsdataid)
		o2 := models.CrNewsData{Id: newsdataid2}
		o2.Body = c.GetString("bz")
		var num int64
		num, err2 := or.Update(&o2, "body")
		fmt.Println(err)
		fmt.Println(err2)
		fmt.Println(num)
		if err != nil || err2 != nil {
			or.Rollback()
			c.RtJsonString(-1, "更新失败", "")
		} else {
			or.Commit()
			c.RtJsonString(1, "更新成功", "更新成功")
		}

	}
	nid, err := c.GetInt("nid")
	if err != nil {
		c.RtJsonString(-1, "更新失败", err.Error())
	}
	var m []orm.Params
	var d []orm.Params
	u := models.CrNews{}
	uc := models.CrNewsData{}
	orm := orm.NewOrm()
	orm.Raw("select * from cr_newscategory order by sorts desc").Values(&m)
	data := lib.Tree(m, 0, 1, &d)
	orm.Raw("select * from cr_news where nid=? ", nid).QueryRow(&u)
	orm.Raw("select * from cr_news_data where nid=?", nid).QueryRow(&uc)
	c.Data["data"] = &data
	c.Data["u"] = &u
	c.Data["uc"] = &uc
	c.Data["pid"] = strconv.Itoa(u.Cid)
	c.Data["nyuanchuang"] = u.Nyuanchuang
	c.Data["ntype"] = u.Ntype
	fmt.Println(uc)
	c.TplName = "news/news/upnews.html"
}

//上传图片
func (c *NewsController) NewsUploadFile() {
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
	maps := make(map[string]interface{})
	data := make(map[string]interface{})
	data["src"] = beego.AppConfig.String("imgpath") + filename
	data["title"] = ""
	maps["code"] = 0
	maps["data"] = data
	maps["msg"] = "上传成功"
	c.Data["json"] = maps
	c.ServeJSON()
	c.StopRun()

}

//文章分类管理/admin/news/newsgategory
func (c *NewsController) NewsGategory() {
	o := models.GetNewsGateory()
	sortcat := models.NewsSortUserScore(o)
	c.Data["list"] = sortcat
	c.TplName = "news/category/category.html"
}

//文章分类管理排序-/upcatrgorysort
func (c *NewsController) NewsUpCatrgorySort() {
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
	m := models.CrNewscategory{Id: id}
	if err = o.Read(&m); err == nil {
		m.Sorts = val
		num, _ := o.Update(&m, "sorts")
		c.RtAjaxSuccess("更新数量" + strconv.Itoa(int(num)))
	} else {
		c.RtJsonString(-1, "操作失败", err.Error())
	}

}

//文章分类删除
func (c *NewsController) NewsDelCatrgory() {
	c.AjaxCheck()
	id, err := c.GetInt("id")
	if err != nil {
		c.RtJsonString(-1, "id参数错误", err.Error())
	}
	o := orm.NewOrm()
	m := models.CrNewscategory{Id: id}
	//验证是否有子
	var s []orm.Params
	cnum, _ := o.Raw("select id from cr_newscategory where pid=? and pid!=0 ", id).Values(&s)
	if cnum > 0 {
		c.RtJsonString(-1, "删除失败请先删除子节点", "")
	}
	if err = o.Read(&m); err == nil {
		var num int64
		if num, err = o.Delete(&models.CrNewscategory{Id: id}); err == nil {
			c.RtAjaxSuccess("删除数量" + strconv.Itoa(int(num)))
		} else {
			c.RtJsonString(-1, "删除失败", err.Error())
		}
	} else {
		c.RtJsonString(-1, "id删除失败", err.Error())
	}
}

//文章分类添加-/admin/news/newsaddcategory
func (c *NewsController) NewsAddCategory() {
	o := orm.NewOrm()
	if c.IsAjax() {
		s := models.CrNewscategory{}
		c.ParseForm(&s)
		_, err := o.Insert(&s)
		if err != nil {
			c.RtJsonString(-1, "添加失败", err.Error())

		} else {
			c.RtAjaxSuccess("添加成功")
		}

	}
	var m []orm.Params
	o.Raw("select id,category from cr_newscategory where pid=0 order by sorts desc").Values(&m)
	c.Data["data"] = m
	c.TplName = "news/category/addcategory.html"

}

//编辑文章分类-index/editcategory
func (c *NewsController) NewsEditCategory() {
	o := orm.NewOrm()
	s := models.CrNewscategory{}

	if c.IsAjax() {
		c.ParseForm(&s)
		err := models.UpdateCrNewscategoryById(&s)
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
	o.Raw("select id,category from cr_newscategory where pid=0 order by sorts desc").Values(&m)
	o.Raw("select * from cr_newscategory where id=?", id).QueryRow(&s)
	c.Data["data"] = m  //主导航
	c.Data["data2"] = s //导航数据
	c.Data["pid"] = strconv.Itoa(s.Pid)
	c.TplName = "news/category/editcategory.html"

}
