package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getPendingAttorney(c *gin.Context) {
	var attorney struct {
		Number  string
		CarType string
	}
	database.Raw("select attorney.number as number, type as car_type from attorney inner join vehicle on attorney.vehicle_number = vehicle.number where progress = '待处理'").Scan(&attorney)
	c.JSON(http.StatusOK, attorney)
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
