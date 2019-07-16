package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type ArticleController struct {
	beego.Controller
}

//展示首页
func(this*ArticleController)ShowIndex(){
	//获取数据并展示
	//高级查询
	//查看登录状态
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}



	o := orm.NewOrm()
	//获取类型
	typeName := this.GetString("select")




	var articles []models.Article


	qs := o.QueryTable("Article")  //queryseter  查询集合
	//qs.All(&articles)

	//获取总记录数
	var count int64
	if typeName == ""{
		//获取所有带类型文章的个数
		count,_ = qs.RelatedSel("ArticleType").Count()
	}else {
		//获取选中类型的文章个数
		//select * from artice where articleType.TypeName = typeName
		count,_ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	}
	//获取总页数
	pageSize := 2
	//math

	pageCount := math.Ceil(float64(count) / float64(pageSize))

	//获取当前页码
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	//limit  获取部分数据
	//一对多查询的时候，orm默认是惰性查询,relatedSel，一旦加上relatedSel之后，将只查询外键有值的数据
	if typeName == "" {
		qs.Limit(pageSize,( pageIndex - 1) * pageSize ).RelatedSel("ArticleType").All(&articles)
	}else {
		qs.Limit(pageSize,( pageIndex - 1) * pageSize ).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

	}



	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count
	this.Data["articles"] = articles
	this.Data["pageIndex"] = pageIndex




	//把数据存到redis中  第一次尝试从redis中获取数据,如果获取到数据，说明不是第一次，如果没有数据说明是第一次访问
	conn,err := redis.Dial("tcp",":6379")
	if err != nil {
		fmt.Println("redis数据库链接失败",err)
		return
	}

	result,err := redis.Bytes(conn.Do("get","articleTypes"))
	var articleTypes []models.ArticleType

	if len(result) == 0{
		//获取所有类型
		o.QueryTable("ArticleType").All(&articleTypes)


		//把数据写入到redis数据库中
		var buffer bytes.Buffer
		//获取编码器
		enc := gob.NewEncoder(&buffer)
		//编码
		enc.Encode(articleTypes)
		conn.Do("set","articleTypes",buffer.Bytes())

		fmt.Println("从mysql中获取数据")
	}else {
		//获取解码器
		dec := gob.NewDecoder(bytes.NewReader(result))

		dec.Decode(&articleTypes)
		fmt.Println("从redis中获取数据",articleTypes)
	}






	//序列化
	/*//要有一个容器接受序列化之后的值
	var buffer bytes.Buffer
	//要有一个编码器
	enc := gob.NewEncoder(&buffer)
	//编码
	enc.Encode(articleTypes)
	conn.Do("set","articleType",buffer.Bytes())


	//获取数据  反序列化
	var newTypes []models.ArticleType
	result,err := redis.Bytes(conn.Do("get","articleType"))

	dec := gob.NewDecoder(bytes.NewReader(result))
	dec.Decode(&newTypes)
	fmt.Println("获取到的数据为:",newTypes)*/




	errmsg := this.GetString("errmsg")


	this.Data["articleTypes"] = articleTypes
	this.Data["typeName"] = typeName
	this.Data["errmsg"] = errmsg

	this.Layout = "layout.html"
	this.TplName = "index.html"
}


//处理首页数据
func(this*ArticleController)HandleIndex(){

}

//添加文章页面
func(this*ArticleController)ShowAddArticle(){
	//获取所有类型
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes

	this.Layout="layout.html"
	this.TplName = "add.html"
}

//添加文章业务处理   插入数据
func(this*ArticleController)HandleAddArticle(){
	//1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	_,head ,err:= this.GetFile("uploadname")
	typeName := this.GetString("select")

	//校验数据
	if articleName == "" || content == "" || err != nil {
		fmt.Println("获取数据错误",err)
		this.TplName = "add.html"
		return
	}

	//上传文件一般需要校验
	//1.文件大小
	if head.Size > 10000000{
		this.Data["errmsg"] = "图片太大，请重新选择"
		this.TplName = "add.html"
		return
	}
	//2.校验文件格式
	ext := path.Ext(head.Filename)
	fmt.Println("当前文章格式为:",ext)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "文件格式错误，请重新选择"
		this.TplName = "add.html"
		return
	}
 	//3.防止重名
 	fileName := time.Now().Format("20060102150405")
 	err = this.SaveToFile("uploadname","static/img/"+fileName+ext)
 	if err != nil {
		this.Data["errmsg"] = "存储文件失败，请重新插入文章"
		this.TplName = "add.html"
		return
	}


	//把数据添加到数据库
	o := orm.NewOrm()
	var article models.Article
	article.Title  = articleName
	article.Content = content
	article.Image = "/static/img/"+fileName+ext

	//类型名称  根据类型名称获取类型对象
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")


	article.ArticleType = &articleType
	_, err = o.Insert(&article)
	if err != nil {
		fmt.Println("插入文章失败")
		this.TplName = "add.html"
		return
	}


	//成功回到首页
	this.Redirect("/article/index",302)
}

