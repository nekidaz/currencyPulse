package controllers

import (
	"HalykTZ/data_fetcher"
	"HalykTZ/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetAllCurrency(c *gin.Context) {
	// Получаем данные курсов валют из Redis или обновляем их при необходимости
	jsonData, err := data_fetcher.GetCurrencyRatesFromRedis()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting currency data")
		return
	}

	// Отправляем клиенту данные
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(jsonData))
}

func GetCurrencyByCode(c *gin.Context) {
	// Получаем параметр "code" из URL
	code := strings.ToUpper(c.Param("code"))

	// Получаем данные курсов валют из Redis или обновляем их при необходимости
	jsonData, err := data_fetcher.GetCurrencyRatesFromRedis()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting currency data")
		return
	}

	// Декодируем JSON и ищем курс валюты по коду
	var currencyRatesJSON struct {
		Items []models.CurrencyRate `json:"items"`
	}
	if err := json.Unmarshal(jsonData, &currencyRatesJSON); err != nil {
		c.String(http.StatusInternalServerError, "Error unmarshalling JSON")
		return
	}

	for _, rate := range currencyRatesJSON.Items {
		if rate.Title == code {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusOK, gin.H{
				"Title":       rate.Title,
				"FullName":    rate.Fullname,
				"Description": rate.Description,
				"Quant":       rate.Quant,
				"Index":       rate.Index,
				"Change":      rate.Change,
			})
			return
		}
	}

	c.String(http.StatusNotFound, "Currency not found")
}

func UpdateData(c *gin.Context) {
	// Обновляем данные курсов валют
	if err := data_fetcher.UpdateCurrencyData(); err != nil {
		c.String(http.StatusInternalServerError, "Error updating currency data")
		return
	}

	c.String(http.StatusOK, "Currency data updated successfully")
}
