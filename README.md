# Описание GoStorage v0.17.1
Этот репозиторий содержит описание библиотеки GoStorage.

## Статус библиотеки
Библиотека находится в стадии разработки.

## Описание библиотеки
Библиотека для работы с хранилищами данных.
На данный момент реализованы адаптеры для следующих клиентов:
- rabbitmq - `mrrabbitmq` (`amqp091-go`);
- redis - `mrredis` (`go-redis/v9`) + распределённые блокировки `mrredis/locker` (`redislock`);
- S3 minio - `mrminio` (`minio-go/v7`) + `FileProvider`;
- Native File System - `mrfilestorage` + `FileProvider`;

### Пакеты библиотеки

#### `mrrabbitmq` - RabbitMQ (`amqp091-go`)
`ConnAdapter` управляет соединением с сервером RabbitMQ (AMQP 0.9.1):
- `New()` - создаёт адаптер без активного соединения;
- `Connect(ctx, Options)` - подключается по схеме `amqp://User:Password@Host:Port/`;
- `Cli()` - возвращает нативное соединение `*amqp.Connection`, через которое создаются каналы
  для публикации и потребления сообщений;
- `Close()` - закрывает соединение.

#### `mrredis` - Redis (`go-redis/v9`)
`ConnAdapter` управляет соединением с Redis:
- `New(tracer)` - создаёт адаптер без активного соединения;
- `Connect(ctx, Options)` - подключается по `DSN` либо по `Host`/`Port`, пароль (`Password`)
  переопределяет пароль из DSN, таймауты чтения/записи (`ReadTimeout`/`WriteTimeout`) по умолчанию 5s;
- `Ping(ctx)`, `Close()` - проверка и закрытие соединения;
- `Cli()` - возвращает нативный клиент `redis.UniversalClient`;
- `GetStruct` / `SetStruct` / `Delete` - чтение и запись структур по ключу (с временем жизни ключа),
  удаление одного или нескольких ключей.

#### `mrredis/locker` - распределённые блокировки (`bsm/redislock`)
`Adapter` реализует интерфейс `mrlock.Locker` из GoCore поверх Redis:
- `NewAdapter(conn, logger, tracer)` - создаёт адаптер на основе `redis.UniversalClient`
  (например, полученного через `mrredis.ConnAdapter.Cli()`);
- `Lock(ctx, key)` / `LockWithExpiry(ctx, key, expiry)` - захватывают блокировку ключа
  (со временем жизни по умолчанию или заданным) и возвращают функцию её освобождения.

#### `mrminio` - S3-совместимое хранилище MinIO (`minio-go/v7`)
`ConnAdapter` управляет соединением с MinIO:
- `New(createBuckets, mimeTypes, tracer)` - создаёт адаптер; `createBuckets` разрешает автоматически
  создавать отсутствующие бакеты, `mimeTypes` используется для определения типа контента;
- `Connect(ctx, Options)` - подключается по `DSN` (`host:port`) либо по `Host`/`Port`,
  с учётом `UseSSL` и учётных данных `User`/`Password`;
- `Ping(ctx)`, `Close()`, `Cli()` (нативный `*minio.Client`);
- `InitBucket(ctx, bucketName)` - проверяет наличие бакета и при необходимости создаёт его.

`FileProvider` (`NewFileProvider(conn, bucketName)`) реализует `mrstorage.FileProvider` для одного бакета:
- `Info` - метаинформация файла (размер, тип контента, даты, оригинальное имя);
- `Download` / `DownloadFile` - содержимое файла вместе с метаинформацией или только `io.ReadCloser`;
- `Upload` - сохраняет файл, определяя `ContentType` по расширению (если он не указан)
  и записывая оригинальное имя файла в `Content-Disposition`;
- `Remove` - удаляет файл;
- `Ping`, `Close`, `Cli` - унаследованы от встроенного `ConnAdapter`.

#### `mrfilestorage` - файловое хранилище на локальной файловой системе
- `FileSystem` (`New(dirMode, createDirs, mimeTypes)`) - работа с директориями: `InitRootDir`
  инициализирует (и при `createDirs` создаёт) корневую директорию, `CreateDirIfNotExists` создаёт
  вложенные директории;
- `FileProvider` (`NewFileProvider(fs, tracer, rootDir)`) реализует `mrstorage.FileProvider`
  относительно `rootDir`: `Info`, `Download` / `DownloadFile`, `Upload` (создаёт недостающие директории),
  `Remove`, `Ping` (проверяет возможность записи, создавая и сразу удаляя тестовый файл), `Close`.
  Пути короче 3 символов или содержащие `..` отклоняются с ошибкой `ErrInternalInvalidPath`.

Благодаря общему интерфейсу `mrstorage.FileProvider` провайдеры `mrminio` и `mrfilestorage`
взаимозаменяемы. В `mrredis`, `mrredis/locker`, `mrminio.FileProvider` и `mrfilestorage.FileProvider`
ошибки драйверов транслируются в ошибки GoCore (`go-core/errors`), а операции трассируются через `mrtrace`.
Примеры использования находятся в директории `examples/` (`mrminio`, `mrrabbitmq`, `mrredis`).

