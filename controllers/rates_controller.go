package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nekidaz/currencyPulse/data_fetcher"
	"github.com/nekidaz/currencyPulse/models"
	"net/http"
	"strings"
)

func GetAllCurrency(c *gin.Context) {
	jsonData, err := data_fetcher.GetCurrencyRatesFromRedis()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting currency data")
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(jsonData))
}

func GetCurrencyByCode(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))

	jsonData, err := data_fetcher.GetCurrencyRatesFromRedis()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting currency data")
		return
	}

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
				"FullName":    rate.FullName,
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
	if err := data_fetcher.UpdateCurrencyData(); err != nil {
		c.String(http.StatusInternalServerError, "Error updating currency data")
		return
	}

	c.String(http.StatusOK, "Currency data updated successfully")
}
