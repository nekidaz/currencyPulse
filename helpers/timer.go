// Файл HalykTZ/helpers/helpers.go

package helpers

import (
	"fmt"
	"time"
)

// тут таймер который обновляет бд каждые 6 часов
func UpdateCurrencyDataPeriodically(interval time.Duration, updateFunc func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("8")
			if err := updateFunc(); err != nil {
				fmt.Println("Error updating currency data:", err)
			}
		}
	}
}

// получаем актуальную дату
func GetTodayDate() string {
	todayday := time.Now()
	return todayday.Format("02.01.2006")
}
