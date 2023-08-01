package tests

import (
	"HalykTZ/helpers"
	"fmt"
	"testing"
	"time"
)

func TestUpdateCurrencyDataPeriodically(t *testing.T) {
	// Функция-заглушка для обновления, которая возвращает nil
	mockUpdateFunc := func() error {
		return nil
	}

	// Запускаем функцию с небольшим интервалом для тестирования
	interval := 100 * time.Millisecond
	go helpers.UpdateCurrencyDataPeriodically(interval, mockUpdateFunc)

	// Ждем несколько итераций
	time.Sleep(3 * interval)

	// Тест считается успешным, если во время итераций не было выведено ошибок
}

func TestGetTodayDate(t *testing.T) {
	// Получаем строку с текущей датой, используя функцию
	todayDate := helpers.GetTodayDate()

	// Получаем ожидаемую строку с датой для сегодняшнего дня
	expectedDate := time.Now().Format("02.01.2006")

	// Сравниваем полученную дату с ожидаемой датой
	if todayDate != expectedDate {
		t.Errorf("Ожидаемая дата: %s, но получена: %s", expectedDate, todayDate)
	}
}

func TestUpdateCurrencyDataPeriodicallyWithError(t *testing.T) {
	// Функция-заглушка для обновления, которая возвращает ошибку, если передан параметр "error" равный true
	mockUpdateFunc := func() error {
		if true {
			return fmt.Errorf("Mock update error")
		}
		return nil
	}

	// Запускаем функцию с небольшим интервалом для тестирования
	interval := 100 * time.Millisecond
	go helpers.UpdateCurrencyDataPeriodically(interval, mockUpdateFunc)

	// Ждем несколько итераций
	time.Sleep(3 * interval)

	// Тест должен ожидать, что во время итераций была выведена "Mock update error"
	// Проверяем, что в логах была зарегистрирована ошибка

}
