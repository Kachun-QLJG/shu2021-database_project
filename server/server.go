package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type User struct { //客户表
	Number        string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //客户编号
	Password      string `gorm:"type:varchar(255);not null"`                        //密码
	Name          string `gorm:"type:varchar(100)"`                                 //客户名称
	Property      string `gorm:"type:varchar(4)"`                                   //性质
	DiscountRate  int    `gorm:"type:int(2)"`                                       //折扣率
	ContactPerson string `gorm:"type:varchar(10)"`                                  //联系人
	ContactTel    string `gorm:"type:varchar(20);not null;unique_index"`            //联系电话
}

type Salesman struct { //业务员表
	Number   string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //工号
	Name     string `gorm:"type:varchar(20);not null"`                         //姓名
	Password string `gorm:"type:varchar(255);not null"`                        //密码
}

type Repairman struct { //维修员表
	Number          string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //工号
	Name            string `gorm:"type:varchar(20);not null"`                         //姓名
	Password        string `gorm:"type:varchar(255);not null"`                        //密码
	Type            string `gorm:"type:varchar(10);not null"`                         //工种
	CurrentWorkHour int    `gorm:"type:int(4);not null"`                              //当前工时
	Status          string `gorm:"type:varchar(10);not null"`                         //工人状态
}

type TypeOverview struct { //工种总览表
	ProjectNumber string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //维修项目编号
	ProjectName   string `gorm:"type:varchar(100);not null"`                        //维修项目名称
	Type          string `gorm:"type:varchar(10);not null;primary_key"`             //工种
}

type PartsOverview struct { //零件总览表
	PartsNumber string  `gorm:"type:varchar(8);unique_index;not null;primary_key"` //零件编号
	PartsName   string  `gorm:"type:varchar(50);not null"`                         //零件名称
	Unit        string  `gorm:"type:varchar(6);not null"`                          //计量单位
	PartsCost   float64 `gorm:"type:double(8,2);not null"`                         //零件价格
}

type Vehicle struct { //车辆表
	Number        string `gorm:"type:varchar(17);not null;primary_key"` //车架号
	LicenseNumber string `gorm:"type:varchar(10);not null"`             //车牌号
	UserID        string `gorm:"type:varchar(8);not null;primary_key"`  //客户编号
	Color         string `gorm:"type:varchar(10);not null"`             //车辆颜色
	Model         string `gorm:"type:varchar(40);not null"`             //车型
	Type          string `gorm:"type:varchar(10);not null"`             //车辆类别
	Time          string `gorm:"type:varchar(20);not null"`             //绑定时间
}

type Notification struct { //通知表
	Number  string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //通知编号
	UserID  string `gorm:"type:varchar(8);not null"`                          //用户号
	Title   string `gorm:"type:varchar(50);not null"`                         //通知标题
	Content string `gorm:"type:varchar(255);not null"`                        //通知内容
	Status  string `gorm:"type:varchar(4);not null"`                          //接收状态
}

type Attorney struct { //委托书表
	Number            string  `gorm:"type:varchar(11);not null;unique_index;primary_key"` //订单号
	UserID            string  `gorm:"type:varchar(8);not null"`                           //客户编号
	VehicleNumber     string  `gorm:"type:varchar(17);not null"`                          //车架号
	RepairType        string  `gorm:"type:varchar(4)"`                                    //维修类型
	Classification    string  `gorm:"type:varchar(4)"`                                    //作业分类
	PayMethod         string  `gorm:"type:varchar(4)"`                                    //结算方式
	StartTime         string  `gorm:"type:varchar(20);not null"`                          //进场时间
	SalesmanID        string  `gorm:"type:varchar(8)"`                                    //业务员编号
	PredictFinishTime string  `gorm:"type:varchar(20)"`                                   //预计完工时间时间
	ActualFinishTime  string  `gorm:"type:varchar(20)"`                                   //实际完工时间时间
	RoughProblem      string  `gorm:"type:varchar(255);not null"`                         //粗略故障描述
	SpecificProblem   string  `gorm:"type:varchar(255);not null"`                         //详细故障描述
	Progress          string  `gorm:"type:varchar(10);not null"`                          //进展
	TotalCost         float64 `gorm:"type:double(6,2);not null"`                          //总价
	StartPetrol       float64 `gorm:"type:double(5,2);not null"`                          //进厂油量
	StartMile         float64 `gorm:"type:double(8,2);not null"`                          //进厂里程
	EndPetrol         float64 `gorm:"type:double(5,2)"`                                   //出厂油量
	EndMile           float64 `gorm:"type:double(8,2)"`                                   //出厂里程
	OutRange          string  `gorm:"type:varchar(255)"`                                  //非维修范围
}

