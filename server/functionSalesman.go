package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func saveAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.PostForm("attorney_no")
	repairType := c.PostForm("repair_type")
	repairClassification := c.PostForm("repair_classification")
	specificProblem := c.PostForm("specific_problem")
	outRange := c.PostForm("out_range")
	predictFinishTime := c.PostForm("predict_finish_time")
	endPetrol, _ := strconv.ParseFloat(c.PostForm("end_petrol"), 64)
	endMile, _ := strconv.ParseFloat(c.PostForm("end_mile"), 64)
	var attorney Attorney
	result := database.First(&attorney, "number = ? and salesman_id = ?", attorneyNo, username)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	database.Model(&attorney).Update("repair_type", repairType)
	database.Model(&attorney).Update("classification", repairClassification)
	database.Model(&attorney).Update("predict_finish_time", predictFinishTime)
	database.Model(&attorney).Update("specific_problem", specificProblem)
	database.Model(&attorney).Update("end_petrol", endPetrol)
	database.Model(&attorney).Update("end_mile", endMile)
	database.Model(&attorney).Update("out_range", outRange)
	c.String(http.StatusOK, "保存成功！")
}

// 获取项目详情
func getArrangements(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.Query("attorney_no")
	var attorney Attorney
	result := database.First(&attorney, "number = ? and salesman_id = ?", attorneyNo, username)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	type repairman struct {
		Type   string
		Number string
		Name   string
	}
	var arrangement []struct {
		ProjectNumber   string
		ProjectName     string
		ProjectTime     string
		ProjectProgress string
		Repairman       []repairman
	}
	timeType := ""
	var vehicle Vehicle
	database.First(&vehicle, "number = ?", attorney.VehicleNumber)
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
	database.Raw("select distinct arrangement.project_number as project_number, project_name as project_name, "+timeType+" as project_time, progress as project_progress\n"+
		"from arrangement inner join time_overview on arrangement.project_number = time_overview.project_number\n"+
		"where order_number = ?", attorneyNo).Scan(&arrangement)
	for i := range arrangement {
		database.Raw("select type as type, name as name, repairman.number as number\n"+
			"from arrangement inner join repairman on arrangement.repairman_number = repairman.number\n"+
			"where order_number = ? and project_number = ?", attorneyNo, arrangement[i].ProjectNumber).Scan(&arrangement[i].Repairman)
	}
	c.JSON(http.StatusOK, arrangement)
}

func changeDiscountRate(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.PostForm("attorney_no")
	client := c.PostForm("user_id")
	fmt.Println(client)
	fmt.Println(attorneyNo)
	discountRate, _ := strconv.Atoi(c.PostForm("discount_rate"))
	var attorney Attorney
	result := database.First(&attorney, "number = ? and salesman_id = ? and user_id = ?", attorneyNo, username, client)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	var user User
	database.First(&user, "number = ?", client)
	curDiscountRate := user.DiscountRate
	if curDiscountRate <= discountRate || curDiscountRate-discountRate > 10 || discountRate < 70 {
		c.String(http.StatusForbidden, "折扣率设置不合法！")
		return
	}
	database.Model(&user).Update("discount_rate", discountRate)
	c.String(http.StatusOK, "成功")
}

