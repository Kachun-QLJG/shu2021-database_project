package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func checkRegister(c *gin.Context) {
	phoneNumber := c.Query("contact_tel")
	var user User
	c.String(http.StatusOK, strconv.FormatInt(database.First(&user, "contact_tel = ?", phoneNumber).RowsAffected, 10))
}

func generateLicensePlate(licenseNumber string) {
	color := "blue"
	count := 0
	for _, char := range licenseNumber {
		if count == 0 && char >= 'A' && char <= 'Z' {
			color = "white"
		}
		if count == 0 && string(char) == "使" {
			color = "black"
		}
		if count == 6 && string(char) == "警" {
			color = "white"
		}
		if count == 6 && string(char) == "领" || string(char) == "港" {
			color = "black"
		}
		if count == 6 && string(char) == "学" {
			color = "yellow"
		}
		count++
	}
	if count == 8 {
		color = "green_car"
	}
	args := []string{"--plate-number", licenseNumber, "--bg-color", color}
	cmd := exec.Command("generate_special_plate.exe", args...)
	fmt.Println(cmd.Args)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}
func checkVehicle(c *gin.Context) {
	carNumber := c.Query("number")
	var vehicle Vehicle
	c.String(http.StatusOK, strconv.FormatInt(database.First(&vehicle, "number = ?", carNumber).RowsAffected, 10))
}

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
	number := c.PostForm("number")                //获取车架号
	licenseNumber := c.PostForm("license_number") //获取车牌号
	phoneNumber := c.MustGet("username").(string) //获取用户名
	color := c.PostForm("color")
	model := c.PostForm("model")
	carType := c.PostForm("type")
	sTime := time.Now().Format("2006-01-02 15:04:05")
	generateLicensePlate(licenseNumber)
	var user User
	database.First(&user, "contact_tel = ?", phoneNumber) //查找用户
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
			"尊敬的用户" + oldUser.Name + "您好，您的车辆" + vehicle.LicenseNumber + "已被手机尾号为" + phoneNumber[7:] + "的用户绑定。",
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

func showPlate(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	licenseNumber := c.Query("license_number")
	if group != "普通用户" {
		c.JSON(http.StatusForbidden, gin.H{"data": "错误!"})
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", number)
	var vehicle Vehicle
	database.First(&vehicle, "license_number = ?", licenseNumber).Order("time desc")
	if vehicle.UserID == user.Number {
		c.File("./html/statics/license plates/" + licenseNumber + ".jpg")
		return
	} else {
		c.String(http.StatusForbidden, "无权查看车牌")
		return
	}
}

func getVehicle(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	page := c.Query("pg")
	numPage, _ := strconv.Atoi(page)
	if group != "普通用户" {
		c.JSON(http.StatusForbidden, gin.H{"data": "错误!"})
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", number)
	var vehicle []Vehicle
	database.Limit(3).Offset(3*(numPage-1)).Find(&vehicle, "user_id = ?", user.Number)
	fmt.Println(vehicle)
	var result [4]struct {
		Number        string
		LicenseNumber string
		Color         string
		Model         string
		Type          string
	}
	count := 0
	for row := range vehicle {
		result[count].Number = vehicle[row].Number
		result[count].LicenseNumber = vehicle[row].LicenseNumber
		result[count].Color = vehicle[row].Color
		result[count].Model = vehicle[row].Model
		result[count].Type = vehicle[row].Type
		count++
	}
	c.JSON(http.StatusOK, result)
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
