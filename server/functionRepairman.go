package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

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
