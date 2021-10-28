# Distributed queue (XK6-SQS)

Модуль для нагрузочного тестирования SQS через K6

## Сборка
Для сборки `k6` с модулем xk6-sqs, убедитесь, что у вас установлено:
- [Go toolchain](https://golang.org/doc/install)
- Git

1. Установите `xk6`
```bash
  go install github.com/k6io/xk6/cmd/xk6@latest
```
2. Соберите бинарь `k6` с модулем `xk6-sqs`
```bash
  xk6 build  --with github.com/tarantool/xk6-sqs
```
или локально
```bash
  xk6 build  --with github.com/tarantool/xk6-sqs=.
```

## Использование

- Конфигурация модуля
    * url (optional): URL SQS
    * user_id (optional): ID пользователя SQS
    * access_key_id (optional): ID ключа SQS
    * secret_access_key (optional): Ключ SQS
```javascript
import sqs from 'k6/x/sqs';
const client = new sqs.Client({
    url: "http://localhost:8081",
    user_id: "1234567890",
    access_key_id: "ACCESS-KEY-ID",
    secret_access_key: "SECRET-ACCESS-KEY"
});
```

- Создание очереди
```javascript
client.createQueue("test")
```

- Удаление очереди
```javascript
client.deleteQueue("test")
```

* Добавить сообщение в очередь
    -  queue_name: Название очереди
    -  message_body: Сообщение
    -  delay_seconds (optional): Задержка перед обработкой
```javascript
client.sendMessage({
    queue_name: "test",
    message_body: "Message"
})
```
* Добавить партию сообщений в очередь
    - queue_name: Название очереди
    - entries: Массив объектов
        - id: ID сообщения
        - message_body: Сообщение
```javascript
client.sendMessageBatch({
    queue_name: "test",
    entries: [{
        id: "1",
        message_body: "Message1"
    },
    {
        id: "2",
        message_body: "Message2"
    }]
})
```

* Получить сообщения
    - queue_name: Название очереди
    - max_number_of_messages (optional): Кол-во сообщений
    - wait_time_seconds (optional): Время ожидания сообщений
    - visibility_timeout (optional): Время скрытия полученных сообщений
```javascript
client.receiveMessage({
    queue_name: "test",
    max_number_of_messages: 5,
	wait_time_seconds: 5,
    visibility_timeout: 10
})
```

* Удалить сообщение
    - queue_name: Название очереди
    - receipt_handle: ReceiptHandleId сообщения
```javascript
client.deleteMessage({
    queue_name: "test",
    receipt_handle: "ReceiptHandle"
})
```

* Удалить партию сообщений
    - queue_name: Название очереди
    - entries: Массив объектов
        - id: ID сообщения
        - receipt_handle: ReceiptHandleId сообщения
```javascript
client.deleteMessageBatch({
    queue_name: "test",
    entries: [
        {
            id: "MessageId",
            receipt_handle: "ReceiptHandle"
        }
    ]
})
```