type Arrangement struct { //派工单表
	OrderNumber     string `gorm:"type:varchar(11);not null;primary_key"` //订单号
	ProjectNumber   string `gorm:"type:varchar(8);not null;primary_key"`  //维修项目编号
	PredictTime     int    `gorm:"type:int(3);not null"`                  //预计工时
	ActualTIme      int    `gorm:"type:int(3);not null"`                  //实际工时
	RepairmanNumber string `gorm:"type:varchar(8);not null;primary_key"`  //维修工工号
	PartsNumber     string `gorm:"type:varchar(8);not null"`              //零件号
	PartsCount      int    `gorm:"type:int(2);not null"`                  //零件数量
	Progress        string `gorm:"type:varchar(6);not null"`              //进展
}

type AuthSession struct { //登录表
	TimeHash  string `gorm:"type:varchar(64);not null;unique_index;primary_key"` //时间戳的哈希值
	LastVisit string `gorm:"type:varchar(30);not null"`                          //最后一次访问的时间戳（精确到秒）
	Username  string `gorm:"type:varchar(20);not null"`                          //当前session对应的用户信息
}

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
	r.Static("/statics", "./html/statics")
	r.Static("/src", "./html/src") //将相对html的路径替换成相对工程的路径
	r.LoadHTMLFiles("./html/login.html", "./html/logout.html", "./html/register.html", "./html/error.html",
		"./html/success.html", "./html/index.html", "./html/change_password.html") //加载html模板
	r.Use(Session("SHU")) //验证码生成会使用SHU作为密钥生成session
	// 	GET：请求方式；/index：请求的路径
	// 	当客户端以GET方法请求/index路径时，会执行后面的匿名函数
	//	authMiddleWare是一个中间件，用以检查cookie来判断用户是否登录。并将用户名提交给下一个中间件checkPermission。
	//	checkPermission也是一个中间件，用以判断已登录用户所处的组别。并将用户名与组别提交给下一个匿名函数。
	r.GET("/index", welcome)                                                            //用户登录
	r.GET("/register", startRegister)                                                   //用户注册
	r.GET("/login", startLogin)                                                         //用户登录
	r.GET("/logout", authMiddleWare(), checkPermission(), startLogout)                  //用户登出
	r.GET("/addVehicle", addVehicle)                                                    //用户登出
	r.GET("/changePassword", authMiddleWare(), checkPermission(), startChangePassword)  //更改密码
	r.GET("/checkNotification", authMiddleWare(), checkPermission(), checkNotification) //检查通知
	r.GET("/checkGroup", authMiddleWare(), checkPermission(), checkGroup)               //返回用户组
	r.GET("/checkStatus", authMiddleWare(), checkPermission(), checkStatus)             //返回维修工状态

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

func checkStatus(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	var repairman Repairman
	database.First(&repairman, "number = ?", number)
	c.String(http.StatusOK, repairman.Status)
}

func changeStatus(c *gin.Context) {
	status := c.PostForm("status")
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.String(http.StatusBadRequest, "错误！")
		return
	}
	var repairman Repairman
	database.First(&repairman, "number = ?", number)
	database.Model(&repairman).Update("status", status) //更改状态为传过来的状态
	c.String(http.StatusOK, "修改成功！")
}

func checkGroup(c *gin.Context) {
	group := c.MustGet("group").(string)
	c.String(http.StatusOK, group)
}

func read(c *gin.Context) {
	username := c.MustGet("username").(string)
	var notification Notification
	database.First(&notification, "user_id = ? and status = ?", username, "未读")
	database.Model(&notification).Update("status", "已读") //更改消息为已读
}