func getFullAttorneyS(c *gin.Context) {
	attorneyNo := c.Query("attorney_no")
	var attorney Attorney
	database.First(&attorney, "number = ?", attorneyNo)
	var user User
	database.First(&user, "number = ?", attorney.UserID)
	var vehicle Vehicle
	database.First(&vehicle, "number = ?", attorney.VehicleNumber)
	//定义表的结构
	type repairmanInfo struct {
		RepairmanNumber string
		RepairmanName   string
		Type            string
		Progress        string
	}
	type project struct {
		ProjectNumber string
		ProjectName   string
		ProjectTime   float64
		Repairman     []repairmanInfo
		ProjectRemark string
	}
	var result struct {
		VehicleVin           string
		StartPetrol          float64
		StartMile            float64
		EndPetrol            float64
		EndMile              float64
		PayMethod            string
		DiscountRate         int
		StartTime            string
		PredictFinishTime    string
		RepairType           string
		RepairClassification string
		RoughProblem         string
		OutRange             string
		SpecificProblem      string
		Project              []project
		Remark               string
		Progress             string
	}
	//填充信息
	result.VehicleVin = vehicle.Number
	result.StartPetrol = attorney.StartPetrol
	result.StartMile = attorney.StartMile
	result.EndPetrol = attorney.EndPetrol
	result.EndMile = attorney.EndMile
	result.PayMethod = attorney.PayMethod
	result.DiscountRate = user.DiscountRate
	result.StartTime = attorney.StartTime
	result.PredictFinishTime = attorney.PredictFinishTime
	result.RepairType = attorney.RepairType
	result.RepairClassification = attorney.Classification
	result.RoughProblem = attorney.RoughProblem
	result.OutRange = attorney.OutRange
	result.SpecificProblem = attorney.SpecificProblem
	result.Progress = attorney.Progress
	timeType := ""
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
	database.Raw("select time_overview.project_number as project_number, project_name as project_name, "+timeType+" as project_time, remark as project_remark\n"+
		"from arrangement inner join time_overview on arrangement.project_number = time_overview.project_number\n"+
		"where order_number = ?", attorneyNo).Scan(&result.Project) //找到该委托所有的维修项目
	var remarks []string            //备注字典
	for i := range result.Project { //遍历所有维修项目
		remarkList := strings.Split(result.Project[i].ProjectRemark, "*") //按照*分割备注
		result.Project[i].ProjectRemark = ""                              //处理后的备注，先置空
		if len(remarkList) > 1 {                                          //备注不为空
			for j := 1; j < len(remarkList); j++ { //0是空，1开始有值
				if j == len(remarkList)-1 { //如果是最后一条备注
					result.Project[i].ProjectRemark += "*" + remarkList[j] //不加换行符
				} else {
					result.Project[i].ProjectRemark += "*" + remarkList[j] + "\n" //加换行符
				}
				remarks = append(remarks, "*"+remarkList[j]) //将这条备注添加到字典中
			}
		}
		database.Raw("select repairman_number as repairman_number, name as repairman_name, type as type, progress as progress\n"+
			"from arrangement inner join repairman on arrangement.repairman_number = repairman.number\n"+
			"where order_number = ? and project_number = ?", attorneyNo, result.Project[i].ProjectNumber).Scan(&result.Project[i].Repairman) //找到该委托所有的维修项目
	}
	strRemark := ""        //备注返回的结果
	if len(remarks) != 0 { //如果备注字典中非空
		sort.Strings(remarks)           //按字典顺序排序
		for _, value := range remarks { //遍历字典
			var remarkDatabase Remark
			database.First(&remarkDatabase, "remark_number = ?", value)                    //查找备注表
			strRemark += remarkDatabase.RemarkNumber + " " + remarkDatabase.Content + "\n" //添加到返回结果中
		}
		strRemark = strRemark[:len(strRemark)-1] //处理掉最后一个换行符
	}
	result.Remark = strRemark //将备注结果加到返回结果中
	c.JSON(http.StatusOK, result)
}

// 订单添加项目
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
	database.First(&project, "project_number = ?", projectNo)
	var notification Notification
	notificationNumber := database.Find(&notification).RowsAffected + 1
	strNumber := fmt.Sprintf("%08d", notificationNumber) //获取通知序号
	sTime := time.Now().Format("2006-01-02 15:04:05")
	notice := Notification{
		strNumber,
		repairmanNo,
		"【通知】您有新的任务，请查收！",
		"维修员" + repairman.Name + "您好，您有新的任务：" + project.ProjectName + "，工时定额为" + c.PostForm("work_hour") + "工时，请及时处理！\n祝您工作愉快！",
		"未读",
		sTime,
	}
	database.Create(&notice) //添加到通知表
	c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "添加成功"})
}

