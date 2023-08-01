package models

import "gorm.io/gorm"

// модель валюты
type CurrencyRate struct {
	gorm.Model
	Fullname    string  `xml:"fullname"`
	Title       string  `xml:"title"`
	Description float64 `xml:"description"`
	Quant       int     `xml:"quant"`
	Index       string  `xml:"index"`
	Change      float64 `xml:"change"`
}