func checkNotification(c *gin.Context) {
	username := c.MustGet("username").(string)
	var notification Notification
	result := database.First(&notification, "user_id = ? and status = ?", username, "未读")
	if result.RowsAffected == 1 {
		c.JSON(http.StatusOK, gin.H{"title": notification.Title, "content": notification.Content})
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func startChangePassword(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	c.HTML(http.StatusOK, "change_password.html", gin.H{"username": username, "group": group})
}

func changePassword(c *gin.Context) {
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group == "普通用户" {
		var user User
		database.First(&user, "contact_tel = ?", username)
		if CheckPasswordHash(oldPassword, user.Password) { //密码比对通过
			secretPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			database.Model(&user).Update("password", secretPassword) //用一次session，更新一次时间。
			username1 := "[out]" + username
			var session AuthSession
			database.First(&session, "username = ?", username)
			database.Model(&session).Update("username", username1) //在session表中将用户的账号前加入[out]标识
			c.SetCookie("sessionId", "", 0, "", "", false, true)   //清除浏览器中的cookie
			c.HTML(http.StatusOK, "success.html", gin.H{"data": "密码更改成功！", "website": "/login", "webName": "登录页面"})
		} else {
			c.HTML(http.StatusOK, "error.html", gin.H{"data": "密码错误！", "website": "/change_password", "webName": "修改密码页面"})
		}
	} else if group == "业务员" {
		var user Salesman
		database.First(&user, "number = ?", username)
		if CheckPasswordHash(oldPassword, user.Password) { //密码比对通过
			secretPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			database.Model(&user).Update("password", secretPassword) //用一次session，更新一次时间。
			username1 := "[out]" + username
			var session AuthSession
			database.First(&session, "username = ?", username)
			database.Model(&session).Update("username", username1) //在session表中将用户的账号前加入[out]标识
			c.SetCookie("sessionId", "", 0, "", "", false, true)   //清除浏览器中的cookie
			c.HTML(http.StatusOK, "success.html", gin.H{"data": "密码更改成功！", "website": "/login", "webName": "登录页面"})
		} else {
			c.HTML(http.StatusOK, "error.html", gin.H{"data": "密码错误！", "website": "/change_password", "webName": "修改密码页面"})
		}
	} else {
		var user Repairman
		database.First(&user, "number = ?", username)
		if CheckPasswordHash(oldPassword, user.Password) { //密码比对通过
			secretPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			database.Model(&user).Update("password", secretPassword) //用一次session，更新一次时间。
			username1 := "[out]" + username
			var session AuthSession
			database.First(&session, "username = ?", username)
			database.Model(&session).Update("username", username1) //在session表中将用户的账号前加入[out]标识
			c.SetCookie("sessionId", "", 0, "", "", false, true)   //清除浏览器中的cookie
			c.HTML(http.StatusOK, "success.html", gin.H{"data": "密码更改成功！", "website": "/login", "webName": "登录页面"})
		} else {
			c.HTML(http.StatusOK, "error.html", gin.H{"data": "密码错误！", "website": "/change_password", "webName": "修改密码页面"})
		}
	}
}

func addVehicle(c *gin.Context) {
	number := c.Query("number")
	licenseNumber := c.Query("license_number")
	userId := c.Query("user_id")
	color := c.Query("color")
	model := c.Query("model")
	carType := c.Query("type")
	sTime := time.Now().Format("2006-01-02 15:04:05")
	data := Vehicle{number, licenseNumber, userId, color, model, carType, sTime}
	err := database.Create(&data)
	strErr := fmt.Sprintf("%v", err.Error)
	if strErr != "<nil>" {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"errdata": "注册失败！" + strErr, "website": "/register", "webName": "注册页面"})
	} else {
		c.HTML(http.StatusOK, "success.html", gin.H{"data": "注册成功！", "website": "/login", "webName": "登录页面"})
	}
}

//localhost:8080/addVehicle?number=1&license_number=1&user_id=1&color=1&model=1&type=1

