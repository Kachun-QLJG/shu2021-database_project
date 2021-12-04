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
	defer database.Close()
	database.SingularTable(true)
	database.InstantSet("gorm:table_options", "ENGINE=InnoDB")
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Salesman{})
	database.AutoMigrate(&Repairman{})
	database.AutoMigrate(&TypeOverview{})
	database.AutoMigrate(&PartsOverview{})
	database.AutoMigrate(&Vehicle{})
	database.AutoMigrate(&Attorney{})
	database.AutoMigrate(&Arrangement{})
	database.AutoMigrate(&AuthSession{})
	database.AutoMigrate(&Notification{})

	database.Model(&Vehicle{}).AddForeignKey("user_id", "user(number)", "RESTRICT", "RESTRICT")
	database.Model(&Attorney{}).AddForeignKey("user_id", "user(number)", "RESTRICT", "RESTRICT")
	database.Model(&Attorney{}).AddForeignKey("vehicle_number", "vehicle(number)", "RESTRICT", "RESTRICT")
	database.Model(&Attorney{}).AddForeignKey("salesman_id", "salesman(number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("order_number", "attorney(number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("repairman_number", "repairman(number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("project_number", "type_overview(project_number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("parts_number", "parts_overview(parts_number)", "RESTRICT", "RESTRICT")
	//database.Model(&Notification{}).AddForeignKey("user_id", "user(number)", "RESTRICT", "RESTRICT")
	//database.Model(&Notification{}).AddForeignKey("user_id", "repairman(number)", "RESTRICT", "RESTRICT")
	//database.Model(&Notification{}).AddForeignKey("user_id", "salesman(number)", "RESTRICT", "RESTRICT")

	database.DB().SetConnMaxLifetime(time.Hour * 24 * 21) //设置超时时间
	// 创建一个默认的路由引擎
	r := gin.Default()
	r.Static("/p", "./html/statics")
	r.Static("/src", "./html/src") //将相对html的路径替换成相对工程的路径
	r.LoadHTMLFiles("./html/login.html", "./html/logout.html", "./html/register.html", "./html/error.html",
		"./html/success.html", "./html/index.html", "./html/change_password.html", "./html/customer_index.html",
		"./html/salesman_index.html", "./html/repairman_index.html") //加载html模板
	r.Use(Session("SHU")) //验证码生成会使用SHU作为密钥生成session
	// 	GET：请求方式；/index：请求的路径
	// 	当客户端以GET方法请求/index路径时，会执行后面的匿名函数
	//	authMiddleWare是一个中间件，用以检查cookie来判断用户是否登录。并将用户名提交给下一个中间件checkPermission。
	//	checkPermission也是一个中间件，用以判断已登录用户所处的组别。并将用户名与组别提交给下一个匿名函数。
	r.GET("/index", welcome)                                                            //用户登录
	r.GET("/customer_index", welcome)                                                   //顾客登录
	r.GET("/repairman_index", welcome)                                                  //维修员登录
	r.GET("/salesman_index", welcome)                                                   //业务员登录
	r.GET("/register", startRegister)                                                   //用户注册
	r.GET("/login", startLogin)                                                         //用户登录
	r.GET("/logout", authMiddleWare(), checkPermission(), startLogout)                  //用户登出
	r.GET("/addVehicle", addVehicle)                                                    //用户登出
	r.GET("/changePassword", authMiddleWare(), checkPermission(), startChangePassword)  //更改密码
	r.GET("/checkNotification", authMiddleWare(), checkPermission(), checkNotification) //检查通知
	r.GET("/checkGroup", authMiddleWare(), checkPermission(), checkGroup)               //返回用户组
	r.GET("/checkStatus", authMiddleWare(), checkPermission(), checkStatus)             //返回维修工状态
	r.GET("/test", authMiddleWare(), checkPermission(), test)                           //测试

	r.POST("/changeStatus", authMiddleWare(), checkPermission(), changeStatus)     //后端处理更改密码
	r.POST("/read", authMiddleWare(), checkPermission(), read)                     //设置已读
	r.POST("/changePassword", authMiddleWare(), checkPermission(), changePassword) //后端处理更改密码
	r.POST("/logout", logout)                                                      //后端处理用户登出
	r.POST("/login", login)                                                        //后端处理用户登陆
	r.POST("/register", register)                                                  //后端处理用户注册
	r.GET("/captcha", func(c *gin.Context) { Captcha(c, 4) })                      //随机生成一个4位数字验证码

	// 启动HTTP服务，在36b1c95548.qicp.vip启动服务
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("启动HTTP服务失败：", err)
	}
}
