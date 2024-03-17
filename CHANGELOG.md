# GoStorage Changelog
Все изменения библиотеки GoStorage будут документироваться на этой странице.

## 2024-03-17
### Added
- Добавлен новый метод `FilterEqualUUID` для интерфейса `mrstorage.SqlBuilderWhere`,
  а также его реализация для `postgres`;

## 2024-03-16
### Changed
- Доработан `mrsql.FieldsForUpdate`, добавлены новые поддерживаемые типы:
  - Slice, Int, Int8, Int16, Uint, Uint8, Uint16, Uint32, Uint64;

## 2024-03-14
### Changed
- Переименовано `ModifiedAt` -> `UpdatedAt`;
- Из `SqlBuilderSelect` выделен интерфейс `SqlBuilderCondition`;
- Метод `SqlBuilderOrderBy.WrapWithDefault()` преобразован в `SqlBuilderOrderBy.DefaultField()`;
- Доработан пример формирования порядка следования;
- Перенесена и доработана следующая логика:
  - `mrsql.BuilderSelect` -> `mrpostgres.SqlBuilderSelect`;
  - `mrsql.BuilderUpdate` -> `mrpostgres.SqlBuilderUpdate`;

### Fixed
- Добавлена дополнительная проверка в `mrminio.DownloadFile` на случай если файл не существует;

## 2024-02-01
### Fixed
- добавлен вызов `Caller(1)` при записи ошибок в лог в некоторых методах;

## 2024-01-30
### Changed
- Для методов `Connect()` адаптеров добавлен параметр `ctx context.Context`;
- Внедрён новый интерфейс логгера, добавлен режим трассировки.

## 2024-01-25
### Added
- Добавлен новый метод `DownloadFile()` для интерфейса `mrstorage.FileProviderAPI`.
  Также добавлена его реализация для провайдеров `mrfilestorage` и `mrminio`;
- Добавлены следующие вспомогательные функции `mrentity.FileMetaToInfo()` и `mrentity.ImageMetaToInfo()`;

### Changed
- Переименовано:
  - ConvertFileMetaToInfo -> FileMetaToInfoPointer
  - ConvertImageMetaToInfo -> ImageMetaToInfoPointer;

### Removed
- Удалён адаптер `mrredsync` т.к. один из его компонентов использует `MPL-2.0 license`;

## 2024-01-22
### Changed
- Обновлены зависимости библиотеки;
- `FactoryErrInternalWithData` было заменено на `FactoryErrInternal.WithAttr(...)`;
- В `NewSqlBuilderWhere` добавлено выделение памяти с помощью `buf.Grow` при формировании некоторых условий;

## 2024-01-18
### Changed
- Переработан метод `SqlBuilderWhere.FilterAnyOf` для того, чтобы избавиться от зависимости `github.com/lib/pq`;
- Теперь сборка условий в `SqlBuilderWhere` происходит с использованием `strings.Builder`;

## 2024-01-16
### Added
- Добавлена ошибка `mrfilestorage.FactoryErrInvalidPath`;
- Добавлены `mrentity.FileMeta` и `mrentity.ImageMeta` для хранения их в БД в виде json.
  Также добавлены функции `ConvertFileMetaToInfo` и `ConvertImageMetaToInfo` для преобразования
  данных в формат для пользователя;

### Changed
- Тип `mrtype.NullableBool` заменён на `*bool`;
- Обновлены зависимости библиотеки;
- Устранена путаница с пакетами `path` и `path/filepath`, теперь используется
  только один из них в рамках одного файла;

### Removed
- Удалён интерфейс `mrstorage.ExtFileProviderAPI` и метод у файловых провайдеров `WithBaseDir`;
- Удалён `mrentity.ZeronullTime`;

## 2023-12-13
### Changed
- В `mrredislock` и `mrredsync` логирование вызова методов теперь происходив в самом начале;

## 2023-12-11
### Added
- В пакетах `mrredislock` и `mrredsync`, добавлено обёртывание ошибок и логирование исполнения команд;

### Changed
- Обновлены зависимости библиотеки;

## 2023-12-10
### Added
- Добавлен интерфейс `ExtFileProviderAPI`, в котором метод `WithBaseDir` позволяет задать
  постоянный префикс ко всем именам файлов используемых в интерфейсе `FileProviderAPI`.
  `ExtFileProviderAPI` интерфейс следует использовать только при инициализации системы;
- Добавлено отладочное логирование вызовов команд в `mrfilestorage` и `mrminio`;

### Changed
- Доработана логика копирования объектов в `BuilderPart.WithPrefix`, `BuilderPart.Param`;
- В `mrminio.ConnAdapter` добавлен флаг `createBuckets`, а из `mrminio.ConnAdapter.InitBucket` этот флаг удалён;
- Переработан пакет `mrfilestorage`, добавлена абстракция `FileSystem`, которая инициализирует базовые директории
  для хранения файлов в рамках файловой системы. Добавлено обёртывание ошибок, поддержка `ExtFileProviderAPI`;
- Доработан пакет `mrminio`, добавлено обёртывание ошибок, поддержка `ExtFileProviderAPI`;

## 2023-12-06
### Changed
- Теперь ошибка `mrcore.FactoryErrStorageRowsNotAffected` формируется для запросов типа INSERT, UPDATE и DELETE;

