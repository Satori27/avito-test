
## Функциональность

В API доступны и работают все ручки описанные в openapi.yml. Форматы и параметры запросов совпадают с форматами и параметрами запросов описанных в openapi.yml

## Запуск проекта

Для запуска приложения используйте Docker Compose.

1. **Клонируйте репозиторий:**

    ```sh
    git clone https://github.com/Satori27/avito-test.git
    cd avito-test
    ```

2. **Запустите приложение с помощью Docker Compose:**

    ```sh
    docker compose --env-file ./.env up
    ```

Команда запустит контейнер backend. Приложение будет доступно по адресу:

     http://localhost:9346/api/ping

В самом контейнере открыт порт `:8080`


## Примеры запросов

### Создание тендера (POST)

**Пример запроса на создание тендера:**

```sh
curl -X POST http://localhost:9346/api/tenders/new \
-H "Content-Type: application/json" \
-d '{
  "name": "string",
  "description": "string",
  "serviceType": "Construction",
  "organizationId": "550e8400-e29b-41d4-a716-446655440020",
  "creatorUsername": "user1"
}'
```

**Ожидаемый ответ от сервера:**
```
{"id":"e2939028-44aa-4002-9d7c-b09eb917ce65","version":1,"createdAt":"2024-09-15T20:46:19Z","name":"string","description":"string","serviceType":"Construction","status":"Created"}
```