//展示文章详情页
func(this*ArticleController)ShowContent(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {

		//this.TplName = "index.html"
		//如果页面有数据，不能直接渲染，适合用跳转
		this.Redirect("/article/index?errmsg=文章详情获取数据失败",302)
		return
	}
	//处理数据  获取数据并展示
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		this.Redirect("/article/index",302)
		return
	}

	//关联多对多关系
	//o.LoadRelated(&article,"Users")
	var users []models.User
	//queryTable  需要获取哪些数据，指定哪个表（多对多属性名__表名__比较的属性）
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"] = users


	//把阅读次数加一
	article.ReadCount += 1
	o.Update(&article)


	//在登录的情况下点击查看详情
	userName := this.GetSession("userName")
	//多对多插入
	//获取插入对象

	//获取多对多操作对象
	m2m := o.QueryM2M(&article,"Users")
	//获取插入对象
	var user models.User
	user.Name = userName.(string)
	o.Read(&user,"Name")
	//插入
	m2m.Add(user)



	//返回数据
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "content.html"
}

//展示编辑页面
func(this*ArticleController)ShowUpdate(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		this.Redirect("/article/index?errmsg=文章编辑页面获取信息错误",302)
		return
	}
	//处理数据
	//查询数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		this.Redirect("/article/index?errmsg=文章编辑页面获取信息错误",302)
		return
	}

	//返回数据
	this.Data["article"] = article
	this.TplName = "update.html"
}

func UploadFile(this*ArticleController,fileImage string)string{
	_,head ,err:= this.GetFile(fileImage)

	if err != nil {
		this.Data["errmsg"] = "获取文件失败，请重新添加"
		this.TplName = "add.html"
		return ""
	}

	//上传文件一般需要校验
	//1.文件大小
	if head.Size > 10000000{
		this.Data["errmsg"] = "图片太大，请重新选择"
		this.TplName = "add.html"
		return ""
	}
	//2.校验文件格式
	ext := path.Ext(head.Filename)
	fmt.Println("当前文章格式为:",ext)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "文件格式错误，请重新选择"
		this.TplName = "add.html"
		return ""
	}
	//3.防止重名
	fileName := time.Now().Format("20060102150405")
	err = this.SaveToFile(fileImage,"static/img/"+fileName+ext)
	if err != nil {
		this.Data["errmsg"] = "存储文件失败，请重新插入文章"
		this.TplName = "add.html"//查看是否有用
		return ""
	}
	return "static/img/"+fileName+ext
}

//处理编辑数据
func(this*ArticleController)HandleUpdate(){
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//函数的封装
	filePath := UploadFile(this,"uploadname")
	id,err := this.GetInt("id")
	//隐藏域传值
	if articleName == "" || content == "" || filePath == "" || err != nil {
		fmt.Println("获取数据错误")
		this.Redirect("/article/update?id="+strconv.Itoa(id),302)
		return
	}

	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		fmt.Println("获取失败",err)
		this.Redirect("/article/update?id="+strconv.Itoa(id),302)
		return
	}
	//赋新值
	article.Title = articleName
	article.Content = content
	article.Image = filePath
	_,err = o.Update(&article)
	if err != nil {
		fmt.Println("更新失败",err)
		this.Redirect("/article/update?id="+strconv.Itoa(id),302)
		return
	}

	//返回数据
	this.Redirect("/article/index",302)

}

//删除文章
func(this*ArticleController)Delete(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		fmt.Println("删除失败",err)
		this.Redirect("/article/index?errmsg=删除文章失败",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	_,err = o.Delete(&article)

	if err != nil {
		fmt.Println("删除失败",err)
		this.Redirect("/article/index?errmsg=删除文章失败",302)
		return
	}


	//返回数据
	this.Redirect("/article/index",302)
}

//展示添加文章类型页面
func(this*ArticleController)ShowAddType(){
	//获取类型数据
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes

	this.Layout = "layout.html"
	this.TplName = "addType.html"
}


//处理添加文章业务
func(this*ArticleController)HandleAddType(){
	//获取数据
	typeName := this.GetString("typeName")
	//校验数据
	if typeName == "" {
		this.Data["errmsg"] = "获取数据失败"
		this.TplName = "addType.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	_,err := o.Insert(&articleType)
	if err != nil {
		this.Data["errmsg"] = "添加文章类型失败"
		this.TplName = "addType.html"
		return
	}

	//返回数据
	this.Redirect("/article/addType",302)
	//this.Dat*/
	//this.TplName = "addType.html"
}

//删除类型操作
func(this*ArticleController)DeleteType(){
	//接受数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		this.Redirect("/article/addType",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType,"Id")

	//返回数据
	this.Redirect("/article/addType",302)

}