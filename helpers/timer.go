// Файл HalykTZ/helpers/helpers.go

package helpers

import (
	"fmt"
	"time"
)

func UpdateCurrencyDataPeriodically(interval time.Duration, updateFunc func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := updateFunc(); err != nil {
				fmt.Println("Error updating currency data:", err)
			}
		}
	}
}

func GetTodayDate() string {
	today := time.Now()
	return today.Format("02.01.2006")
}
