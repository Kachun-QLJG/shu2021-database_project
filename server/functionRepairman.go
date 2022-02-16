package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getFinishedArrangement(c *gin.Context) {
	username := c.MustGet("username").(string)
	type parts struct {
		PartsNumber string
		PartsName   string
		PartsCount  string
	}
	var arrangement []struct {
		OrderNumber   string
		ProjectNumber string
		Parts         []parts
		Vin           string
		Status        string
	}
	database.Raw("select order_number as order_number, project_number as project_number, vehicle_number as vin, arrangement.progress as status\n"+
		"from arrangement inner join attorney on arrangement.order_number = attorney.number\n"+
		"where repairman_number = ? and arrangement.progress = '已完成'", username).Scan(&arrangement)
	for i := range arrangement {
		database.Raw("select repair_parts.parts_number as parts_number, parts_name as parts_name, parts_count as parts_count\n"+
			"from repair_parts inner join parts_overview on repair_parts.parts_number = parts_overview.parts_number\n"+
			"where order_number = ? and project_number = ?", arrangement[i].OrderNumber, arrangement[i].ProjectNumber).Scan(&arrangement[i].Parts)
	}
	c.JSON(http.StatusOK, arrangement)
}

func getProcessingArrangement(c *gin.Context) {
	username := c.MustGet("username").(string)
	type parts struct {
		PartsNumber string
		PartsName   string
		PartsCount  string
	}
	type project struct {
		ProjectNumber string
		ProjectName   string
		ProjectTime   string
		Parts         []parts
	}
	var arrangement []struct {
		OrderNumber     string
		Plate           string
		Vin             string
		Project         []project
		SpecificProblem string
	}
	database.Raw("select distinct order_number as order_number, license_number as plate, vehicle.number as vin, specific_problem as specific_problem\n"+
		"from arrangement inner join attorney inner join vehicle on arrangement.order_number = attorney.number and attorney.vehicle_number = vehicle.number\n"+
		"where repairman_number = ? and arrangement.progress = '维修中'", username).Scan(&arrangement)
	for i := range arrangement {
		timeType := ""
		var vehicle Vehicle
		database.First(&vehicle, "number = ?", arrangement[i].Vin)
		if vehicle.Type == "轿车-A" {
			timeType = "time_a"
		} else if vehicle.Type == "轿车-B" {
			timeType = "time_b"
		} else if vehicle.Type == "轿车-C" {
			timeType = "time_c"
		} else if vehicle.Type == "轿车-D" {
			timeType = "time_d"
		} else {
			timeType = "time_e"
		}
		database.Raw("select arrangement.project_number as project_number, project_name as project_name, "+timeType+" as project_time\n"+
			"from arrangement inner join time_overview on arrangement.project_number = time_overview.project_number\n"+
			"where repairman_number = ? and arrangement.progress = '维修中'", username).Scan(&arrangement[i].Project)
		for j := range arrangement[i].Project {
			database.Raw("select repair_parts.parts_number as parts_number, parts_name as parts_name, parts_count as parts_count\n"+
				"from repair_parts inner join parts_overview on repair_parts.parts_number = parts_overview.parts_number\n"+
				"where project_number = ?", arrangement[i].Project[j].ProjectNumber).Scan(&arrangement[i].Project[j].Parts)
		}
	}
	c.JSON(http.StatusOK, arrangement)
}

func getPendingArrangement(c *gin.Context) {
	username := c.MustGet("username").(string)
	var arrangement []struct {
		OrderNumber     string
		Plate           string
		Vin             string
		ProjectNumber   string
		ProjectName     string
		ProjectTime     string
		SpecificProblem string
		Salesman        string
	}
	var temp struct{ Time string }
	database.Raw("select order_number as order_number, license_number as plate, vehicle.number as vin, arrangement.project_number as project_number,\n"+
		"time_overview.project_name as project_name, specific_problem as specific_problem, salesman_id as salesman\n"+
		"from arrangement inner join attorney inner join vehicle inner join time_overview on arrangement.order_number = attorney.number\n"+
		"and attorney.vehicle_number = vehicle.number and arrangement.project_number = time_overview.project_number\n"+
		"where repairman_number = ? and arrangement.progress = '待确认'", username).Scan(&arrangement)
	for i := range arrangement {
		timeType := ""
		var vehicle Vehicle
		database.First(&vehicle, "number = ?", arrangement[i].Vin)
		if vehicle.Type == "轿车-A" {
			timeType = "time_a"
		} else if vehicle.Type == "轿车-B" {
			timeType = "time_b"
		} else if vehicle.Type == "轿车-C" {
			timeType = "time_c"
		} else if vehicle.Type == "轿车-D" {
			timeType = "time_d"
		} else {
			timeType = "time_e"
		}
		database.Raw("select "+timeType+" as time from time_overview where project_number = ?", arrangement[i].ProjectNumber).Scan(&temp)
		arrangement[i].ProjectTime = temp.Time
	}
	c.JSON(http.StatusOK, arrangement)
}

