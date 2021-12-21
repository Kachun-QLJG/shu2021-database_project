package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

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
	c.HTML(http.StatusOK, "index.html", nil)
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
	phoneNumber := c.PostForm("phone")
	password := c.PostForm("pswd")
	value := c.PostForm("ver")
	if CaptchaVerify(c, value) { //验证码通过
		var user User
		number := database.Find(&user).RowsAffected + 1
		strNumber := fmt.Sprintf("%08d", number)
		secretPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		data := User{strNumber, string(secretPassword), "", "", 100, "", phoneNumber}
		err := database.Create(&data)
		strErr := fmt.Sprintf("%v", err.Error)
		if strErr != "<nil>" {
			c.JSON(http.StatusOK, gin.H{"status": "失败", "data": strErr})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "注册成功！"})
		}
	} else { //验证码错误
		c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "验证码"})
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
	sessionId, err := c.Cookie("sessionId")
	if err == nil { //已登录
		var session AuthSession
		result := database.First(&session, "time_hash=?", sessionId)
		if result.RowsAffected == 1 { //找到了信息
			c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "已登录"})
			return
		}
	}
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
				c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "用户名或密码错误！"})
				return
			}
		} else { //不是用户表里的，找业务员表和维修员表
			var repairMan Repairman
			result := database.First(&repairMan, "number=?", username)
			if result.RowsAffected == 1 { //只找到一条数据，用户名存在，比对密码
				if CheckPasswordHash(password, repairMan.Password) { //密码比对通过
					goto CORRECT
				} else { //密码比对不通过
					fmt.Println("@2")
					c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "用户名或密码错误！"})
					return
				}
			} else { //不是用户表和维修员表里的，找业务员表
				var salesman Salesman
				result := database.First(&salesman, "number=?", username)
				if result.RowsAffected == 1 { //只找到一条数据，用户名存在，比对密码
					if CheckPasswordHash(password, salesman.Password) { //密码比对通过
						goto CORRECT
					} else { //密码比对不通过
						fmt.Println()
						c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "用户名或密码错误！"})
						return
					}
				} else {
					fmt.Println("4")
					c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "用户名或密码错误！"})
					return
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
			c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "/index"})
		}
	} else { //验证码错误
		c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "验证码"})
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