func welcome(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	var username, group string
	var session AuthSession
	if err == nil {
		result := database.First(&session, "time_hash=?", sessionId)
		if result.RowsAffected == 1 { //找到了信息
			username = session.Username
			sTime := time.Now().Format("2006-01-02 15:04:05")
			database.Model(&session).Update("last_visit", sTime) //用一次session，更新一次时间。
		}
	}
	if strings.HasPrefix(username, "[out]") || username == "" {
		username = "未登录"
	}
	var user User
	result := database.First(&user, "contact_tel = ?", username)
	if result.RowsAffected == 1 {
		group = "普通用户"
	} else {
		var repairman Repairman
		result := database.First(&repairman, "number = ?", username)
		if result.RowsAffected == 1 {
			group = "维修员"
		} else {
			var salesman Salesman
			result := database.First(&salesman, "number = ?", username)
			if result.RowsAffected == 1 {
				group = "业务员"
			} else {
				group = "未登录"
			}
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"username": username, "group": group})
}

//-----------------------------登录与登出---------------------------------------------
func startLogout(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	c.HTML(http.StatusOK, "logout.html", gin.H{"username": username, "group": group})
}
func logout(c *gin.Context) {
	sessionId, _ := c.Cookie("sessionId")
	var username string
	var session AuthSession
	result := database.First(&session, "time_hash=?", sessionId)
	if result.RowsAffected == 1 { //找到了信息
		username = session.Username
		username1 := "[out]" + username
		database.Model(&session).Update("username", username1) //在session表中将用户的账号前加入[out]标识
		c.SetCookie("sessionId", "", 0, "", "", false, true)   //清除浏览器中的cookie
	}
	c.HTML(http.StatusOK, "success.html", gin.H{"data": "用户" + username + "退出登录成功！", "website": "/index", "webName": "主页"})
}

func CheckPasswordHash(password, hash string) bool { //bcrypt比对密钥函数，用以比对两个字符串中的明文是否一致
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func startRegister(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err == nil { //已登录
		var session AuthSession
		result := database.First(&session, "time_hash=?", sessionId)
		if result.RowsAffected == 1 { //找到了信息
			c.HTML(http.StatusForbidden, "error.html", gin.H{"errdata": "已登录，请勿在登录状态注册账号", "website": "/index", "webName": "主页"})
			return
		}
	}
	c.HTML(http.StatusOK, "register.html", nil)
}

func register(c *gin.Context) {
	phoneNumber := c.PostForm("phone_number")
	password := c.PostForm("password")
	value := c.PostForm("ver_code")
	if CaptchaVerify(c, value) { //验证码通过
		var user User
		number := database.Find(&user).RowsAffected + 1
		strNumber := fmt.Sprintf("%08d", number)
		secretPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		data := User{strNumber, string(secretPassword), "", "", 100, "", phoneNumber}
		err := database.Create(&data)
		strErr := fmt.Sprintf("%v", err.Error)
		if strErr != "<nil>" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{"errdata": "注册失败！" + strErr, "website": "/register", "webName": "注册页面"})
		} else {
			c.HTML(http.StatusOK, "success.html", gin.H{"data": "注册成功！", "website": "/login", "webName": "登录页面"})
		}
	} else { //验证码错误
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"errdata": "验证码错误！", "website": "/login", "webName": "登录页面"})
	}
}

