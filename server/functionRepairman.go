package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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

func startChangeStatus(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.HTML(http.StatusOK, "error.html", gin.H{"data": "无权限！", "website": "/index", "webName": "主页"})
		return
	}
	c.HTML(http.StatusOK, "repairman_status_change.html", gin.H{"username": username, "group": group})
}

//func startRCheckOrders(c *gin.Context) {
//	username := c.MustGet("username").(string)
//	group := c.MustGet("group").(string)
//	if group != "维修员" {
//		c.HTML(http.StatusOK, "error.html", gin.H{"data": "无权限！", "website": "/index", "webName": "主页"})
//		return
//	}
//	c.HTML(http.StatusOK, "repairman_check_orders.html", gin.H{"username": username, "group": group})
//}

func checkStatus(c *gin.Context) {
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.String(http.StatusForbidden, "错误！")
		return
	}
	var repairman Repairman
	database.First(&repairman, "number = ?", number)
	c.String(http.StatusOK, repairman.Status)
}

func changeStatus(c *gin.Context) {
	status := c.PostForm("status")
	number := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.String(http.StatusBadRequest, "错误！")
		return
	}
	var repairman Repairman
	database.First(&repairman, "number = ?", number)
	database.Model(&repairman).Update("status", status) //更改状态为传过来的状态
	c.String(http.StatusOK, "修改成功！")
}
