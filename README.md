# Сервис динамического сегментирования пользователей

## Сборка проекта
```bash
docker compose --project-directory . -f deployments/docker-compose.yml build
```

## Запуск проекта

```bash
docker compose --project-directory . -f deployments/docker-compose.yml up
```

## Описание  api
Файл openapi: [openapi.yml](./api/openapi.yaml)

### Запрос создания сегмента

```bash
curl -X 'POST' \
  'http://localhost:8082/segment' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "slug": "AVITO_VOICE_MESSAGES"
}'
```

### Запрос удаления сегмента
```bash
curl -X 'DELETE' \
  'http://localhost:8082/segment' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "slug": "AVITO_VOICE_MESSAGES"
}'
```

### Запрос добавления пользователя в сегмент
```bash
curl -X 'POST' \
  'http://localhost:8082/user' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "user_id": 1000,
  "add_segments": [
    "AVITO_VOICE_MESSAGES"
  ],
  "delete_segments": [
    "AVITO_VOICE_MESSAGES_2"
  ]
}'
```

### Запрос получения активных сегментов пользователя
```bash
curl -X 'GET' \
  'http://localhost:8082/user/10000' \
  -H 'accept: application/json'
```