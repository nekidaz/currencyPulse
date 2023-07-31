package main

import (
	"HalykTZ/controllers"
	"HalykTZ/data_fetcher"
	"HalykTZ/helpers"
	"HalykTZ/initializers"
	"github.com/gin-gonic/gin"
	"time"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToPostgres()
	initializers.ConnectToRedis()
	initializers.SyncDatabase()
	data_fetcher.UpdateCurrencyData()

	interval := time.Hour * 6
	// Запуск периодического обновления данных каждые 24 часа
	go helpers.UpdateCurrencyDataPeriodically(interval, data_fetcher.UpdateCurrencyData)
}

func main() {
	r := gin.Default()

	r.GET("/rates", controllers.GetAllCurrency)
	r.GET("/rates/:code", controllers.GetCurrencyByCode)

	r.Run(":8080")
}
