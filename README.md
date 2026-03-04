# url-shortener
тестовое задание в ozon

## Укорачиватель ссылок

Запуск
```
docker compose up --build
```

Пример запроса
```
curl -i -X POST "http://localhost:8080/url" \
  -u myuser:mypass \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","alias":""}'
```

Пример ответа
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 04 Mar 2026 01:48:18 GMT
Content-Length: 37

{"status":"OK","alias":"ACDEFGHIJK"}
```

Тесты
```
go test ./tests -v
```
