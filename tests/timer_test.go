package tests

import (
	"HalykTZ/helpers"
	"testing"
	"time"
)

func TestUpdateCurrencyDataPeriodically(t *testing.T) {
	// Mock update function that returns nil
	mockUpdateFunc := func() error {
		return nil
	}

	// Run the function with a small interval for testing purposes
	interval := 100 * time.Millisecond
	go helpers.UpdateCurrencyDataPeriodically(interval, mockUpdateFunc)

	// Wait for a few iterations to happen
	time.Sleep(3 * interval)

	// The test should pass if no error was printed during the iterations
}

func TestGetYesterdayDate(t *testing.T) {
	// Get the date string using the function
	yesterdayDate := helpers.GetYesterdayDate()

	// Get the expected date string for yesterday
	expectedDate := time.Now().Add(-24 * time.Hour).Format("02.01.2006")

	// Compare the obtained date with the expected date
	if yesterdayDate != expectedDate {
		t.Errorf("Expected date: %s, but got: %s", expectedDate, yesterdayDate)
	}
}
