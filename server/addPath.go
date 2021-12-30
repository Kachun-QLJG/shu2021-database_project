package main

import "github.com/gin-gonic/gin"

func addPath(r *gin.Engine) {
	// 	GET：请求方式；/index：请求的路径
	// 	当客户端以GET方法请求/index路径时，会执行后面的匿名函数
	//	authMiddleWare是一个中间件，用以检查cookie来判断用户是否登录。并将用户名提交给下一个中间件checkPermission。
	//	checkPermission也是一个中间件，用以判断已登录用户所处的组别。并将用户名与组别提交给下一个匿名函数。
	r.GET("/index", welcome)                                                              //根据登录状态返回不同主页
	r.GET("/register", startRegister)                                                     //用户注册
	r.GET("/check_register", checkRegister)                                               //用户注册
	r.GET("/login", startLogin)                                                           //用户登录
	r.GET("/logout", authMiddleWare(), checkPermission(), startLogout)                    //用户登出
	r.GET("/change_password", authMiddleWare(), checkPermission(), startChangePassword)   //更改密码
	r.GET("/check_notification", authMiddleWare(), checkPermission(), checkNotification)  //检查通知
	r.GET("/check_group", authMiddleWare(), checkPermission(), checkGroup)                //返回用户组
	r.GET("/get_username", authMiddleWare(), checkPermission(), getUsername)              //返回用户名
	r.GET("/get_group", authMiddleWare(), checkPermission(), getGroup)                    //返回用户组
	r.GET("/get_pdf", authMiddleWare(), checkPermission(), getPdf)                        //生成并获取PDF文档
	r.GET("/show_pdf", authMiddleWare(), checkPermission(), showPdf)                      //打印车牌
	r.GET("/show_plate", authMiddleWare(), checkPermission(), showPlate)                  //返回车牌
	r.GET("/get_vehicle", authMiddleWare(), checkPermission(), getVehicle)                //获取车辆（所有）
	r.GET("/repair_history", authMiddleWare(), checkPermission(), getRepairHistory)       //获取某一车辆历史维修记录
	r.GET("/search_for_projects", authMiddleWare(), checkPermission(), searchForProjects) //寻找维修项目
	r.GET("/search_for_parts", authMiddleWare(), checkPermission(), searchForParts)       //寻找零件

	r.GET("/userinfo", authMiddleWare(), checkPermission(), userinfo)                   //查询用户个人信息
	r.GET("/change_userinfo", authMiddleWare(), checkPermission(), startChangeUserinfo) //进入更改个人信息界面
	r.GET("/add_vehicle", authMiddleWare(), checkPermission(), startAddVehicle)         //用户添加车辆
	r.GET("/check_vehicle", authMiddleWare(), checkPermission(), checkVehicle)          //用户添加车辆
	r.GET("/u_check_orders", authMiddleWare(), checkPermission(), startUCheckOrders)    //进入用户查看订单页面

	r.GET("/check_status", authMiddleWare(), checkPermission(), checkStatus)         //返回维修工状态
	r.GET("/change_status", authMiddleWare(), checkPermission(), startChangeStatus)  //进入更改工作状态界面
	r.GET("/r_check_orders", authMiddleWare(), checkPermission(), startRCheckOrders) //进入维修员查看订单页面

	r.GET("/s_check_orders", authMiddleWare(), checkPermission(), startSCheckOrders) //进入业务员查看订单页面
	r.GET("/take_orders", authMiddleWare(), checkPermission(), startTakeOrders)      //进入接单界面

	// 表单方式
	r.POST("/read", authMiddleWare(), checkPermission(), read)                      //设置已读
	r.POST("/change_password", authMiddleWare(), checkPermission(), changePassword) //后端处理更改密码
	r.POST("/logout", logout)                                                       //后端处理用户登出
	r.POST("/login", login)                                                         //后端处理用户登陆
	r.POST("/register", register)                                                   //后端处理用户注册
	r.GET("/captcha", func(c *gin.Context) { Captcha(c, 4) })                       //随机生成一个4位数字验证码

	r.POST("/change_userinfo", authMiddleWare(), checkPermission(), changeUserinfo) //后端处理更改用户个人信息
	r.POST("/add_vehicle", authMiddleWare(), checkPermission(), addVehicle)         //后端处理用户添加车辆
	r.POST("/change_status", authMiddleWare(), checkPermission(), changeStatus)     //后端处理维修员更改工作状态
}
