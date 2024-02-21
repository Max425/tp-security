# Безопасность интернет - приложений
ДЗ по курсу "Безопасность интернет - приложений"

[Описание дз](https://docs.google.com/document/d/1QaQ-Nc_eE4dBKZwQbA4E2o8pOJ3CktgsKDAn375iY24/edit?usp=sharing)

Студент: Сиканов Максим `Вариант 2`

Группа WEB-31

## Запуск 
```bash
docker-compose up --build
```

При первом запуске будет ошибка из-за инициализации БД. Решение - перезапустить)

## Api

после запуска можно подергать апи в свагере http://localhost:8000/swagger/index.html#/

### Пример работы сканера 
```bash
curl -X 'GET' \
'http://localhost:8000/scan/a51a325c-4457-436c-8a75-7dd640c839d8' \
-H 'accept: application/json'
```
P. s. `a51a325c-4457-436c-8a75-7dd640c839d8` - uid уже проксированного запроса


Response body
```json
{
  "status": 200,
  "message": "success",
  "payload": [
    "Параметр 'Accept' уязвим для SQL инъекций (двойная кавычка)\n",
    "Параметр 'Accept' уязвим для SQL инъекций (одинарная кавычка)\n",
    "Параметр 'User-Agent' уязвим для SQL инъекций (одинарная кавычка)\n",
    "Параметр 'User-Agent' уязвим для SQL инъекций (двойная кавычка)\n"
  ]
}
```