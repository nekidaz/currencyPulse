package data_fetcher

import (
	"HalykTZ/helpers"
	"HalykTZ/initializers"
	"HalykTZ/models"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func UpdateCurrencyData() error {
	yesterdayDate := helpers.GetTodayDate()
	url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", yesterdayDate)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching currency rates: %v", err)
	}
	defer resp.Body.Close()

	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %v", err)
	}

	var currencyRatesXML struct {
		Items []models.CurrencyRate `xml:"item"`
	}
	if err := xml.Unmarshal(xmlData, &currencyRatesXML); err != nil {
		return fmt.Errorf("Error unmarshalling XML: %v", err)
	}

	// Сохраняем данные в базу данных
	for _, rate := range currencyRatesXML.Items {
		var existingRate models.CurrencyRate
		if err := initializers.DB.Where("title = ?", rate.Title).First(&existingRate).Error; err == nil {
			// Запись с таким же title найдена, обновляем ее поля
			existingRate.Fullname = rate.Fullname
			existingRate.Description = rate.Description
			existingRate.Quant = rate.Quant
			existingRate.Index = rate.Index
			existingRate.Change = rate.Change

			if err := initializers.DB.Save(&existingRate).Error; err != nil {
				return fmt.Errorf("Failed to update data in database: %v", err)
			}
		} else {
			// Запись с таким title не найдена, создаем новую запись
			if err := initializers.DB.Create(&rate).Error; err != nil {
				return fmt.Errorf("Failed to save data to database: %v", err)
			}
		}
	}

	jsonData, err := json.Marshal(currencyRatesXML)
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v", err)
	}

	// Сохраняем данные в Redis на 6 часов
	ctx := context.Background()
	err = initializers.RDB.Set(ctx, "currency_rates", jsonData, 6*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("Failed to save data to Redis: %v", err)
	}

	return nil
}
