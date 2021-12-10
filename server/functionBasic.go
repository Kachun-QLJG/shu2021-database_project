package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func test(c *gin.Context) {
	text := c.Query("text")
	searchText := "%" + strings.ToLower(text) + "%"
	var typeOverview []TypeOverview
	if searchText[1] >= 'a' && searchText[1] <= 'z' {
		database.Limit(10).Where("project_spelling LIKE ?", searchText).Find(&typeOverview)
	} else {
		database.Limit(10).Where("project_name LIKE ?", searchText).Find(&typeOverview)
	}
	var result [10]struct {
		Name string
		Id   string
	}
	count := 0
	for row := range typeOverview {
		result[count].Id = typeOverview[row].ProjectNumber
		result[count].Name = typeOverview[row].ProjectName
		count++
	}
	c.JSON(http.StatusOK, result)
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
		c.HTML(http.StatusOK, "customer_index.html", gin.H{"username": username, "group": group})
	} else {
		var repairman Repairman
		result := database.First(&repairman, "number = ?", username)
		if result.RowsAffected == 1 {
			group = "维修员"
			c.HTML(http.StatusOK, "repairman_index.html", gin.H{"username": username, "group": group})
		} else {
			var salesman Salesman
			result := database.First(&salesman, "number = ?", username)
			if result.RowsAffected == 1 {
				group = "业务员"
				c.HTML(http.StatusOK, "salesman_index.html", gin.H{"username": username, "group": group})
			} else {
				group = "未登录"
				c.HTML(http.StatusOK, "index.html", gin.H{"username": username, "group": group})
			}
		}
	}
	//c.HTML(http.StatusOK, "index.html", gin.H{"username": username, "group": group})
}
