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
			fmt.Println("8")
			if err := updateFunc(); err != nil {
				fmt.Println("Error updating currency data:", err)
			}
		}
	}
}

func GetTodayDate() string {
	todayday := time.Now()
	return todayday.Format("02.01.2006")
}
