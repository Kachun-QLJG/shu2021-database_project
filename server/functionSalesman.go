package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func addProjectForAttorney(c *gin.Context) {
	attorneyNo := c.PostForm("attorney_no")
	projectNo := c.PostForm("project_no")
	repairmanNo := c.PostForm("repairman_no")
	workHour, _ := strconv.ParseFloat(c.PostForm("work_hour"), 64)
	data := Arrangement{attorneyNo, projectNo, repairmanNo, "待确认"}
	addArrangementResult := database.Create(&data) //添加到派工单表
	if addArrangementResult.Error != nil {
		c.JSON(http.StatusOK, gin.H{"status": "错误", "data": addArrangementResult.Error})
		return
	}
	var repairman Repairman
	database.First(&repairman, "number = ?", repairmanNo)
	database.Model(&repairman).Update("current_work_hour", repairman.CurrentWorkHour+workHour)
	var project TimeOverview
	database.First(&project, "number = ?", projectNo)
	var notification Notification
	notificationNumber := database.Find(&notification).RowsAffected + 1
	strNumber := fmt.Sprintf("%08d", notificationNumber) //获取通知序号
	sTime := time.Now().Format("2006-01-02 15:04:05")
	notice := Notification{
		strNumber,
		repairmanNo,
		"【通知】您有新的任务，请查收！",
		"维修员" + repairman.Name + "您好，您有新的任务：" + project.ProjectName + "，工时定额为" + c.PostForm("work_hour") + "小时，请及时处理！\n祝您工作愉快！",
		"未读",
		sTime,
	}
	database.Create(&notice) //添加到通知表
}

func getCorrespondingRepairman(c *gin.Context) {
	repairmanType := c.Query("type")
	var result []struct {
		Number          string
		Name            string
		CurrentWorkHour float64
	}
	database.Table("repairman").Order("current_work_hour").Find(&result, "type = ? and status = '正常'", repairmanType)
	c.JSON(http.StatusOK, result)
}

func getRelatingAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	type pending struct {
		Number  string
		CarType string
		Status  string
	}
	type own struct {
		Number  string
		CarType string
		Status  string
	}
	var attorney struct {
		Pending []pending
		Own     []own
	}
	database.Raw("select attorney.number as number, type as car_type, progress as status from attorney inner join vehicle on attorney.vehicle_number = vehicle.number where progress = '待处理'").Scan(&attorney.Pending)
	database.Raw("select attorney.number as number, type as car_type, progress as status from attorney inner join vehicle on attorney.vehicle_number = vehicle.number where salesman_id = ?", username).Scan(&attorney.Own)
	c.JSON(http.StatusOK, attorney)
}

func getFinishedAttorneyS(c *gin.Context) {
	username := c.MustGet("username").(string)
	var result []struct {
		OrderNumber       string
		Vin               string
		RoughProblem      string
		SpecificProblem   string
		PredictFinishTime string
		Progress          string
		Username          string
	}
	database.Raw("select attorney.number as order_number, vehicle_number as vin, rough_problem as rough_problem, specific_problem as specific_problem, predict_finish_time as predict_finish_time, progress as progress, contact_tel as username\n"+
		"from attorney inner join user on attorney.user_id = user.number\n"+
		"where salesman_id = ? and progress = '已完成'", username).Scan(&result)
	c.JSON(http.StatusOK, result)
}

func setAttorneyFinished(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.PostForm("attorney_no")
	endPetrol := c.PostForm("end_petrol")
	endMile := c.PostForm("end_mile")
	var attorney Attorney
	result := database.First(&attorney, "number = ? and salesman_id = ?", attorneyNo, username)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	sTime := time.Now().Format("2006-01-02 15:04:05")
	database.Model(&attorney).Update("actual_finish_time", sTime)
	database.Model(&attorney).Update("end_petrol", endPetrol)
	database.Model(&attorney).Update("end_mile", endMile)
	database.Model(&attorney).Update("progress", "已完成")
	genPdf(attorney.UserID, attorneyNo)
	c.JSON(http.StatusOK, gin.H{"url": "/show_pdf?attorney_no=" + attorneyNo + "&user_id=" + attorney.UserID})
}

func getSalesmanInfo(c *gin.Context) {
	username := c.MustGet("username").(string)
	var salesman struct {
		Number string
		Name   string
	}
	database.Table("salesman").Where("number = ?", username).Select("number as number, name as name").Scan(&salesman)
	c.JSON(http.StatusOK, salesman)
}

func startSCheckOrders(c *gin.Context) {
	c.HTML(http.StatusOK, "salesman_check_orders.html", nil)
}

func startTakeOrders(c *gin.Context) {
	c.HTML(http.StatusOK, "salesman_take_orders.html", nil)
}
