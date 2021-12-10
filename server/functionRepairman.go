package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func startChangeStatus(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.HTML(http.StatusOK, "error.html", gin.H{"data": "无权限！", "website": "/index", "webName": "主页"})
		return
	}
	c.HTML(http.StatusOK, "repairman_status_change.html", gin.H{"username": username, "group": group})
}

func startCheckOrders(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "维修员" {
		c.HTML(http.StatusOK, "error.html", gin.H{"data": "无权限！", "website": "/index", "webName": "主页"})
		return
	}
	c.HTML(http.StatusOK, "repairman_check.html", gin.H{"username": username, "group": group})
}

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
