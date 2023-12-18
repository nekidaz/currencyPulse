package initializers

import "github.com/nekidaz/currencyPulse/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.CurrencyRate{})
}
