package models

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
)

func init(){
	//连接数据库
	conn,err := redis.Dial("tcp",":6379")
	if err != nil {
		fmt.Println("redis链接失败")
		return

	}
	//操作数据
	/*conn.Send("set","aa","bb")
	conn.Flush()
	conn.Receive()*/
	resp,err := conn.Do("mget","userName","pCount")
	//回复助手函数(类型转换)
	result,err := redis.Values(resp,err)

	//scan函数
	var userName string
	var pCount int
	redis.Scan(result,&userName,&pCount)
	fmt.Println("userName: ",userName,"   pCount: ",pCount)

}