### Fixed
- В методе `SqlBuilderSet.Fields` если параметр names пустой, то возвращается nil;

## 2023-12-04
### Added
- Обновлены зависимости библиотеки;

### Fixed
- Исправлена работоспособность некоторых примеров использования библиотеки;

## 2023-11-27
### Added
- Добавлен `FileProviderPool`;

## 2023-11-26
### Added
- В `SqlBuilderWhere` добавлены новые методы: `Less`, `LessOrEqual`, `Greater`, `GreaterOrEqual`;

## 2023-11-23
### Changed
- Обновлены зависимости библиотеки;

## 2023-11-20
### Changed
- В некоторых местах оптимизирована конкатенация строк (`Sprintf` заменён на нативный "+");
- Обновлены зависимости библиотеки;
- Обновлён `.editorconfig`;

## 2023-11-13
### Added
- Для `SqlBuilderOrderBy` добавлена обёртка `WrapWithDefault` для того можно было формировать сортировку по умолчанию;
- Для `BuilderUpdate` добавлен метод `SetFromEntityWith`, чтобы можно было формировать список полей для обновления не только в рамках указанной структуры;
- Добавлен интерфейс для работы с внешними файловыми хранилищами `FileProviderAPI`, реализован данный интерфейс в пакетах `mrfilestorage` и `mrminio`;
- В `SqlBuilderWhere` добавлены новые методы для фильтрации данных `FilterEqualString`, `FilterEqualBool`;

### Changed
- Для `mrminio` доработан ConnAdapter, добавлен метод `InitBucket`, который проверяет существование указанного бакета и создаёт его при необходимости;
- Доработан `EntityMetaOrderBy`, который управляет сортируемыми полями структуры через тег `sort`, его метод `CheckField` проверяет возможность по указанному полю проведение сортировки. Метод `DefaultSort` возвращает поле для сортировки по умолчанию;
- Доработан `EntityMetaUpdate`, который управляет обновляемыми полями структуры через тег `upd`, его метод FieldsForUpdate формирует запрос на обновление только заполненных полей;
- Переименованы некоторые переменные и функции (типа Id -> ID) в соответствии с code style языка go;
- Обновлены зависимости библиотеки;
- Все файлы библиотеки были пропущены через `gofmt`;

## 2023-11-01
### Added
- Добавлен пакет `mrsql` (`SqlBuilder`) для генерации SQL фрагментов, которые можно подключать к основному SQL запросу;
- Добавлена обработка тегов `fieldTagFreeUpdate`, `fieldTagSortByField`;
- Добавлены новые сущности:
    - `mrentity.ListSorter` + `mrreq.ParseListSorter`;
    - `mrentity.ListPager` + `mrreq.ParseListPager`;
    - `mrentity.RangeInt64` + `mrreq.ParseRangeInt64`;
    - `mrentity.SortDirection`;

### Changed
- Обновлены зависимости библиотеки;
- Переименована сущность `mrstorage.File` -> `mrentity.File`;
- Переработана `FilledFieldsToUpdate` и перенос кода в `NewEntityMetaUpdate`;

### Removed
- После внедрения `SqlBuilder` была удалена функциональность связанная с `DbSqlizer`:
    - в пакете `mrstorage` удалены `SqQuery`, `SqQueryRow`, `SqExec`;
    - в `exec_helper_sqlizer.go` удалены `sqQuery`, `sqQueryRow`, `sqExec`, `parseSql`;
    - в сущности `mrpostgres.Transaction` удалены `SqQuery`, `SqQueryRow`, `SqExec`;
- Удалена ошибка не используемая `FactoryErrInternalListOfFieldsIsEmpty`;

## 2023-10-08
### Added
- В пакет `mrpostgres` в `Options` добавлен `AfterConnectFunc`;

### Changed
- Обновлены зависимости библиотеки;
- Обработка ошибок приведена к более компактному виду;

## 2023-09-20
### Added
- Добавлены интерфейсы для БД: `DbConn`, `DbTransaction`, `DbQuery` и др.;

### Changed
- Заменены `tabs` на пробелы в коде;
- Переработан адаптер для `postgres`, под новые интерфейсы, добавлена сущность транзакции;
- Обновлены зависимости библиотеки;

## 2023-09-16
### Added
- Реализован интерфейс `mrcore.locker` на базе `redsync` и `redislock`;
- Добавлен интерфейс `mrstorage.FileProvider` для работы с файлами;
- Добавлена поддержка работы с `minio`;

### Changed
- Из адаптера `redis` удалён `redsync`;

## 2023-09-13
### Changed
- Обновлены зависимости библиотеки;
- Переименованы `Connection` -> `ConnAdapter`;
- Добавлен интерфейс `mrstorage.Sqlizer`, для того чтобы снять зависимость от `squirrel`;

## 2023-09-12
### Changed
- Все часто используемые ошибки теперь подключаются из пакета `mrcore`;
- Формат глобальных `const`, `type`, `var` приведён к общему виду;
- Некоторые названия ошибок переименованы для поддержки обновлённой версии `go-webcore`;

## 2023-09-11
### Changed
- Обновлены зависимости библиотеки;
- Доработаны методы управления соединениями хранилищ;

## 2023-09-10
### Changed
- Обновлены зависимости библиотеки;
- Доработаны методы управления соединениями хранилищ;