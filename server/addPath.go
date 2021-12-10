package main

import "github.com/gin-gonic/gin"

func addPath(r *gin.Engine) {
	// 	GET：请求方式；/index：请求的路径
	// 	当客户端以GET方法请求/index路径时，会执行后面的匿名函数
	//	authMiddleWare是一个中间件，用以检查cookie来判断用户是否登录。并将用户名提交给下一个中间件checkPermission。
	//	checkPermission也是一个中间件，用以判断已登录用户所处的组别。并将用户名与组别提交给下一个匿名函数。
	r.GET("/index", welcome)           //用户登录
	r.GET("/customer_index", welcome)  //顾客登录
	r.GET("/repairman_index", welcome) //维修员登录
	r.GET("/salesman_index", welcome)  //业务员登录

	r.GET("/register", startRegister)                                                   //用户注册
	r.GET("/login", startLogin)                                                         //用户登录
	r.GET("/logout", authMiddleWare(), checkPermission(), startLogout)                  //用户登出
	r.GET("/changePassword", authMiddleWare(), checkPermission(), startChangePassword)  //更改密码
	r.GET("/checkNotification", authMiddleWare(), checkPermission(), checkNotification) //检查通知
	r.GET("/checkGroup", authMiddleWare(), checkPermission(), checkGroup)               //返回用户组

	r.GET("/test", authMiddleWare(), checkPermission(), test) //测试

	r.GET("/userinfo", authMiddleWare(), checkPermission(), userinfo)                   //查询用户个人信息
	r.GET("/change_userinfo", authMiddleWare(), checkPermission(), startChangeUserinfo) //进入更改工作状态界面
	r.GET("/addVehicle", addVehicle)                                                    //用户登出

	r.GET("/checkStatus", authMiddleWare(), checkPermission(), checkStatus)         //返回维修工状态
	r.GET("/change_status", authMiddleWare(), checkPermission(), startChangeStatus) //进入更改工作状态界面
	r.GET("/check_orders", authMiddleWare(), checkPermission(), startCheckOrders)   //进入更改工作状态界面

	// 表单方式
	r.POST("/read", authMiddleWare(), checkPermission(), read)                     //设置已读
	r.POST("/changePassword", authMiddleWare(), checkPermission(), changePassword) //后端处理更改密码
	r.POST("/logout", logout)                                                      //后端处理用户登出
	r.POST("/login", login)                                                        //后端处理用户登陆
	r.POST("/register", register)                                                  //后端处理用户注册
	r.GET("/captcha", func(c *gin.Context) { Captcha(c, 4) })                      //随机生成一个4位数字验证码

	r.POST("/changeUserinfo", authMiddleWare(), checkPermission(), changeUserinfo)

	r.POST("/changeStatus", authMiddleWare(), checkPermission(), changeStatus) //后端处理更改密码
}
