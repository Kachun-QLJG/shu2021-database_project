package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiechui1994/gopdf"
	"github.com/tiechui1994/gopdf/core"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	time2 "time"
)

const (
	FONT_KT = "楷体"
	FONT_ST = "宋体"
	FONT_TN = "Times New Roman"
)

var (
	largeFont           = core.Font{Family: FONT_KT, Size: 15}
	textFont            = core.Font{Family: FONT_ST, Size: 10}
	usernameForPdfGen   string
	attorneyNoForPdfGen string
)

func generatePdf(c *gin.Context) {
	usernameForPdfGen = c.MustGet("username").(string)
	attorneyNoForPdfGen = c.Query("attorney_no")
	filename := genPdf(usernameForPdfGen, attorneyNoForPdfGen)
	c.File(filename)
}

func checkPermissionForPdf(username string, group string, attorneyNo string) string {
	if group == "维修员" {
		return ""
	}
	if group == "普通用户" {
		var user User
		database.First(&user, "contact_tel = ?", username)
		var attorney Attorney
		result := database.First(&attorney, "user_id = ? and number = ?", user.Number, attorneyNo)
		if result.RowsAffected == 1 {
			return "./files/generatedPDF/" + username + "/" + attorneyNo + "/" + attorneyNo + ".pdf"
		}
	} else if group == "业务员" {
		var attorney Attorney
		result := database.First(&attorney, "salesman_id = ? and number = ?", username, attorneyNo)
		if result.RowsAffected == 1 {
			var res struct{ Username string }
			database.Raw("select contact_tel as username from user where number = ?", attorney.UserID).Scan(&res)
			username = res.Username
			return "./files/generatedPDF/" + username + "/" + attorneyNo + "/" + attorneyNo + ".pdf"
		}
	}
	return ""
}
func showPdf(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	attorneyNo := c.Query("attorney_no")
	path := checkPermissionForPdf(username, group, attorneyNo)
	if path != "" {
		c.File(path)
	} else {
		c.String(http.StatusForbidden, "无权限")
	}
}

func downloadPdf(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	attorneyNo := c.Query("attorney_no")
	path := checkPermissionForPdf(username, group, attorneyNo)
	if path != "" {
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "委托书No"+attorneyNo+".pdf")) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(path)
	} else {
		c.String(http.StatusForbidden, "无权限")
	}
}

func genPdf(username string, attorneyNo string) string {
	usernameForPdfGen = username
	attorneyNoForPdfGen = attorneyNo
	pdf := core.CreateReport()
	font1 := core.FontMap{
		FontName: FONT_KT,
		FileName: "./files/fonts/kaiti.ttf",
	}
	font2 := core.FontMap{
		FontName: FONT_ST,
		FileName: "./files/fonts/songti.ttf",
	}
	font3 := core.FontMap{
		FontName: FONT_TN,
		FileName: "./files/fonts/times.ttf",
	}
	pdf.SetFonts([]*core.FontMap{&font1, &font2, &font3})
	pdf.SetPage("A4", "P")
	pdf.FisrtPageNeedHeader = true
	pdf.FisrtPageNeedFooter = true
	pdf.RegisterExecutor(ComplexReportExecutor, core.Detail)
	pdf.RegisterExecutor(ComplexReportFooterExecutor, core.Footer)
	pdf.RegisterExecutor(ComplexReportHeaderExecutor, core.Header)
	path := "./files/generatedPDF/" + username + "/" + attorneyNo
	os.MkdirAll(path, os.ModePerm)
	path = path + "/" + attorneyNo + ".pdf"
	pdf.Execute(path)
	return path
}

