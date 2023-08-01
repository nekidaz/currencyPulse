# Руководство по запуску проекта 
## Сначала клонируем проект 
```
git clone https://github.com/nekidaz/HalykBankTZ.git
```

## Открываем папку с проектом 
```
cd HalykBankTZ
```
## Запускаем Docker-compose
```bash
 docker-compose up
```

## Эндпоиниты для проверки

### Получить все валюты 
```bash
http://localhost:8080/rates
```
### Получить валюту по коду  
```bash
GET: http://localhost:8080/rates/:code 
```
### Получить валюту по коду - Например  
```bash
GET: http://localhost:8080/rates/usd
```
### Обновить данные в кэше и в бд
```bash
POST: http://localhost:8080/update
```


