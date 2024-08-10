# Описание GoStorage v0.11.11
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
`go get -u github.com/mondegor/go-storage@v0.11.11`

## Установка библиотеки для её локальной разработки

- Выбрать рабочую директорию, где должна быть расположена библиотека
- `mkdir go-storage && cd go-storage` // создать и перейти в директорию проекта
- `git clone git@github.com:mondegor/go-storage.git .`
- `cp .env.dist .env`

### Консольные команды используемые при разработке библиотеки

> Перед запуском консольных скриптов сервиса необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

- `mrcmd go help` - выводит список всех доступных go команд (docker версия);
- `mrcmd go-dev help` // выводит список всех доступных go-dev команд (локальная версия);
- `mrcmd go-dev fmt` // исправляет форматирование кода (gofumpt -l -w -extra ./)
- `mrcmd go-dev goimports-fix` // исправление imports, если это требуется (goimports -d -local ${GO_DEV_LOCAL_PACKAGE} ./)
- `mrcmd go-dev check` // статический анализ кода библиотеки (линтеры: govet, staticcheck, errcheck)
- `mrcmd go-dev test` // запуск тестов библиотеки
- `mrcmd golangci-lint check` // запуск линтеров для проверки кода (на основе `.golangci.yaml`)
- `mrcmd plantuml build-all` // генерирует файлы изображений из `.puml` [подробнее](https://github.com/mondegor/mrcmd-plugins/blob/master/plantuml/README.md#%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0-%D1%81-%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D0%B5%D0%B9-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0-markdown--plantuml)