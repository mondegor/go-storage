# Описание GoStorage v0.10.1
Этот репозиторий содержит описание библиотеки GoStorage.

## Статус библиотеки
Библиотека находится в стадии разработки.

## Описание библиотеки
Библиотека для работы с хранилищами данных.
На данный момент реализованы адаптеры для следующих клиентов:
- postgres (`pgx/v5`);
- rabbitmq (`amqp-go/v1`);
- redis (`go-redis/v9` + `redislock/v0.9`);
- `S3 minio` + `FileProvider`;
- `Native File System` + `FileProvider`;

## Подключение библиотеки
`go get github.com/mondegor/go-storage`