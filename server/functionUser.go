package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func startAddVehicle(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	c.HTML(http.StatusOK, "user_car_register.html", gin.H{"username": number, "group": group})
}

func startUCheckOrders(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	c.HTML(http.StatusOK, "user_check_orders.html", gin.H{"username": number, "group": group})
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

func userinfo(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", number)
	c.JSON(http.StatusOK, gin.H{
		"number":         user.Number,
		"name":           user.Name,
		"property":       user.Property,
		"discount_rate":  user.DiscountRate,
		"contact_person": user.ContactPerson,
		"contact_tel":    user.ContactTel})
}

func changeUserinfo(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	name := c.PostForm("name")
	property := c.PostForm("property")
	contactPerson := c.PostForm("contact_person")
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", number)
	database.Model(&user).Update("name", name)                    //更改状态为传过来的状态
	database.Model(&user).Update("property", property)            //更改状态为传过来的状态
	database.Model(&user).Update("contact_person", contactPerson) //更改状态为传过来的状态
	c.String(http.StatusOK, "成功！")
}

func startChangeUserinfo(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	c.HTML(http.StatusOK, "user_change_profile.html", gin.H{"username": number, "group": group})
}
