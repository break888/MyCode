package main

import (
	_ "newsWeb/models"
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
)

func main() {

	beego.AddFuncMap("prePage",ShowPrePage)
	beego.AddFuncMap("nextPage",ShowNextPage)
	beego.Run()
}

func ShowPrePage(pageIndex int)int{
	if pageIndex <= 1{
		return 1
	}
	return pageIndex - 1
}

//下一页函数
func ShowNextPage(pageIndex int,pageCount int)int{
	if pageIndex >= pageCount{
		return pageCount
	}
	return pageIndex + 1
}

