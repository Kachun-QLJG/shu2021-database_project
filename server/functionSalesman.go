package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func startSCheckOrders(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "业务员" {
		c.HTML(http.StatusOK, "error.html", gin.H{"data": "无权限！", "website": "/index", "webName": "主页"})
		return
	}
	c.HTML(http.StatusOK, "salesman_check_orders.html", gin.H{"username": username, "group": group})
}

func startTakeOrders(c *gin.Context) {
	username := c.MustGet("username").(string)
	group := c.MustGet("group").(string)
	if group != "业务员" {
		c.HTML(http.StatusOK, "error.html", gin.H{"data": "无权限！", "website": "/index", "webName": "主页"})
		return
	}
	c.HTML(http.StatusOK, "salesman_take_orders.html", gin.H{"username": username, "group": group})
}
