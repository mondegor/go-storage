# GoStorage Changelog
Все изменения библиотеки GoStorage будут документироваться на этой странице.

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