func ComplexReportExecutor(report *core.Report) {
	var user User
	database.First(&user, "contact_tel = ?", usernameForPdfGen)
	var attorney Attorney
	database.First(&attorney, "number = ?", attorneyNoForPdfGen)
	var salesman Salesman
	database.First(&salesman, "number = ?", attorney.SalesmanID)
	var vehicle Vehicle
	database.First(&vehicle, "number = ? and user_id = ?", attorney.VehicleNumber, attorney.UserID)
	var arrangement []struct {
		OrderNumber   string
		ProjectNumber string
	}
	database.Raw("select distinct order_number as order_number, project_number as project_number\n"+
		"from arrangement where order_number = ?", attorney.Number).Scan(&arrangement) //查这个订单有多少维修项目
	columnCount := 0
	for _, data := range arrangement {
		var repairParts []RepairParts
		result := database.Find(&repairParts, "order_number = ? and project_number = ?", data.OrderNumber, data.ProjectNumber)
		if result.RowsAffected == 0 {
			columnCount++
		} else {
			columnCount += int(result.RowsAffected)
		}
	}

	var remarks []string
	lineSpace := 1.0
	lineHeight := 16.0

	title1 := gopdf.NewDiv(lineHeight, lineSpace, report)
	title1.SetFont(largeFont)
	title1.HorizontalCentered().SetContent("502汽车维修站\n维修委托书").GenerateAtomicCell()
	date := time2.Now().Format("2006年1月2日")

	no := gopdf.NewDiv(lineHeight, lineSpace, report)
	no.SetFont(textFont).SetContent("No." + attorneyNoForPdfGen).GenerateAtomicCell()
	time := gopdf.NewDiv(lineHeight, lineSpace, report)
	time.SetMarign(core.Scope{Top: lineHeight * -1})
	time.SetFont(textFont).RightAlign().SetContent("生成日期：" + date).GenerateAtomicCell()
	border := core.NewScope(4.0, 4.0, 4.0, 0) //表格内的margin
	rows, cols := 21+columnCount, 23
	table := gopdf.NewTable(cols, int(rows), 450, lineHeight, report)

	userinfo := table.NewCellByRange(3, 2)
	cell := gopdf.NewTextCell(table.GetColWidth(0, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("客户信息")
	userinfo.SetElement(cell)
	userNumber := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(0, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("客户编号")
	userNumber.SetElement(cell)
	contactPerson := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(0, 6), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("联系人")
	contactPerson.SetElement(cell)
	contactTel := table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(0, 9), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("联系方式")
	contactTel.SetElement(cell)
	salesmanInfo := table.NewCellByRange(3, 2)
	cell = gopdf.NewTextCell(table.GetColWidth(0, 13), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("业务员\n信息")
	salesmanInfo.SetElement(cell)
	salesmanNumber := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(0, 16), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("编号")
	salesmanNumber.SetElement(cell)
	salesmanName := table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(0, 19), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("姓名")
	salesmanName.SetElement(cell)
	userNumberBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(1, 3), lineHeight, lineSpace, report)

	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(user.Number) //填客户编号
	userNumberBlank.SetElement(cell)
	contactPersonBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(1, 6), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(user.ContactPerson) //填联系人
	contactPersonBlank.SetElement(cell)
	contactTelBlank := table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(1, 9), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(user.ContactTel) //填联系方式
	contactTelBlank.SetElement(cell)
	salesmanNumberBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(1, 16), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(salesman.Number) //填编号
	salesmanNumberBlank.SetElement(cell)
	salesmanNameBlank := table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(1, 19), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(salesman.Name) //填姓名
	salesmanNameBlank.SetElement(cell)

	vehicleInfo := table.NewCellByRange(3, 4)
	cell = gopdf.NewTextCell(table.GetColWidth(2, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("车辆信息")
	vehicleInfo.SetElement(cell)
	licenseNumber := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(2, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("车牌号")
	licenseNumber.SetElement(cell)
	vehicleNumber := table.NewCellByRange(7, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(2, 8), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("车架号")
	vehicleNumber.SetElement(cell)
	vehicleModel := table.NewCellByRange(8, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(2, 15), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("车型")
	vehicleModel.SetElement(cell)
	licenseNumberBlank := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(3, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(vehicle.LicenseNumber) //填车牌号
	licenseNumberBlank.SetElement(cell)
	vehicleNumberBlank := table.NewCellByRange(7, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(3, 8), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(vehicle.Number) //填车架号
	vehicleNumberBlank.SetElement(cell)
	vehicleModelBlank := table.NewCellByRange(8, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(3, 15), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(vehicle.Model) //填车型
	vehicleModelBlank.SetElement(cell)

	vehicleType := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("车辆类型")
	vehicleType.SetElement(cell)
	startPetrol := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 6), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("进场油量")
	startPetrol.SetElement(cell)
	endPetrol := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 9), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("出场油量")
	endPetrol.SetElement(cell)
	startMile := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 12), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("进场里程数")
	startMile.SetElement(cell)
	endMile := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 15), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("出场里程数")
	endMile.SetElement(cell)
	taskInfo := table.NewCellByRange(2, 4)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 18), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("作业\n信息")
	taskInfo.SetElement(cell)
	repairType := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(4, 20), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("维修类型")
	repairType.SetElement(cell)

	vehicleTypeBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(5, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(vehicle.Type) //填车辆类型
	vehicleTypeBlank.SetElement(cell)
	startPetrolBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(5, 6), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(attorney.StartPetrol, 'f', -1, 64)) //填进场油量
	startPetrolBlank.SetElement(cell)
	endPetrolBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(5, 9), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(attorney.EndPetrol, 'f', -1, 64)) //填出场油量
	endPetrolBlank.SetElement(cell)
	startMileBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(5, 12), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(attorney.StartMile, 'f', -1, 64)) //填进场里程数
	startMileBlank.SetElement(cell)
	endMileBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(5, 15), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(attorney.EndMile, 'f', -1, 64)) //填出场里程数
	endMileBlank.SetElement(cell)
	repairTypeBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(5, 20), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.RepairType) //填维修类型
	repairTypeBlank.SetElement(cell)

	checkoutInfo := table.NewCellByRange(3, 2)
	cell = gopdf.NewTextCell(table.GetColWidth(6, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("结算信息")
	checkoutInfo.SetElement(cell)
	payWay := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(6, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("结算方式")
	payWay.SetElement(cell)
	discountRate := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(6, 6), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("折扣率")
	discountRate.SetElement(cell)
	predictFinishTime := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(6, 8), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("预计完工时间")
	predictFinishTime.SetElement(cell)
	actualFinishTime := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(6, 13), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("实际完工时间")
	actualFinishTime.SetElement(cell)
	taskClassification := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(6, 20), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("作业分类")
	taskClassification.SetElement(cell)

	payWayBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(7, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.PayMethod) //填结算方式
	payWayBlank.SetElement(cell)
	discountRateBlank := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(7, 6), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.Itoa(user.DiscountRate) + "%") //填折扣率
	discountRateBlank.SetElement(cell)
	predictFinishTimeBlank := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(7, 8), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.PredictFinishTime) //填预计完工时间
	predictFinishTimeBlank.SetElement(cell)
	actualFinishTimeBlank := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(7, 13), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.ActualFinishTime) //填实际完工时间
	actualFinishTimeBlank.SetElement(cell)
	taskClassificationBlank := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(7, 20), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.Classification) //填作业分类
	taskClassificationBlank.SetElement(cell)

	userDescription := table.NewCellByRange(15, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(8, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("客户报障描述")
	userDescription.SetElement(cell)
	outRange := table.NewCellByRange(8, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(8, 15), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("非维修范围")
	outRange.SetElement(cell)
	userDescriptionBlank := table.NewCellByRange(15, 2)
	cell = gopdf.NewTextCell(table.GetColWidth(9, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.RoughProblem) //填客户报障描述
	userDescriptionBlank.SetElement(cell)
	outRangeBlank := table.NewCellByRange(8, 2)
	cell = gopdf.NewTextCell(table.GetColWidth(9, 15), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.OutRange) //填非维修范围
	outRangeBlank.SetElement(cell)

	ErrorDescription := table.NewCellByRange(23, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(11, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("检修故障描述")
	ErrorDescription.SetElement(cell)
	ErrorDescriptionBlank := table.NewCellByRange(23, 2)
	cell = gopdf.NewTextCell(table.GetColWidth(12, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(attorney.SpecificProblem) //填检修故障描述
	ErrorDescriptionBlank.SetElement(cell)

	repairNumber := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("项目编号")
	repairNumber.SetElement(cell)
	repairName := table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("维修项目名称")
	repairName.SetElement(cell)
	repairTime := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 7), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("工时")
	repairTime.SetElement(cell)
	partsNumber := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 9), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("零件号")
	partsNumber.SetElement(cell)
	partsName := table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 12), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("零件名")
	partsName.SetElement(cell)
	num := table.NewCellByRange(1, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 16), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("数量")
	num.SetElement(cell)
	price := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 17), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("单价")
	price.SetElement(cell)
	totalPrice := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 19), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("总价")
	totalPrice.SetElement(cell)
	remark := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(14, 21), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("备注")
	remark.SetElement(cell)
	currentRow := 15
	timeCount := 0.0
	priceCount := 0.0
	type Project struct {
		ProjectNum string
	}
	var project []Project
	database.Select("distinct(project_number) as project_num").Table("arrangement").Where("order_number = ?", attorney.Number).Scan(&project)
	for projectNum := range project {
		var repairParts []RepairParts
		result := database.Find(&repairParts, "order_number = ? and project_number = ?", attorney.Number, project[projectNum].ProjectNum) //查每一个维修项目的零件
		resultNum := int(result.RowsAffected)
		noParts := false
		if resultNum == 0 {
			resultNum = 1
			noParts = true
		}
		var timeOverview TimeOverview
		database.Find(&timeOverview, "project_number = ?", project[projectNum].ProjectNum)
		repairNumber = table.NewCellByRange(3, resultNum)
		cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 0), lineHeight, lineSpace, report)
		cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(timeOverview.ProjectNumber) //填项目编号
		repairNumber.SetElement(cell)
		repairName = table.NewCellByRange(4, resultNum)
		cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 3), lineHeight, lineSpace, report)
		cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(timeOverview.ProjectName) //填维修项目名称
		repairName.SetElement(cell)
		repairTime = table.NewCellByRange(2, resultNum)
		cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 7), lineHeight, lineSpace, report)
		if vehicle.Type == "轿车-A" {
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(timeOverview.TimeA, 'f', -1, 64)) //填工时
			timeCount += timeOverview.TimeA
		} else if vehicle.Type == "轿车-B" {
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(timeOverview.TimeB, 'f', -1, 64)) //填工时
			timeCount += timeOverview.TimeB
		} else if vehicle.Type == "轿车-C" {
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(timeOverview.TimeC, 'f', -1, 64)) //填工时
			timeCount += timeOverview.TimeC
		} else if vehicle.Type == "轿车-D" {
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(timeOverview.TimeD, 'f', -1, 64)) //填工时
			timeCount += timeOverview.TimeD
		} else {
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(timeOverview.TimeE, 'f', -1, 64)) //填工时
			timeCount += timeOverview.TimeE
		}
		repairTime.SetElement(cell)
		flag := true
		if noParts {
			partsNumber = table.NewCellByRange(3, 1)
			cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 9), lineHeight, lineSpace, report)
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("") //填零件号
			partsNumber.SetElement(cell)
			partsName = table.NewCellByRange(4, 1)
			cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 12), lineHeight, lineSpace, report)
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("") //填零件名
			partsName.SetElement(cell)
			num = table.NewCellByRange(1, 1)
			cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 16), lineHeight, lineSpace, report)
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("") //填数量
			num.SetElement(cell)
			price = table.NewCellByRange(2, 1)
			cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 17), lineHeight, lineSpace, report)
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("") //填单价
			price.SetElement(cell)
			totalPrice = table.NewCellByRange(2, 1)
			cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 19), lineHeight, lineSpace, report)
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("") //填总价
			totalPrice.SetElement(cell)
			remark = table.NewCellByRange(2, resultNum)
			cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 21), lineHeight, lineSpace, report)
			remarkList := strings.Split(timeOverview.Remark, "*")
			remarkOnTable := ""
			if len(remarkList) > 1 { //备注不为空
				for i := 1; i < len(remarkList); i++ { //0是空，1开始有值
					remarkOnTable += "*" + remarkList[i] + "\n"
					remarks = append(remarks, "*"+remarkList[i])
				}
			}
			cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(remarkOnTable) //填备注
			remark.SetElement(cell)
			currentRow++
		} else {
			for _, line := range repairParts {
				partsNumber = table.NewCellByRange(3, 1)
				cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 9), lineHeight, lineSpace, report)
				cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(line.PartsNumber) //填零件号
				partsNumber.SetElement(cell)
				var parts PartsOverview
				database.First(&parts, "parts_number = ?", line.PartsNumber)
				partsName = table.NewCellByRange(4, 1)
				cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 12), lineHeight, lineSpace, report)
				cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(parts.PartsName) //填零件名
				partsName.SetElement(cell)
				num = table.NewCellByRange(1, 1)
				cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 16), lineHeight, lineSpace, report)
				cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.Itoa(line.PartsCount)) //填数量
				num.SetElement(cell)
				price = table.NewCellByRange(2, 1)
				cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 17), lineHeight, lineSpace, report)
				cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(parts.PartsCost, 'f', -1, 64)) //填单价
				price.SetElement(cell)
				totalPrice = table.NewCellByRange(2, 1)
				cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 19), lineHeight, lineSpace, report)
				cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(float64(line.PartsCount)*parts.PartsCost, 'f', -1, 64)) //填总价
				priceCount += float64(line.PartsCount) * parts.PartsCost
				totalPrice.SetElement(cell)
				if flag {
					remark = table.NewCellByRange(2, resultNum)
					cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 21), lineHeight, lineSpace, report)
					remarkList := strings.Split(timeOverview.Remark, "*")
					remarkOnTable := ""
					if len(remarkList) > 1 { //备注不为空
						for i := 1; i < len(remarkList); i++ { //0是空，1开始有值
							remarkOnTable += "*" + remarkList[i] + "\n"
							remarks = append(remarks, "*"+remarkList[i])
						}
					}
					cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(remarkOnTable) //填备注
					remark.SetElement(cell)
					flag = false
				}
				currentRow++
			}
		}
	}

	blank := table.NewCellByRange(5, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("")
	blank.SetElement(cell)
	totalTime := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 5), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("总工时")
	totalTime.SetElement(cell)
	totalTimeBlank := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 7), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(timeCount, 'f', -1, 64)) //填总工时
	totalTimeBlank.SetElement(cell)
	blank = table.NewCellByRange(8, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 9), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("")
	blank.SetElement(cell)
	totalCost := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 17), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("总计")
	totalCost.SetElement(cell)
	totalCostBlank := table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 19), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strconv.FormatFloat(priceCount, 'f', -1, 64)) //填总计
	totalCostBlank.SetElement(cell)
	blank = table.NewCellByRange(2, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 21), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("")
	blank.SetElement(cell)
	currentRow++

	allRemark := table.NewCellByRange(23, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("备注")
	allRemark.SetElement(cell)
	currentRow++

	strRemark := ""
	if len(remarks) != 0 {
		sort.Strings(remarks)
		for _, value := range remarks {
			var remarkDatabase Remark
			database.First(&remarkDatabase, "remark_number = ?", value)
			strRemark += remarkDatabase.RemarkNumber + " " + remarkDatabase.Content + "\n"
		}
		strRemark = strRemark[:len(strRemark)-1]
	}
	allRemarkBlank := table.NewCellByRange(23, 2)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent(strRemark) //填备注
	allRemarkBlank.SetElement(cell)
	currentRow += 2

	timeCost := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("工时单价")
	timeCost.SetElement(cell)
	timeCostBlank := table.NewCellByRange(1, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 3), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("24")
	timeCostBlank.SetElement(cell)
	allPrice := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 4), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("维修总价")
	allPrice.SetElement(cell)
	allPriceBlank := table.NewCellByRange(16, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 7), lineHeight, lineSpace, report)
	payPrice := (timeCount*24 + priceCount) * float64(user.DiscountRate) / 100
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("(" + strconv.FormatFloat(timeCount, 'f', -1, 64) + "*24+" + strconv.FormatFloat(priceCount, 'f', -1, 64) + ")*" + strconv.Itoa(user.DiscountRate) + "%=" + strconv.FormatFloat(payPrice, 'f', -1, 64)) //填维修总价
	allPriceBlank.SetElement(cell)
	currentRow++

	blank = table.NewCellByRange(4, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 0), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("")
	blank.SetElement(cell)
	calculateMethod := table.NewCellByRange(3, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 4), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("计价方式")
	calculateMethod.SetElement(cell)
	calculateMethodBlank := table.NewCellByRange(16, 1)
	cell = gopdf.NewTextCell(table.GetColWidth(currentRow, 7), lineHeight, lineSpace, report)
	cell.SetFont(textFont).SetBorder(border).HorizontalCentered().VerticalCentered().SetContent("（总工时*工时单价+材料费）*折扣率")
	calculateMethodBlank.SetElement(cell)

	err := table.GenerateAtomicCell()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func ComplexReportFooterExecutor(report *core.Report) {
	content := fmt.Sprintf("第 %v / {#TotalPage#} 页", report.GetCurrentPageNo())
	footer := gopdf.NewSpan(10, 0, report)
	footer.SetFont(textFont)
	footer.SetFontColor("0, 0, 0")
	footer.RightAlign().SetContent(content).GenerateAtomicCell()
}
func ComplexReportHeaderExecutor(report *core.Report) {
	content := "维修委托书"
	header := gopdf.NewSpan(10, 0, report)
	header.SetFont(textFont)
	header.SetFontColor("0,0,0")
	header.SetBorder(core.Scope{Top: 35})
	header.HorizontalCentered().SetContent(content).GenerateAtomicCell()
	line := gopdf.NewHLine(report)
	line.SetMargin(core.Scope{Top: 5}).GenerateAtomicCell()
}
