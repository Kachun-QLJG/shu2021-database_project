package main

import "github.com/gin-gonic/gin"

func addPath(r *gin.Engine) {
	//==================================GET方法==================================
	//	authMiddleWare是一个中间件，用以检查cookie来判断用户是否登录。并将用户名提交给下一个中间件checkPermission。
	//	checkPermission也是一个中间件，用以判断已登录用户所处的组别。并将用户名与组别提交给下一个匿名函数。
	//----------全局路径----------
	r.GET("/index", welcome)                                                             //根据登录状态返回不同主页
	r.GET("/register", startRegister)                                                    //用户注册
	r.GET("/check_register", checkRegister)                                              //检查用户名是否重复
	r.GET("/login", startLogin)                                                          //用户登录
	r.GET("/logout", authMiddleWare(), checkPermission(), startLogout)                   //用户登出
	r.GET("/check_notification", authMiddleWare(), checkPermission(), checkNotification) //检查通知
	r.GET("/get_username", authMiddleWare(), checkPermission(), getUsername)             //返回用户的用户名
	r.GET("/download_pdf", authMiddleWare(), checkPermission(), downloadPdf)             //生成并下载PDF文档
	r.GET("/show_pdf", authMiddleWare(), checkPermission(), showPdf)                     //返回pdf文档
	r.GET("/generate_pdf", authMiddleWare(), checkPermission(), generatePdf)             //【临时】生成PDF
	r.GET("/captcha", func(c *gin.Context) { Captcha(c, 4) })                            //随机生成一个4位数字验证码

	//----------用户路径----------
	r.GET("/get_user_info", authMiddleWare(), checkPermission(), filterUser(), getUserInfo)                         //查询用户个人信息
	r.GET("/show_plate", authMiddleWare(), checkPermission(), filterUser(), showPlate)                              //返回车牌
	r.GET("/get_vehicle", authMiddleWare(), checkPermission(), filterUser(), getVehicle)                            //获取车辆（所有）
	r.GET("/repair_history", authMiddleWare(), checkPermission(), filterUser(), getRepairHistory)                   //获取某一车辆历史维修记录
	r.GET("/get_processing_attorney", authMiddleWare(), checkPermission(), filterUser(), getProcessingAttorney)     //寻找进行中的委托
	r.GET("/get_finished_attorney", authMiddleWare(), checkPermission(), filterUser(), getFinishedAttorney)         //寻找已完成的委托（用户）
	r.GET("/get_attorney_detail", authMiddleWare(), checkPermission(), filterUser(), getAttorneyDetail)             //获取某一委托的详情
	r.GET("/check_vehicle", authMiddleWare(), checkPermission(), filterUser(), checkVehicle)                        //检查车辆是否已被绑定
	r.GET("/u_check_orders_ongoing", authMiddleWare(), checkPermission(), filterUser(), startUCheckOrdersOngoing)   //进入用户查看进行中订单页面
	r.GET("/u_check_orders_finished", authMiddleWare(), checkPermission(), filterUser(), startUCheckOrdersFinished) //进入用户查看已完成订单页面
	r.GET("/get_full_attorney", authMiddleWare(), checkPermission(), filterUser(), getFullAttorney)                 //返回完整委托书·所有内容（用戶）

	//----------业务员路径----------
	r.GET("/get_project_time", authMiddleWare(), checkPermission(), filterSalesman(), getProjectTime)                       //业务员根据维修项目来获取预计维修时间
	r.GET("/search_for_projects", authMiddleWare(), checkPermission(), filterSalesman(), searchForProjects)                 //寻找维修项目
	r.GET("/get_salesman_info", authMiddleWare(), checkPermission(), filterSalesman(), getSalesmanInfo)                     //获取业务员个人信息
	r.GET("/get_relating_attorney", authMiddleWare(), checkPermission(), filterSalesman(), getRelatingAttorney)             //业务员获与自己相关的订单
	r.GET("/s_check_orders", authMiddleWare(), checkPermission(), filterSalesman(), startSCheckOrders)                      //进入业务员查看订单页面
	r.GET("/take_orders", authMiddleWare(), checkPermission(), filterSalesman(), startTakeOrders)                           //进入接单界面
	r.GET("/get_corresponding_repairman", authMiddleWare(), checkPermission(), filterSalesman(), getCorrespondingRepairman) //根据维修员种类返回相应种类的维修员信息
	//r.GET("/get_finished_attorney_s", authMiddleWare(), checkPermission(), filterSalesman(), getFinishedAttorneyS)          //寻找已完成的委托（业务员）
	r.GET("/get_full_attorney_s", authMiddleWare(), checkPermission(), filterSalesman(), getFullAttorneyS) //获得委托的详细信息（业务员）

	//----------维修员路径----------
	r.GET("/search_for_parts", authMiddleWare(), checkPermission(), filterRepairman(), searchForParts)                     //寻找零件
	r.GET("/get_repairman_info", authMiddleWare(), checkPermission(), filterRepairman(), getRepairmanInfo)                 //获取业务员个人信息
	r.GET("/check_status", authMiddleWare(), checkPermission(), filterRepairman(), checkStatus)                            //返回维修工状态
	r.GET("/get_processing_arrangement", authMiddleWare(), checkPermission(), filterRepairman(), getProcessingArrangement) //返回正在处理的维修项目
	r.GET("/get_finished_arrangement", authMiddleWare(), checkPermission(), filterRepairman(), getFinishedArrangement)     //返回已完成的维修项目

	//==================================POST方法==================================
	//----------全局路径----------
	r.POST("/read", authMiddleWare(), checkPermission(), read)                      //设置消息已读
	r.POST("/change_password", authMiddleWare(), checkPermission(), changePassword) //后端处理更改密码
	r.POST("/logout", logout)                                                       //后端处理用户登出
	r.POST("/login", login)                                                         //后端处理用户登陆
	r.POST("/register", register)                                                   //后端处理用户注册

	//----------用户路径----------
	r.POST("/create_attorney", authMiddleWare(), checkPermission(), filterUser(), createAttorney) //用户创建新的委托
	r.POST("/change_userinfo", authMiddleWare(), checkPermission(), filterUser(), changeUserinfo) //后端处理更改用户个人信息
	r.POST("/add_vehicle", authMiddleWare(), checkPermission(), filterUser(), addVehicle)         //后端处理用户添加车辆

	//----------业务员路径----------
	r.POST("/add_project_for_attorney", authMiddleWare(), checkPermission(), filterSalesman(), addProjectForAttorney) //业务员为委托添加维修项目
	r.POST("/set_attorney_finished", authMiddleWare(), checkPermission(), filterSalesman(), setAttorneyFinished)      //业务员结算委托书并设置进度为已完成
	r.POST("/receive_attorney", authMiddleWare(), checkPermission(), filterSalesman(), receiveAttorney)               //接单
	r.POST("/change_discount_rate", authMiddleWare(), checkPermission(), filterSalesman(), changeDiscountRate)        //更改用户折扣率

	//----------维修员路径----------
	r.POST("/change_status", authMiddleWare(), checkPermission(), filterRepairman(), changeStatus)                //后端处理维修员更改工作状态
	r.POST("/add_parts_for_project", authMiddleWare(), checkPermission(), filterRepairman(), addPartsForProject)  //维修工为项目添加零件
	r.POST("/change_repair_status", authMiddleWare(), checkPermission(), filterRepairman(), changeRepairProgress) //后端处理维修员更改维修项目状态【开始工作、已完成等】
}