Абстрактные интерфейсы, которые реализуют адаптеры (`mrstorage`, `mrlock`), модели файлов (`mrmodel/media`),
а также адаптер для PostgreSQL (`mrpostgres`) находятся в библиотеке
[GoCore](https://github.com/mondegor/go-core).

Также реализован вспомогательный пакет `mrtests` для тестирования приложения без необходимости
разворачивания инфраструктуры (БД и т.д.), для этого используется библиотека `testcontainers`.
Пакет разделён по сервисам, в каждом из которых есть контейнер (`Container`) и объект для работы
с ним в тестах (`Tester`):
- `mrtests/pgtest` - postgres; `pgtest.Tester` применяет миграции (`golang-migrate`),
  загружает и очищает фикстуры (`testfixtures`);
- `mrtests/redistest` - redis (запускается с паролем).

По умолчанию поднимаются публичные образы Docker Hub (`postgres:18.3-alpine3.23`, `redis:7.4.8-alpine3.21`),
а вся настройка контейнера задаётся в коде. Образ можно переопределить переменными окружения
`MRTESTS_POSTGRES_DOCKER_IMAGE` и `MRTESTS_REDIS_DOCKER_IMAGE`; итоговый образ возвращают функции
`pgtest.DockerImage()` и `redistest.DockerImage()` (пригодится при прямом вызове `NewContainer`).

## Подключение библиотеки
`go get -u github.com/mondegor/go-storage@v0.17.1`

Требуется Go 1.26 или новее.

## Установка библиотеки для её локальной разработки
- Выбрать рабочую директорию, где должна быть расположена библиотека
- `mkdir go-storage && cd go-storage` // создать и перейти в директорию проекта
- `git clone git@github.com:mondegor/go-storage.git .`
- `cp .env.dist .env`
- `mrcmd go-dev deps` // загрузка зависимостей проекта
- Для работы утилит `gofumpt`, `goimports`, `gci`, `golangci-lint`, `mockgen` необходимо запустить
  `mrcmd go-dev install-tools`. По умолчанию `gofumpt`, `goimports`, `gci`, `golangci-lint` устанавливаются
  последних версий; чтобы закрепить версию, раскомментируйте переменную `GO_DEV_TOOLS_INSTALL_*` в `.env`.
  `mockgen` в go-dev по умолчанию выключен, но в `.env.dist` для него задана версия, поэтому он тоже
  будет установлен (сейчас в библиотеке нет моков, он понадобится только при их добавлении);

### Запуск тестов
Тесты библиотеки запускаются командой `go test ./...`. Интеграционные тесты пакетов `mrtests/*`
поднимают контейнеры и требуют запущенный Docker, без него `go test ./...` завершится с ошибкой;
это же требование действует при использовании `mrtests` (контейнеры `testcontainers`) в тестах приложения.

### Консольные команды используемые при разработке библиотеки

> Перед запуском консольных скриптов библиотеки необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

- `mrcmd go-dev help` // выводит список всех доступных go-dev команд;
- `mrcmd go-dev generate` // генерирует go файлы через встроенный механизм go:generate;
- `mrcmd go-dev gofumpt-fix` // исправляет форматирование кода (`gofumpt -l -w -extra ./`);
- `mrcmd go-dev goimports-fix` // исправляет imports, если это требуется (`goimports -l -w -local ${GO_DEV_IMPORTS_LOCAL_PREFIXES}` для всех go файлов, кроме сгенерированных);
- `mrcmd go-dev gci-fix` // упорядочивает imports (`gci write --skip-generated -s standard -s default -s prefix(${GO_DEV_IMPORTS_LOCAL_PREFIXES}) .`);
- `mrcmd go-dev lint` // запускает линтеры для проверки кода (на основе `.golangci.yaml`);
- `mrcmd go-dev test` // запускает тесты библиотеки;
- `mrcmd go-dev test-report` // запускает тесты библиотеки с формированием отчёта о покрытии кода (`test-coverage-full.html`);

#### Короткий вариант выше приведённых команд (Makefile)
- `make deps` // аналог `mrcmd go-dev deps`
- `make deps-upgrade` // аналог `mrcmd go-dev get -u ./...` + `mrcmd go-dev tidy`
- `make generate` // аналог `mrcmd go-dev generate`
- `make lint` // аналог `mrcmd go-dev gofumpt-fix` + `goimports-fix` + `gci-fix` + `lint`
- `make test` // аналог `mrcmd go-dev test`
- `make test-report` // аналог `mrcmd go-dev test-report`

Дополнительные команды (Makefile.mk):
- `make check-and-fix` // generate + форматирование + lint + test;
- `make full` // `make deps` + `make check-and-fix`;
- `make full2` // `make deps-upgrade` + `make check-and-fix`;
- `make archive` // упаковывает проект в `../go-storage.tar.gz`;

> Чтобы расширить список команд, необходимо создать Makefile.mk и добавить
> туда дополнительные команды, все они будут добавлены в единый список команд make утилиты.
