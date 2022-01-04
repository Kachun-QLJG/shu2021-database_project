package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getProcessingAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误!")
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", username)
	type information struct {
		OrderNumber       string
		Vin               string
		RoughProblem      string
		SpecificProblem   string
		PredictFinishTime string
		Progress          string
	}
	var result []information
	database.Table("attorney").Select("number as order_number, vehicle_number as vin, rough_problem as rough_problem, specific_problem as specific_problem, predict_finish_time as predict_finish_time, progress as progress").Where("user_id = ? and progress != '已完成'", user.Number).Scan(&result) //查询该用户的，状态不是已完成的委托
	c.JSON(http.StatusOK, result)
}

func getFinishedAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误!")
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", username)
	type information struct {
		OrderNumber       string
		Vin               string
		RoughProblem      string
		SpecificProblem   string
		PredictFinishTime string
		Progress          string
		Owner             string
	}
	var result []information
	database.Table("attorney").Select("number as order_number, vehicle_number as vin, rough_problem as rough_problem, specific_problem as specific_problem, predict_finish_time as predict_finish_time, progress as progress").Where("user_id = ? and progress = '已完成'", user.Number).Scan(&result) //查询该用户的，状态是已完成的委托
	for index := range result {
		var vehicle Vehicle
		res := database.First(&vehicle, "number = ? and user_id = ?", result[index].Vin, user.Number)
		if res.RowsAffected == 1 {
			result[index].Owner = "是"
		} else {
			result[index].Owner = "否"
		}
	}
	c.JSON(http.StatusOK, result)
}

func getAttorneyDetail(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误!")
		return
	}
	attorneyNumber := c.Query("no")
	var user User
	database.First(&user, "contact_tel = ?", username)
	var attorney Attorney
	result := database.First(&attorney, "number = ? and user_id = ?", attorneyNumber, user.Number)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限查看")
	} else {
		c.JSON(http.StatusOK, attorney)
	}
}

func createAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误!")
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", username)
	var attorney Attorney
	number := database.Find(&attorney).RowsAffected + 1
	strNumber := fmt.Sprintf("%08d", number)
	vin := c.PostForm("vin")
	payMethod := c.PostForm("pay_method")
	startTime := c.PostForm("start_time")
	roughProblem := c.PostForm("rough_problem")
	startPetrol, _ := strconv.ParseFloat(c.PostForm("start_petrol"), 64)
	startMile, _ := strconv.ParseFloat(c.PostForm("start_mile"), 64)
	type temp struct {
		Number        string
		UserID        string
		VehicleNumber string
		PayMethod     string
		StartTime     string
		RoughProblem  string
		Progress      string
		StartPetrol   float64
		StartMile     float64
	}
	data := temp{strNumber, user.Number, vin, payMethod, startTime, roughProblem, "待处理", startPetrol, startMile}
	database.Table("attorney").Create(&data)
	c.String(http.StatusOK, "成功")
}

func getRepairHistory(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	vin := c.Query("vin")
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误!")
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", number)
	var vehicle Vehicle
	result := database.First(&vehicle, "user_id = ? and number = ?", user.Number, vin)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "车辆不绑定于您！无权限查看！")
		return
	}
	var data []struct {
		Time    string
		Problem string
	}
	database.Order("time desc").Table("attorney").Select("start_time as time, specific_problem as problem").Where("vehicle_number = ?", vin).Scan(&data)
	c.JSON(http.StatusOK, data)
}

func getVehicle(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.JSON(http.StatusForbidden, gin.H{"data": "错误!"})
		return
	}
	var user User
	database.First(&user, "contact_tel = ?", number)
	var vehicle []Vehicle
	database.Order("time desc").Find(&vehicle, "user_id = ?", user.Number)
	c.JSON(http.StatusOK, vehicle)
}

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

func startUCheckOrdersOngoing(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	c.HTML(http.StatusOK, "user_check_orders_ongoing.html", gin.H{"username": number, "group": group})
}

func startUCheckOrdersFinished(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "普通用户" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	c.HTML(http.StatusOK, "user_check_orders_finished.html", gin.H{"username": number, "group": group})
}

func addVehicle(c *gin.Context) {
	number := c.PostForm("number")                 //获取车架号
	licenseNumber := c.PostForm("license_number")  //获取车牌号
	licenseNumber = strings.ToUpper(licenseNumber) //车牌号改大写
	phoneNumber := c.MustGet("username").(string)  //获取用户名
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
		//检查这辆车是否有正在处理的订单
		var attorney Attorney
		result := database.First(&attorney, "vehicle_number = ? and status != '已完成'", number)
		if result.RowsAffected == 0 {
			//通知原车主
			var notification Notification
			notificationNumber := database.Find(&notification).RowsAffected + 1
			strNumber := fmt.Sprintf("%08d", notificationNumber) //获取通知序号
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
			database.Model(&vehicle).Update("color", color)                  //更改颜色为传过来的颜色
			database.Model(&vehicle).Update("model", model)                  //更改车型为传过来的车型
			database.Model(&vehicle).Update("type", carType)                 //更改类别为传过来的类别
			database.Model(&vehicle).Update("time", sTime)                   //更改时间为传过来的时间
			database.Model(&vehicle).Update("user_id", user.Number)          //更改用户id为传过来的用户id
			c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "成功"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "该车辆有未完成的维修委托！"})
		}
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
	result := database.First(&vehicle, "license_number = ? and user_id = ?", licenseNumber, user.Number)
	if result.RowsAffected == 1 {
		c.File("./html/statics/license plates/" + licenseNumber + ".jpg")
		return
	} else {
		c.String(http.StatusForbidden, "无权查看车牌")
		return
	}
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
