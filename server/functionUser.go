package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
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
	projectResult := database.Raw("select distinct time_overview.project_number as project_number, project_name as project_name, "+timeType+" as project_time, remark as project_remark\n"+
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

func getProcessingAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
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
	database.Table("attorney").Select("number as order_number, vehicle_number as vin, rough_problem as rough_problem, specific_problem as specific_problem, predict_finish_time as predict_finish_time, progress as progress").Where("user_id = ? and progress != '已完成'", user.Number).Order("progress, order_number").Scan(&result) //查询该用户的，状态不是已完成的委托
	c.JSON(http.StatusOK, result)
}

func getFinishedAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	var user User
	database.First(&user, "contact_tel = ?", username)
	type information struct {
		OrderNumber      string
		Vin              string
		RoughProblem     string
		SpecificProblem  string
		ActualFinishTime string
	}
	var result []information
	database.Table("attorney").Select("number as order_number, vehicle_number as vin, rough_problem as rough_problem, specific_problem as specific_problem, actual_finish_time as actual_finish_time").Where("user_id = ? and progress = '已完成'", user.Number).Order("order_number desc").Scan(&result) //查询该用户的，状态是已完成的委托
	c.JSON(http.StatusOK, result)
}

func getAttorneyDetail(c *gin.Context) {
	username := c.MustGet("username").(string)
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

func checkUserinfo(c *gin.Context) {
	username := c.MustGet("username").(string)
	var user User
	database.First(&user, "contact_tel = ?", username)
	if user.Name == "" || user.ContactPerson == "" || user.Property == "" {
		c.String(http.StatusOK, "未完成")
	} else {
		c.String(http.StatusOK, "完成")
	}

}
func createAttorney(c *gin.Context) {
	username := c.MustGet("username").(string)
	vin := c.PostForm("vin")
	payMethod := c.PostForm("pay_method")
	startTime := c.PostForm("start_time")
	roughProblem := c.PostForm("rough_problem")
	var user User
	database.First(&user, "contact_tel = ?", username)
	if user.Name == "" || user.ContactPerson == "" || user.Property == "" {
		c.String(http.StatusOK, "请完善个人信息后再提交新的委托申请！")
		return
	}
	var res Attorney
	result := database.Raw("select * from attorney where progress != '已完成' and vehicle_number = ?", vin).Scan(&res)
	if result.RowsAffected != 0 {
		c.String(http.StatusOK, "新增委托失败！该车辆正在维修中！请耐心等待~")
		return
	}
	date := time.Now().Format("20060102")
	var attorney Attorney
	number := database.Find(&attorney, "number like ?", date+"___").RowsAffected + 1
	strNumber := date + fmt.Sprintf("%03d", number)
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
	vin := c.Query("vin")
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
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("生成车牌失败！" + fmt.Sprint(err) + ": " + stderr.String())
		return
	}
}

func checkVehicle(c *gin.Context) {
	carNumber := c.Query("number")
	var vehicle Vehicle
	c.String(http.StatusOK, strconv.FormatInt(database.First(&vehicle, "number = ?", carNumber).RowsAffected, 10))
}

func startUCheckOrdersOngoing(c *gin.Context) {
	c.HTML(http.StatusOK, "user_check_orders_ongoing.html", nil)
}

func startUCheckOrdersFinished(c *gin.Context) {
	c.HTML(http.StatusOK, "user_check_orders_finished.html", nil)
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
		result := database.First(&attorney, "vehicle_number = ? and progress != '已完成'", number)
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
			c.JSON(http.StatusOK, gin.H{"status": "失败", "data": "该车辆有未完成的维修委托！绑定失败！"})
		}
	}
}

func getUserInfo(c *gin.Context) {
	number := c.MustGet("username").(string)
	var user struct {
		Number        string
		Name          string
		Property      string
		DiscountRate  int
		ContactPerson string
		ContactTel    string
	}
	database.Table("user").Where("contact_tel = ?", number).Select("number as number, name as name, property as property, discount_rate as discount_rate, contact_person as contact_person, contact_tel as contact_tel").Scan(&user)
	c.JSON(http.StatusOK, user)
}

func showPlate(c *gin.Context) {
	number := c.MustGet("username").(string)
	licenseNumber := c.Query("license_number")
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
	name := c.PostForm("name")
	property := c.PostForm("property")
	contactPerson := c.PostForm("contact_person")
	var user User
	database.First(&user, "contact_tel = ?", number)
	database.Model(&user).Update("name", name)                    //更改状态为传过来的状态
	database.Model(&user).Update("property", property)            //更改状态为传过来的状态
	database.Model(&user).Update("contact_person", contactPerson) //更改状态为传过来的状态
	c.String(http.StatusOK, "成功！")
}
