package controllers

import (
	"HalykTZ/data_fetcher"
	"HalykTZ/initializers"
	"HalykTZ/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetAllCurrency(c *gin.Context) {
	ctx := context.Background()

	// Проверяем наличие данных в Redis
	value, err := initializers.RDB.Get(ctx, "currency_rates").Result()
	if err == nil {
		// Данные найдены в Redis, возвращаем их клиенту
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, value)
		return
	}

	// Данные в Redis отсутствуют, проверяем базу данных
	var currencyRates []models.CurrencyRate
	if err := initializers.DB.Find(&currencyRates).Error; err == nil {
		// Данные найдены в базе данных, сохраняем их в Redis и возвращаем клиенту
		jsonData, err := json.Marshal(currencyRates)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error marshalling JSON")
			return
		}

		// Сохраняем данные в Redis на 6 часов
		err = initializers.RDB.Set(ctx, "currency_rates", jsonData, 6*time.Hour).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to save data to Redis")
			return
		}

		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(jsonData))
		return
	}

	// Данные не найдены в Redis и базе данных, обновляем данные и возвращаем клиенту
	if err := data_fetcher.UpdateCurrencyData(); err != nil {
		c.String(http.StatusInternalServerError, "Error updating currency data")
		return
	}

	// Повторно вызываем функцию, чтобы получить данные из Redis
	GetAllCurrency(c)
}

func GetCurrencyByCode(c *gin.Context) {
	ctx := context.Background()

	// Получаем параметр "code" из URL
	code := c.Param("code")

	// Проверяем наличие данных в Redis
	value, err := initializers.RDB.Get(ctx, "currency_rates").Result()
	if err == nil {
		// Данные найдены в Redis, декодируем JSON и ищем курс валюты по коду
		var currencyRatesXML struct {
			Items []models.CurrencyRate `json:"items"`
		}
		if err := json.Unmarshal([]byte(value), &currencyRatesXML); err != nil {
			c.String(http.StatusInternalServerError, "Error unmarshalling JSON")
			return
		}

		for _, rate := range currencyRatesXML.Items {
			if rate.Title == code {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusOK, rate)
				return
			}
		}

		c.String(http.StatusNotFound, "Currency not found")
		return
	}

	// Данные в Redis отсутствуют, проверяем базу данных
	var rate models.CurrencyRate
	if err := initializers.DB.Where("title = ?", code).First(&rate).Error; err == nil {
		// Курс валюты найден в базе данных, сохраняем данные в Redis и возвращаем клиенту
		jsonData, err := json.Marshal(rate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error marshalling JSON")
			return
		}

		// Сохраняем данные в Redis на 6 часов
		err = initializers.RDB.Set(ctx, "currency_rates", jsonData, 6*time.Hour).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to save data to Redis")
			return
		}

		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(jsonData))
		return
	}

	// Курс валюты не найден в Redis и базе данных, обновляем данные и возвращаем клиенту
	if err := data_fetcher.UpdateCurrencyData(); err != nil {
		c.String(http.StatusInternalServerError, "Error updating currency data")
		return
	}

	// Повторно вызываем функцию, чтобы получить данные из Redis
	GetCurrencyByCode(c)
}
