package main

import (
	"fmt"
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
func getGroup(c *gin.Context) {
	group := c.MustGet("group").(string)
	c.String(http.StatusOK, group)
}

func searchForProjects(c *gin.Context) {
	text := c.Query("text")
	carType := c.Query("type")
	dbType := strings.ToLower(carType)
	dbType = "time_" + dbType
	text = strings.Replace(text, " ", "%", -1)
	searchText := "%" + strings.ToLower(text) + "%"
	var timeOverview []struct {
		Name string
		Id   string
		Time float64
	}
	if searchText[1] >= 'a' && searchText[1] <= 'z' || searchText[1] >= '0' && searchText[1] <= '9' {
		database.Table("time_overview").Select("project_name as name, project_number as id, "+dbType+" as time").Limit(20).Where("project_spelling LIKE ? and ? != ''", searchText, dbType).Scan(&timeOverview)
		//database.Limit(20).Where("project_spelling LIKE ? and time_? != ''" , searchText, dbType).Find(&timeOverview)
	} else {
		database.Table("time_overview").Select("project_name as name, project_number as id, "+dbType+" as time").Limit(20).Where("project_name LIKE ?and ? != ''", searchText, dbType).Scan(&timeOverview)
		//database.Limit(20).Where("project_name LIKE ?and time_? != ''" , searchText, dbType).Find(&timeOverview)
	}
	var result [20]struct {
		Name string
		Id   string
		Time float64
	}
	count := 0
	for row := range timeOverview {
		result[count].Id = timeOverview[row].Id
		result[count].Name = timeOverview[row].Name
		result[count].Time = timeOverview[row].Time
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
		c.JSON(http.StatusOK, gin.H{"title": notification.Title, "content": notification.Content, "time": notification.Time})
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
	oldPassword := c.PostForm("oldpswd")
	newPassword := c.PostForm("pswd")
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	fmt.Println(oldPassword, newPassword)
	if oldPassword == "" {
		fmt.Println("empty!")
	}
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
	fmt.Println(passwordChange)
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

//localhost:8080/add_vehicle?number=1&license_number=1&user_id=1&color=1&model=1&type=1

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
		result := database.First(&repairman, "number = ?", username)
		if result.RowsAffected == 1 {
			c.HTML(http.StatusOK, "repairman_index.html", nil)
		} else {
			var salesman Salesman
			result := database.First(&salesman, "number = ?", username)
			if result.RowsAffected == 1 {
				c.HTML(http.StatusOK, "salesman_index.html", nil)
			} else {
				c.HTML(http.StatusOK, "index.html", nil)
			}
		}
	}
	//c.HTML(http.StatusOK, "index.html", gin.H{"username": username, "group": group})
}
