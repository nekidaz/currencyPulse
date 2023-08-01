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
	//подключение к постгре
	initializers.ConnectToPostgres()
	//подключение к редис
	initializers.ConnectToRedis()
	//запуск миграции
	initializers.SyncDatabase()
	//при запуске кода парсит данные и сохраняет в бд
	data_fetcher.UpdateCurrencyData()
	//период обновление
	period := time.Hour * 6
	// Запуск периодического обновления данных по interval
	go helpers.UpdateCurrencyDataPeriodically(period, data_fetcher.UpdateCurrencyData)
}

func main() {
	r := gin.Default()
	//получение всех валют
	r.GET("/rates", controllers.GetAllCurrency)
	//получение валюты по коду
	r.GET("/rates/:code", controllers.GetCurrencyByCode)
	//обновление данных в кэше и в бд
	r.POST("/update", controllers.UpdateData)

	r.Run(":8080")
}
