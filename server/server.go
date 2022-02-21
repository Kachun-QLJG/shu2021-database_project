package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var database, databaseERR = gorm.Open("mysql", "admin:123456@(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local")

//连接mysql数据库

func main() {
	if databaseERR != nil {
		panic(databaseERR)
	}
	defer func(database *gorm.DB) {
		err := database.Close()
		if err != nil {
		}
	}(database)
	connectToSql(database)
	// database.LogMode(true)
	database.DB().SetConnMaxLifetime(time.Hour * 24 * 21) //设置连接数据库超时时间
	// 创建一个默认的路由引擎
	r := gin.Default()
	r.Static("/statics", "./html/statics")
	r.Static("/src", "./html/src") //将相对html的路径替换成相对工程的路径
	r.LoadHTMLFiles("./html/login.html", "./html/register.html",
		"./html/user/user_index.html", "./html/user/user_check_orders_ongoing.html", "./html/user/user_check_orders_finished.html",
		"./html/salesman/salesman_index.html", "./html/repairman/repairman_index.html") //加载html模板
	r.Use(Session("SHU")) //验证码生成会使用SHU作为密钥生成session
	addPath(r)

	// 启动HTTP服务，在36b1c95548.qicp.vip启动服务
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("启动HTTP服务失败：", err)
	}
}
