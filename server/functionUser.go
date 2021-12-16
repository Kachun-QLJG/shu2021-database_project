package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
func checkVehicle(c *gin.Context) {
	carNumber := c.Query("number")
}*/
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
	number := c.PostForm("number")                 //获取车架号
	licenseNumber := c.PostForm("license_number")  //获取车牌号
	phone_number := c.MustGet("username").(string) //获取用户名
	color := c.PostForm("color")
	model := c.PostForm("model")
	carType := c.PostForm("type")
	sTime := time.Now().Format("2006-01-02 15:04:05")
	var user User
	database.First(&user, "contact_tel = ?", phone_number) //查找用户
	var vehicle Vehicle
	res := database.First(&vehicle, "number = ?", number) //根据输入的车架号，查这辆车是否已被绑定
	if res.RowsAffected == 0 {                            //如果没有被绑定
		data := Vehicle{number, licenseNumber, user.Number, color, model, carType, sTime} //新建元组
		err := database.Create(&data)                                                     //添加到数据库
		strErr := fmt.Sprintf("%v", err.Error)                                            //将错误（若有）转换成字符串
		if strErr != "<nil>" {
			c.JSON(http.StatusOK, gin.H{"status": "失败", "data": strErr})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "成功"})
		}
	} else { //这辆车已经被绑定
		//通知原车主
		var notification Notification
		number := database.Find(&notification).RowsAffected + 1
		strNumber := fmt.Sprintf("%08d", number) //获取通知序号
		var oldUser User
		database.First(&oldUser, "number = ?", vehicle.UserID) //找到原车主信息
		data := Notification{
			strNumber,
			oldUser.ContactTel,
			"【通知】您的车辆" + vehicle.LicenseNumber + "已被他人绑定",
			"尊敬的用户" + oldUser.Name + "您好，您的车辆" + vehicle.LicenseNumber + "已被手机尾号为" + phone_number[7:] + "的用户绑定。",
			"未读",
			sTime,
		}
		database.Create(&data) //添加到通知表
		//更新
		database.Model(&vehicle).Update("license_number", licenseNumber) //更改车牌号为传过来的车牌号
		database.Model(&vehicle).Update("user_id", user.Number)          //更改用户id为传过来的用户id
		database.Model(&vehicle).Update("color", color)                  //更改颜色为传过来的颜色
		database.Model(&vehicle).Update("model", model)                  //更改车型为传过来的车型
		database.Model(&vehicle).Update("car_type", carType)             //更改类别为传过来的类别
		database.Model(&vehicle).Update("time", sTime)                   //更改时间为传过来的时间
		c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "成功"})
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
