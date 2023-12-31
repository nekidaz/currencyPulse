# currencyPulse Руководство

currencyPulse - тестовое задание для Халык Банка, веб-приложение с мощным API для отображения актуальных курсов валют, получаемых из Национального Банка Казахстана и хранящихся в базе данных PostgreSQL с использованием кэширования Redis для оптимизации производительности.

## Особенности

- Курсы валют в режиме реального времени от Национального Банка Казахстана.
- Автоматическое обновление данных каждые 6 часов для точной и актуальной информации.
- Эффективное кэширование с помощью Redis для улучшения времени ответа API.
- Полная интеграция с PostgreSQL для надежного хранения и получения данных.

## Требования

Для запуска HalykTZ вам потребуются следующие компоненты:

- Go (версия 1.15 или выше)
- PostgreSQL
- Redis

## Установка

1. Клонируйте репозиторий HalykTZ:

```bash
git clone https://github.com/your-username/HalykTZ.git
cd HalykTZ
```

2. Запустите приложение:

```bash
docker-compose up
```

Приложение будет доступно по адресу `http://localhost:8080`.

## Инициализация и Обновление данных |  ```main.go```

Во время инициализации HalykTZ выполняет следующие важные задачи:

1. Подключается к базе данных PostgreSQL и Redis для безупречного управления данными.
2. Выполняет необходимые миграции базы данных для создания необходимых таблиц.
3. Получает данные о курсах валют от Национального Банка Казахстана и сохраняет их в базе данных.
4. Настраивает периодическое обновление данных о курсах валют каждые 6 часов, чтобы пользователи получали последние данные.
```bash
func init() {
	// Подключение к PostgreSQL
	initializers.ConnectToPostgres()
	// Подключение к Redis
	initializers.ConnectToRedis()
	// Запуск миграции
	initializers.SyncDatabase()
	// При запуске кода парсит данные и сохраняет в базе данных
	data_fetcher.UpdateCurrencyData()
	// Период обновления
	period := time.Hour * 6
	// Запуск периодического обновления данных в фоне 
	go helpers.UpdateCurrencyDataPeriodically(period, data_fetcher.UpdateCurrencyData)
}
```
## API Endpoint'ы

API HalykTZ предоставляет следующие Endpoint'ы для доступа к данным о курсах валют:

1. `GET /rates`: Получает все доступные курсы валют.
2. `GET /rates/:code`: Получает курс валюты для определенного кода валюты.
3. `POST /update`: Вручную запускает обновление данных о курсах валют в кэше и базе данных.

## Получение и Хранение данных | ```data_fetcher.go```

Приложение получает актуальные данные о курсах валют в формате XML от RSS-ленты Национального Банка Казахстана. Затем оно обрабатывает XML и сохраняет курсы валют как в базу данных, так и в кэш Redis.

1. `FetchCurrencyData(url string) ([]byte, error)`: Это функция для получения данных о курсах валют от Национального Банка Казахстана. Она принимает URL в качестве аргумента, выполняет HTTP GET-запрос к указанному URL и получает XML-данные с курсами валют. В случае ошибки во время запроса или чтения данных, функция возвращает ошибку.


2. `UnmarshalXML(xmlData []byte) ([]models.CurrencyRate, error)`: Эта функция выполняет разбор XML-данных о курсах валют и возвращает структуры models.CurrencyRate, представляющие курсы валют. Она принимает входные данные в виде среза байтов и использует пакет encoding/xml для разбора XML. Если происходит ошибка разбора, функция возвращает ошибку.


3. `UpdateOrCreateCurrencyRate(rate *models.CurrencyRate) error`: Эта функция обновляет или создает запись о курсе валюты в базе данных PostgreSQL. Она принимает указатель на структуру `models.CurrencyRate`, представляющую курс валюты. Функция ищет существующую запись с тем же названием в базе данных, и если такая запись существует, обновляет ее значения. В противном случае, создает новую запись. Если происходит ошибка при обновлении или создании записи, функция возвращает ошибку.


4. `UpdateCurrencyData() error`: Эта функция обновляет данные о курсах валют от Национального Банка Казахстана и сохраняет их в базе данных PostgreSQL. Она вызывает `FetchCurrencyData` для получения актуальных XML-данных о курсах валют. Затем она вызывает `UnmarshalXML` для разбора XML и получения структур `models.CurrencyRate`. Далее, она использует функцию `UpdateOrCreateCurrencyRate` для обновления или создания записей о курсах валют в базе данных. Наконец, функция сохраняет данные о курсах валют в кэше Redis, вызывая функцию `saveCurrencyRatesToRedis`. Если происходит ошибка при получении, разборе, обновлении или сохранении данных, функция возвращает ошибку.
  

5. `saveCurrencyRatesToRedis(currencyRates []models.CurrencyRate) error`: Эта функция сохраняет данные о курсах валют в кэше Redis. Она принимает срез структур `models.CurrencyRate` и преобразует их в формат JSON. Затем она использует пакет `github.com/go-redis/redis/v8` для установки данных в Redis с ключом "currency_rates". Данные в Redis хранятся в течение 6 часов. Если происходит ошибка при маршалинге JSON или сохранении данных в Redis, функция возвращает ошибку.