func changeRepairProgress(c *gin.Context) {
	attorneyNo := c.PostForm("attorney_no")
	projectNo := c.PostForm("project_no")
	username := c.MustGet("username").(string)
	progress := c.PostForm("progress")
	var arrangement Arrangement
	database.First(&arrangement, "order_number = ? and project_number = ? and repairman_number = ?", attorneyNo, projectNo, username)
	database.Model(&arrangement).Update("progress", progress)
	if progress == "已完成" {
		result := database.Find(&arrangement, "order_number = ? and progress != '已完成'", attorneyNo)
		if result.RowsAffected == 0 {
			var attorney Attorney
			database.First(&attorney, "number = ?", attorneyNo)
			database.Model(&arrangement).Update("progress", "待结算")
			database.Model(&arrangement).Update("actual_finish_time", time.Now().Format("2006-01-02"))
			var user User
			database.First(&user, "number = ?", attorney.UserID)
			var vehicle Vehicle
			database.First(&vehicle, "number = ?", attorney.VehicleNumber)
			var notification Notification
			notificationNumber := database.Find(&notification).RowsAffected + 1
			strNumber := fmt.Sprintf("%08d", notificationNumber) //获取通知序号
			data := Notification{                                //给用户发通知
				strNumber,
				user.ContactTel,
				"【通知】您的车辆" + vehicle.LicenseNumber + "维修完成",
				"尊敬的用户" + user.Name + "您好，您的车辆" + vehicle.LicenseNumber + "已完成订单号为" + attorney.Number + "的维修委托。请您保持电话畅通，与您对接的业务员将与您联系，以确认后续事宜。\n最后感谢您选择了我们为您提供车辆维修服务，祝您一路顺风、生活愉快。",
				"未读",
				time.Now().Format("2006-01-02 15:04:05"),
			}
			database.Create(&data) //添加到通知表
			data = Notification{   //给业务员发通知
				fmt.Sprintf("%08d", notificationNumber+1),
				attorney.SalesmanID,
				"【通知】订单" + attorney.Number + "维修完成，请及时联系客户",
				"业务员您好，您对接的" + attorney.Number + "号维修委托已完成维修，请您及时与客户" + user.ContactPerson + "取得联系（联系电话：" + user.ContactTel + "），提醒客户及时结算、取车。",
				"未读",
				time.Now().Format("2006-01-02 15:04:05"),
			}
			database.Create(&data) //添加到通知表
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": "成功"})
}

func getRepairmanInfo(c *gin.Context) {
	username := c.MustGet("username").(string)
	var repairman struct {
		Number          string
		Name            string
		Type            string
		CurrentWorkHour float64
		Status          string
	}
	database.Table("repairman").Where("number = ?", username).Select("number as number, name as name, type as type, current_work_hour as current_work_hour, status as status").Scan(&repairman)
	c.JSON(http.StatusOK, repairman)
}

func addPartsForProject(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.PostForm("attorney_no")
	projectNo := c.PostForm("project_no")
	partsNo := c.PostForm("parts_no")
	partsCount, _ := strconv.Atoi(c.PostForm("number"))
	var arrangement Arrangement
	result := database.First(&arrangement, "order_number = ? and project_number = ? and repairman_number = ?", attorneyNo, projectNo, username)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	var repairParts RepairParts
	partsResult := database.First(&repairParts, "order_number = ? and project_number = ? and parts_number = ?", attorneyNo, projectNo, partsNo)
	if partsResult.RowsAffected == 0 { //新建零件
		data := RepairParts{attorneyNo, projectNo, partsNo, partsCount}
		createPartsResult := database.Create(&data)
		if createPartsResult.Error != nil {
			c.JSON(http.StatusOK, gin.H{"status": "错误", "data": createPartsResult.Error})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "成功", "data": ""})
		}
	} else { //更改数量
		database.Model(&repairParts).Update("parts_count", partsCount)
		c.JSON(http.StatusOK, gin.H{"status": "成功", "data": ""})
	}
}

func searchForParts(c *gin.Context) {
	text := c.Query("text")
	text = strings.Replace(text, " ", "%", -1)
	searchText := "%" + strings.ToLower(text) + "%"
	var parts []PartsOverview
	if searchText[1] >= 'a' && searchText[1] <= 'z' || searchText[1] >= '0' && searchText[1] <= '9' {
		database.Limit(20).Where("parts_spelling LIKE ?", searchText).Find(&parts)
	} else {
		database.Limit(20).Where("parts_name LIKE ?", searchText).Find(&parts)
	}
	var result [20]struct {
		Name  string
		Id    string
		Price float64
	}
	count := 0
	for row := range parts {
		result[count].Id = parts[row].PartsNumber
		result[count].Name = parts[row].PartsName
		result[count].Price = parts[row].PartsCost
		count++
	}
	c.JSON(http.StatusOK, result)
}

func checkStatus(c *gin.Context) {
	number := c.MustGet("username").(string)
	var repairman Repairman
	database.First(&repairman, "number = ?", number)
	c.String(http.StatusOK, repairman.Status)
}

func changeStatus(c *gin.Context) {
	status := c.PostForm("status")
	number := c.MustGet("username").(string)
	fmt.Println("status: ", status)
	var repairman Repairman
	database.First(&repairman, "number = ?", number)
	database.Model(&repairman).Update("status", status) //更改状态为传过来的状态
	c.String(http.StatusOK, "修改成功！")
}
