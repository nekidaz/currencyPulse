package initializers

import (
	"HalykTZ/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.CurrencyRate{})
}