// 通过维修工编号获取维修工信息
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
		Number       string
		Plate        string
		Vin          string
		CarModel     string
		CarType      string
		RoughProblem string
	}
	type processing struct {
		Number          string
		Plate           string
		Vin             string
		CarModel        string
		CarType         string
		RoughProblem    string
		SpecificProblem string
		Progress        string
	}
	type finished struct {
		Number          string
		Plate           string
		Vin             string
		CarModel        string
		CarType         string
		SpecificProblem string
		ClientPhone     string
	}
	var attorney struct {
		Pending    []pending
		Processing []processing
		Finished   []finished
	}
	database.Raw("select attorney.number as number, license_number as plate, vehicle.number as vin, model as car_model, type as car_type, rough_problem as rough_problem\n" +
		"from attorney inner join vehicle on attorney.vehicle_number = vehicle.number\n" +
		"where progress = '待处理' order by attorney.number").Scan(&attorney.Pending)
	database.Raw("select attorney.number as number, license_number as plate, vehicle.number as vin, model as car_model, type as car_type, rough_problem as rough_problem, specific_problem as specific_problem, progress as progress\n"+
		"from attorney inner join vehicle on attorney.vehicle_number = vehicle.number\n"+
		"where salesman_id = ? and progress != '已完成' order by progress desc, attorney.number", username).Scan(&attorney.Processing)
	database.Raw("select attorney.number as number, license_number as plate, vehicle.number as vin, model as car_model, type as car_type, specific_problem as specific_problem, contact_tel as client_phone\n"+
		"from attorney inner join vehicle inner join user on attorney.vehicle_number = vehicle.number and user.number = attorney.user_id\n"+
		"where salesman_id = ? and progress = '已完成' order by attorney.number desc", username).Scan(&attorney.Finished)
	c.JSON(http.StatusOK, attorney)
}

func receiveAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.PostForm("attorney_no")
	var attorney Attorney
	result := database.Find(&attorney, "number = ? and progress = '待处理'", attorneyNo)
	if result.RowsAffected != 1 {
		c.JSON(http.StatusOK, gin.H{"status": "错误", "data": "订单号错误或该订单已被其他用户接单！"})
		return
	}
	database.Model(&attorney).Update("progress", "处理中")
	database.Model(&attorney).Update("salesman_id", username)
	c.JSON(http.StatusOK, gin.H{"status": "成功", "data": "操作成功！"})
}

func setAttorneyFinished(c *gin.Context) {
	username := c.MustGet("username").(string)
	attorneyNo := c.PostForm("attorney_no")
	//endPetrol := c.PostForm("end_petrol")
	//endMile := c.PostForm("end_mile")
	var attorney Attorney
	result := database.First(&attorney, "number = ? and salesman_id = ?", attorneyNo, username)
	if result.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	sTime := time.Now().Format("2006-01-02 15:04:05")
	database.Model(&attorney).Update("actual_finish_time", sTime)
	//database.Model(&attorney).Update("end_petrol", endPetrol)
	//database.Model(&attorney).Update("end_mile", endMile)
	database.Model(&attorney).Update("progress", "已完成")
	genPdf(attorney.UserID, attorneyNo)
	var user User
	database.First(&user, "number = ?", attorney.UserID)
	c.JSON(http.StatusOK, gin.H{"url": "/show_pdf?attorney_no=" + attorneyNo, "discountRate": user.DiscountRate, "user_id": user.Number, "attorney_no": attorneyNo})
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
		database.Limit(20).Table("time_overview").Select("project_name as name, project_number as id, "+dbType+" as time").Limit(20).Where("project_spelling LIKE ? and ? != ''", searchText, dbType).Scan(&timeOverview)
	} else {
		database.Limit(20).Table("time_overview").Select("project_name as name, project_number as id, "+dbType+" as time").Limit(20).Where("project_name LIKE ? and ? != ''", searchText, dbType).Scan(&timeOverview)
	}
	c.JSON(http.StatusOK, timeOverview)
}
