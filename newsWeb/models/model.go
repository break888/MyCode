package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//当没有设置主键的时候，默认以Id，类型为int,int8,int64的字段当主键
type User struct {
	Id int
	Name string
	Pwd string

	Articles []*Article `orm:"reverse(many)"`
}

//文章表  orm建表的时候默认非空
type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(50);unique"`
	Content string `orm:"size(500)"`
	Time time.Time `orm:"type(datetime);auto_now_add"`
	ReadCount int `orm:"default(0)"`
	Image string  `orm:"null"`

	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`
	Users []*User `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(50)"`

	Articles []*Article `orm:"reverse(many)"`
}

func init(){
	//1,注册数据库
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsWeb")
	//2.注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//3.跑起来
	orm.RunSyncdb("default",false,true)
}

