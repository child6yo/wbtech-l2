# HTTP-сервер «Календарь»

Простой сервис для управления CRUD-операциями вокруг событий (создание, обновление, удаление, выдача за периоды времени).

## Запуск сервиса

1. Склонируйте репозиторий через `git clone`.
2. Запустите сервис через `go run cmd/main.go -port PORT`, где PORT - ваше кастомное значение.

## API

Cтатус-коды:

    200 OK для успешных запросов;
    400 для ошибок ввода (например, некорректный date);
    503 для ошибок бизнес-логики (например, попытка удалить несуществующее событие);
    500 для прочих ошибок.

#### POST /create_event
-> создает новое событие.

**Request body**
```
{
    "user_id": "user1",
    "event": {
        "date": "2025-02-15",
        "event": "test"
    }
}
```

#### POST /update_event
-> обновляет существующее событие.

**Request body**
```
{
    "user_id": "user1",
    "event": {
        "id": "6c8b7c3f-2310-4d8b-bcd0-c62397253136",
        "date": "2025-02-15",
        "event": "emm3"
    }
}
```

#### POST /delete_event
`/delete_event?user_id=USER_ID&&id=EVENT_ID` -> удаляет событие.
#### GET /events_for_day
`/events_for_day?user_id=USER_ID&&date=YYYY-MM-DD` -> возвращает события за день.
#### GET /events_for_week
`/events_for_day?user_id=USER_ID&&date=YYYY-MM-DD` -> возвращает события на 7 дней, начиная с переданного дня.
#### GET /events_for_month
`/events_for_day?user_id=USER_ID&&date=YYYY-MM-DD` -> возвращает события на месяц, переданный в MM, DD может быть любой.