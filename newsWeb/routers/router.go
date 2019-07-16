package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    //路由过滤器
    beego.InsertFilter("/article/*",beego.BeforeExec,filters)

    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleReg")
    //登录业务
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //首页展示
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex;post:HandleIndex")
    //添加文章
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    //文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    //编辑文章
    beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    //删除文章
    beego.Router("/article/delete",&controllers.ArticleController{},"get:Delete")
    //添加文章类型
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    //退出登录
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
    //删除类型
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
}

func filters(ctx *context.Context){
    //检查是否登录
    userName := ctx.Input.Session("userName")
    if userName == nil{
        ctx.Redirect(302,"/login")
        return
    }
}
