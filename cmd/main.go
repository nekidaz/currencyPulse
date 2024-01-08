package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nekidaz/currencyPulse/controllers"
	"github.com/nekidaz/currencyPulse/data_fetcher"
	"github.com/nekidaz/currencyPulse/helpers"
	"github.com/nekidaz/currencyPulse/initializers"
	"time"
)

func init() {
	initializers.ConnectToPostgres()
	initializers.ConnectToRedis()
	initializers.SyncDatabase()
	data_fetcher.UpdateCurrencyData()
	period := time.Hour * 6
	go helpers.UpdateCurrencyDataPeriodically(period, data_fetcher.UpdateCurrencyData)
}

func main() {
	r := gin.Default()
	r.GET("/rates", controllers.GetAllCurrency)
	r.GET("/rates/:code", controllers.GetCurrencyByCode)
	r.POST("/update", controllers.UpdateData)
	x
	r.Run(":8080")
}
