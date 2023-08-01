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

func FetchCurrencyData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching currency rates: %v", err)
	}
	defer resp.Body.Close()

	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	return xmlData, nil
}

func UnmarshalXML(xmlData []byte) ([]models.CurrencyRate, error) {
	var currencyRatesXML struct {
		Items []models.CurrencyRate `xml:"item"`
	}
	if err := xml.Unmarshal(xmlData, &currencyRatesXML); err != nil {
		return nil, fmt.Errorf("Error unmarshalling XML: %v", err)
	}

	return currencyRatesXML.Items, nil
}

func UpdateOrCreateCurrencyRate(rate *models.CurrencyRate) error {
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

	return nil
}

func UpdateCurrencyData() error {
	todaysDate := helpers.GetTodayDate()
	url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", todaysDate)

	xmlData, err := FetchCurrencyData(url)
	if err != nil {
		return err
	}

	currencyRatesXML, err := UnmarshalXML(xmlData)
	if err != nil {
		return err
	}

	// Сохраняем данные в базу данных
	for _, rate := range currencyRatesXML {
		if err := UpdateOrCreateCurrencyRate(&rate); err != nil {
			return err
		}
	}

	if err := saveCurrencyRatesToRedis(currencyRatesXML); err != nil {
		return err
	}

	return nil
}

func saveCurrencyRatesToRedis(currencyRates []models.CurrencyRate) error {
	type CurrencyRatesJSON struct {
		Items []models.CurrencyRate `json:"items"`
	}

	ratesJSON := CurrencyRatesJSON{
		Items: currencyRates,
	}

	// структуру в json
	jsonData, err := json.Marshal(ratesJSON)
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v", err)
	}

	//сохраняем данные в редис в виде json
	ctx := context.Background()
	err = initializers.RDB.Set(ctx, "currency_rates", jsonData, 6*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("Failed to save data to Redis: %v", err)
	}

	return nil
}

// Если данных нет в Redis или они устарели, функция обновляет данные и сохраняет их в Redis.
func GetCurrencyRatesFromRedis() ([]byte, error) {
	ctx := context.Background()

	// Проверяем наличие данных в Redis
	value, err := initializers.RDB.Get(ctx, "currency_rates").Result()
	if err == nil {
		// Данные найдены в Redis, возвращаем их
		return []byte(value), nil
	}

	// Данные отсутствуют или устарели, обновляем данные и сохраняем их в Redis
	if err := UpdateCurrencyData(); err != nil {
		return nil, err
	}

	// Получаем обновленные данные из Redis и возвращаем их
	value, err = initializers.RDB.Get(ctx, "currency_rates").Result()
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}
