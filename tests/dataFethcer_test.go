package tests

import (
	"github.com/nekidaz/currencyPulse/data_fetcher"
	"github.com/nekidaz/currencyPulse/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchCurrencyData(t *testing.T) {
	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отправляем тестовый XML ответ
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<rss version="2.0">
				<channel>
					<item>
						<title>USD</title>
						<fullname>US Dollar</fullname>
						<description>US Dollar Description</description>
						<quant>1</quant>
						<index>100</index>
						<change>0</change>
					</item>
				</channel>
			</rss>`))
	}))
	defer server.Close()

	// Вызываем функцию fetchCurrencyData для теста
	xmlData, err := data_fetcher.FetchCurrencyData(server.URL)
	assert.NoError(t, err, "Error fetching currency data")
	assert.NotNil(t, xmlData, "XML data should not be nil")
}

func TestUnmarshalXML(t *testing.T) {
	// Мокирование XML данных
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0">
				<item>
					<title>USD</title>
					<fullname>US Dollar</fullname>
					<description>5.00</description>
					<quant>1</quant>
					<index>UP</index>
					<change>0</change>
				</item>
		</rss>`)

	// Вызываем функцию unmarshalXML для теста
	currencyRates, err := data_fetcher.UnmarshalXML(xmlData)
	assert.NoError(t, err, "Error unmarshalling XML")
	assert.NotNil(t, currencyRates, "Currency rates should not be nil")
	assert.Len(t, currencyRates, 1, "There should be one currency rate")

	// Проверяем содержание десериализованного валютного курса
	expectedCurrencyRate := models.CurrencyRate{
		Title:       "USD",
		FullName:    "US Dollar",
		Description: 5.0,
		Quant:       1,
		Index:       "UP",
		Change:      0,
	}
	assert.Equal(t, expectedCurrencyRate, currencyRates[0], "Currency rate mismatch")
}
