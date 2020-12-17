package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r:=gin.Default()
	r.GET("/home", func(context *gin.Context) {
		context.JSON(http.StatusOK,"pong")
	})
	err:=r.Run(":8080")
	if err!=nil{
		fmt.Println("gin run error=",err)
		os.Exit(-1)
	}
}
