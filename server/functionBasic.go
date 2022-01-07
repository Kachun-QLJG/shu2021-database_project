package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getFullAttorney(c *gin.Context) {
	number := c.MustGet("username").(string)
	attorneyNo := c.Query("attorney_no")
	var user User
	database.First(&user, "contact_tel = ?", number) //找到用户
	var attorney Attorney
	findAttorneyResult := database.First(&attorney, "number = ? and user_id = ? and progress != '已完成'", attorneyNo, user.Number) //找到该用户的该订单号的委托
	if findAttorneyResult.RowsAffected == 0 {
		c.String(http.StatusForbidden, "无权限！")
		return
	}
	var salesman Salesman
	database.First(&salesman, "number = ?", attorney.SalesmanID) //业务员信息
	var vehicle Vehicle
	database.First(&vehicle, "number = ? and user_id = ?", attorney.VehicleNumber, attorney.UserID) //车辆信息
	//定义表的结构
	type info struct {
		UserNumber           string
		UserContactPerson    string
		UserContactTel       string
		SalesmanNumber       string
		SalesmanName         string
		VehiclePlate         string
		VehicleVin           string
		VehicleModel         string
		VehicleType          string
		StartPetrol          float64
		StartMile            float64
		PayMethod            string
		DiscountRate         int
		PredictFinishTime    string
		RepairType           string
		RepairClassification string
		RoughProblem         string
		OutRange             string
		SpecificProblem      string
	}
	type parts []struct {
		PartsNumber      string
		PartsName        string
		PartsCount       int
		PartsSinglePrice float64
		PartsTotalPrice  float64
	}
	type project struct {
		ProjectNumber     string
		ProjectName       string
		ProjectTime       float64
		ProjectRemark     string
		ProjectPartsCount int
		ProjectParts      parts
	}
	type attorneyProject struct {
		Num            int
		Project        []project
		TotalWorkHour  float64
		TotalPartsCost float64
		Remark         string
	}
	type tail struct {
		WorkHourPrice int
		TotalPrice    string
		CalculateWay  string
	}
	var result struct { //返回的结构体
		Head            info
		AttorneyProject attorneyProject
		Tail            tail
	}
	//填充Head信息
	result.Head.UserNumber = user.Number
	result.Head.UserContactPerson = user.ContactPerson
	result.Head.UserContactTel = user.ContactTel
	result.Head.SalesmanNumber = salesman.Number
	result.Head.SalesmanName = salesman.Name
	result.Head.VehiclePlate = vehicle.LicenseNumber
	result.Head.VehicleVin = vehicle.Number
	result.Head.VehicleModel = vehicle.Model
	result.Head.VehicleType = vehicle.Type
	result.Head.StartPetrol = attorney.StartPetrol
	result.Head.StartMile = attorney.StartMile
	result.Head.PayMethod = attorney.PayMethod
	result.Head.DiscountRate = user.DiscountRate
	result.Head.PredictFinishTime = attorney.PredictFinishTime
	result.Head.RepairType = attorney.RepairType
	result.Head.RepairClassification = attorney.Classification
	result.Head.RoughProblem = attorney.RoughProblem
	result.Head.OutRange = attorney.OutRange
	result.Head.SpecificProblem = attorney.SpecificProblem
	totalWorkHour := 0.0
	totalPartsCost := 0.0
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
	projectResult := database.Raw("select time_overview.project_number as project_number, project_name as project_name, "+timeType+" as project_time, remark as project_remark\n"+
		"from arrangement inner join time_overview on arrangement.project_number = time_overview.project_number\n"+
		"where order_number = ?", attorneyNo).Scan(&result.AttorneyProject.Project) //找到该委托所有的维修项目
	result.AttorneyProject.Num = int(projectResult.RowsAffected) //记录维修项目数量
	var remarks []string                                         //备注字典
	for i := range result.AttorneyProject.Project {              //遍历所有维修项目
		totalWorkHour += result.AttorneyProject.Project[i].ProjectTime                    //总工时增加
		remarkList := strings.Split(result.AttorneyProject.Project[i].ProjectRemark, "*") //按照*分割备注
		result.AttorneyProject.Project[i].ProjectRemark = ""                              //处理后的备注，先置空
		if len(remarkList) > 1 {                                                          //备注不为空
			for j := 1; j < len(remarkList); j++ { //0是空，1开始有值
				if j == len(remarkList)-1 { //如果是最后一条备注
					result.AttorneyProject.Project[i].ProjectRemark += "*" + remarkList[j] //不加换行符
				} else {
					result.AttorneyProject.Project[i].ProjectRemark += "*" + remarkList[j] + "\n" //加换行符
				}
				remarks = append(remarks, "*"+remarkList[j]) //将这条备注添加到字典中
			}
		}
		partsResult := database.Raw("select repair_parts.parts_number as parts_number, parts_name as parts_name, parts_count as parts_count, parts_cost as parts_single_price\n"+
			"from repair_parts inner join parts_overview on repair_parts.parts_number=parts_overview.parts_number\n"+
			"where repair_parts.order_number = ? and project_number = ?", attorneyNo, result.AttorneyProject.Project[i].ProjectNumber).Scan(&result.AttorneyProject.Project[i].ProjectParts) //找这个项目的零件
		result.AttorneyProject.Project[i].ProjectPartsCount = int(partsResult.RowsAffected) //零件数量
		for j := range result.AttorneyProject.Project[i].ProjectParts {                     //遍历零件
			result.AttorneyProject.Project[i].ProjectParts[j].PartsTotalPrice = result.AttorneyProject.Project[i].ProjectParts[j].PartsSinglePrice * float64(result.AttorneyProject.Project[i].ProjectParts[j].PartsCount) //计算该零件的总价
			totalPartsCost += result.AttorneyProject.Project[i].ProjectParts[j].PartsTotalPrice                                                                                                                            //零件费用总计增加
		}
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
	result.AttorneyProject.Remark = strRemark                                                                                                                                                                                                                                                                                                     //将备注结果加到返回结果中
	result.AttorneyProject.TotalWorkHour = totalWorkHour                                                                                                                                                                                                                                                                                          //将总工时加到返回结果中
	result.AttorneyProject.TotalPartsCost = totalPartsCost                                                                                                                                                                                                                                                                                        //将零件总计加到返回结果中
	result.Tail.WorkHourPrice = 24                                                                                                                                                                                                                                                                                                                //将工时单价加到返回结果中
	totalPrice := (result.AttorneyProject.TotalWorkHour*float64(result.Tail.WorkHourPrice) + result.AttorneyProject.TotalPartsCost) * float64(result.Head.DiscountRate) * 0.01                                                                                                                                                                    //计算维修总价
	result.Tail.TotalPrice = "(" + strconv.FormatFloat(result.AttorneyProject.TotalWorkHour, 'f', -1, 64) + "*" + strconv.Itoa(result.Tail.WorkHourPrice) + "+" + strconv.FormatFloat(result.AttorneyProject.TotalPartsCost, 'f', -1, 64) + ")*" + strconv.Itoa(result.Head.DiscountRate) + "% = " + strconv.FormatFloat(totalPrice, 'f', -1, 64) //将计算过程与结果加到返回结果中
	result.Tail.CalculateWay = "（总工时*工时单价+材料费）*折扣率"                                                                                                                                                                                                                                                                                               //将计算方式加到返回结果中
	c.JSON(http.StatusOK, result)
}

func getUsername(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.String(http.StatusOK, username)
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
				c.HTML(http.StatusOK, "index.html", nil)
			}
		}
	}
}
