package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func getUsername(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.String(http.StatusOK, username)
}

func read(c *gin.Context) {
	username := c.MustGet("username").(string)
	var notification Notification
	database.First(&notification, "user_id = ? and status = ?", username, "未读")
	database.Model(&notification).Update("status", "已读") //更改消息为已读
	c.String(http.StatusOK, "成功")
}

func checkNotification(c *gin.Context) {
	username := c.MustGet("username").(string)
	var notification Notification
	result := database.First(&notification, "user_id = ? and status = ?", username, "未读")
	if result.RowsAffected == 1 {
		c.JSON(http.StatusOK, gin.H{"title": notification.Title, "content": notification.Content, "time": notification.Time})
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func changePassword(c *gin.Context) {
	oldPassword := c.PostForm("oldpswd")
	newPassword := c.PostForm("pswd")
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	var passwordChange struct {
		OldPassword string
		Username    string
	}
	if group == "普通用户" {
		database.Table("user").Select("contact_tel as username, password as old_password").Where("contact_tel = ?", username).Scan(&passwordChange).Limit(1)
	} else if group == "业务员" {
		database.Table("salesman").Select("number as username, password as old_password").Where("number = ?", username).Scan(&passwordChange).Limit(1)
	} else {
		database.Table("repairman").Select("number as username, password as old_password").Where("number = ?", username).Scan(&passwordChange).Limit(1)
	}
	if CheckPasswordHash(oldPassword, passwordChange.OldPassword) { //密码比对通过
		secretPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if group == "普通用户" {
			var user User
			database.First(&user, "contact_tel = ?", username)
			database.Model(&user).Update("password", secretPassword)
		} else if group == "业务员" {
			var salesman Salesman
			database.First(&salesman, "number = ?", username)
			database.Model(&salesman).Update("password", secretPassword)
		} else {
			var repairman Repairman
			database.First(&repairman, "number = ?", username)
			database.Model(&repairman).Update("password", secretPassword)
		}
		username1 := "[out]" + username
		var session AuthSession
		database.First(&session, "username = ?", username)
		database.Model(&session).Update("username", username1) //在session表中将用户的账号前加入[out]标识
		c.SetCookie("sessionId", "", 0, "", "", false, true)   //清除浏览器中的cookie
		c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "/index"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "密码错误！"})
	}
}

func welcome(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	var username string
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
		c.HTML(http.StatusOK, "user_index.html", nil)
	} else {
		var repairman Repairman
		result = database.First(&repairman, "number = ?", username)
		if result.RowsAffected == 1 {
			c.HTML(http.StatusOK, "repairman_index.html", nil)
		} else {
			var salesman Salesman
			result = database.First(&salesman, "number = ?", username)
			if result.RowsAffected == 1 {
				c.HTML(http.StatusOK, "salesman_index.html", nil)
			} else {
				c.HTML(http.StatusOK, "login.html", nil)
			}
		}
	}
}