func startLogin(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err == nil { //已登录
		var session AuthSession
		result := database.First(&session, "time_hash=?", sessionId)
		if result.RowsAffected == 1 { //找到了信息
			c.HTML(http.StatusForbidden, "error.html", gin.H{"errdata": "已登录，请勿重复登录", "website": "/index", "webName": "主页"})
			return
		}
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	value := c.PostForm("ver_code")
	if CaptchaVerify(c, value) { //验证码正确
		var user User
		result := database.First(&user, "contact_tel=?", username)
		if result.RowsAffected == 1 { //只找到一条数据，用户名存在，比对密码
			if CheckPasswordHash(password, user.Password) { //密码比对通过
				goto CORRECT
			} else { //密码比对不通过
				c.HTML(http.StatusBadRequest, "error.html", gin.H{"errdata": "登录失败！", "website": "/login", "webName": "登录页面"})
			}
		} else { //不是用户表里的，找业务员表和维修员表
			var repairMan Repairman
			result := database.First(&repairMan, "number=?", username)
			if result.RowsAffected == 1 { //只找到一条数据，用户名存在，比对密码
				if CheckPasswordHash(password, repairMan.Password) { //密码比对通过
					goto CORRECT
				} else { //密码比对不通过
					c.HTML(http.StatusBadRequest, "error.html", gin.H{"errdata": "登录失败！", "website": "/login", "webName": "登录页面"})
				}
			} else { //不是用户表和维修员表里的，找业务员表
				var salesman Salesman
				result := database.First(&salesman, "number=?", username)
				if result.RowsAffected == 1 { //只找到一条数据，用户名存在，比对密码
					if CheckPasswordHash(password, salesman.Password) { //密码比对通过
						goto CORRECT
					} else { //密码比对不通过
						c.HTML(http.StatusBadRequest, "error.html", gin.H{"errdata": "登录失败！", "website": "/login", "webName": "登录页面"})
					}
				}
			}
		}
	CORRECT:
		{
			var session []AuthSession //寻找在数据库里是否有失效的cookie（指的是用户由于关闭了浏览器，已经没有sessionId了）
			result := database.Find(&session, "username=? or username=?", username, "[out]"+username)
			if result.RowsAffected != 0 { //本地没有cookie，需要申请新的cookie时，删除以前保留在服务器中的cookie
				database.Delete(&session)
			}

			curTime := time.Now()
			nanoCurTime := curTime.UnixNano() //获得当前时间（纳秒）
			sTime := curTime.Format("2006-01-02 15:04:05")
			strTime := fmt.Sprintf("%d", nanoCurTime) //时间变为字符串
			tempHash := sha256.Sum256([]byte(strTime))
			timeHash := hex.EncodeToString(tempHash[:]) //计算时间哈希
			data := AuthSession{timeHash, sTime, username}
			database.Create(&data) //在authSession数据库加入一条信息
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "sessionId",
				Value:    timeHash,
				Path:     "/",
				Domain:   "",
				SameSite: http.SameSiteLaxMode,
				Secure:   false,
				HttpOnly: true,
			})
			c.Redirect(http.StatusMovedPermanently, "/index")
		}
	} else { //验证码错误
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"errdata": "验证码错误！", "website": "/login", "webName": "登录页面"})
	}
}
func authMiddleWare() gin.HandlerFunc { //检查登录状态
	return func(c *gin.Context) {
		sessionId, err := c.Cookie("sessionId")
		if err == nil {
			var session AuthSession
			result := database.First(&session, "time_hash=?", sessionId)
			if result.RowsAffected == 1 { //找到了信息
				sTime := time.Now().Format("2006-01-02 15:04:05")
				database.Model(&session).Update("last_visit", sTime) //用一次session，更新一次时间。
				c.Set("username", session.Username)
				c.Next()
				return
			}
		}
		// 返回错误
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"errdata": "未登录", "website": "/login", "webName": "登录页面"})
		c.Abort()
		return
	}
}

func checkPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username").(string)
		var user User
		result := database.First(&user, "contact_tel=?", username)
		if result.RowsAffected == 1 {
			c.Set("username", user.ContactTel)
			c.Set("group", "普通用户")
		} else {
			var repairman Repairman
			result := database.First(&repairman, "number=?", username)
			if result.RowsAffected == 1 {
				c.Set("username", repairman.Number)
				c.Set("group", "维修员")
			} else {
				var salesman Salesman
				result := database.First(&salesman, "number=?", username)
				if result.RowsAffected == 1 {
					c.Set("username", salesman.Number)
					c.Set("group", "业务员")
				} else {
					c.HTML(http.StatusBadRequest, "error.html", gin.H{"errdata": "非法用户！", "errcode": 0, "website": "/index", "webName": "主页"})
					c.Abort()
				}
			}
		}
	}
}

//-----------------------------验证码-------------------------------------------------------------
func Session(keyPairs string) gin.HandlerFunc {
	store := SessionConfig()
	return sessions.Sessions(keyPairs, store)
}
func SessionConfig() sessions.Store {
	sessionMaxAge := 600
	sessionSecret := "SHU"
	var store sessions.Store
	store = cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge: sessionMaxAge, //seconds
		Path:   "/",
	})
	return store
}

func Captcha(c *gin.Context, length ...int) {
	l := captcha.DefaultLen
	w, h := 107, 36
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		w = length[1]
	}
	if len(length) == 3 {
		h = length[2]
	}
	captchaId := captcha.NewLen(l)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	_ = session.Save()
	_ = Serve(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
}
func CaptchaVerify(c *gin.Context, code string) bool {
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
