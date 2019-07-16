package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

//展示注册页面
func(c*UserController)ShowRegister(){
	c.TplName = "register.html"
}

//处理注册业务
func(this*UserController)HandleReg(){
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//fmt.Println("用户名为，",userName,"密码为，",pwd)
	//2.校验数据
	if userName == "" || pwd =="" {
		fmt.Println("用户名或密码不能为空")
		this.TplName = "register.html"
		return
	}
	//3.处理数据 把数据插入数据库
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.Pwd = pwd
	_,err := o.Insert(&user)
	if err != nil {
		fmt.Println("注册失败,",err)
		this.TplName = "register.html"
		return
	}

	//4.返回数据  返回一句话给前段
	//this.TplName = "login.html"
	this.Redirect("/login",302)

}

//展示登录界面
func(this*UserController)ShowLogin(){
	//获取cookie
	userName := this.Ctx.GetCookie("userName")
	if userName != "" {
		dec,_ :=base64.StdEncoding.DecodeString(userName)
		this.Data["userName"] = string(dec)
		this.Data["checked"] = "checked"
	}else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}
	this.TplName = "login.html"
}

//处理登录请求
func(this*UserController)HandleLogin(){
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "用户名或者密码不能为空"
		this.TplName = "login.html"
		return
	}
	//3.处理数据
	//查询操作
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user,"Name")
	if err != nil {
		this.Data["errmsg"] = "当前账户不存在"
		this.TplName = "login.html"
		return
	}

	if user.Pwd != pwd{
		this.Data["errmsg"] = "用户密码错误"
		this.TplName = "login.html"
		return
	}

	//登录成功并且点击记住用户名
	remember := this.GetString("remember")

	if remember == "on"{
		enc := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",enc,60 * 60)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}


	//登录成功存储登录状态
	this.SetSession("userName",userName)
	//this.Ctx.SetCookie("userName",userName,3600)
	//4.返回数据
	//this.TplName = "index.html"
	this.Redirect("/article/index",302)
}

//登录登录
func(this*UserController)Logout(){
	this.DelSession("userName")
	this.Redirect("/login",302)
}