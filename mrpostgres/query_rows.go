package mrpostgres

import (
	"github.com/jackc/pgx/v5"
)

type (
	queryRows struct {
		rows pgx.Rows
	}
)

// Next - переходит к следующей записи.
// Если следующей записи не существует, то возвращает false.
func (qr *queryRows) Next() bool {
	return qr.rows.Next()
}

// Scan - извлекает значения полученные запросом
// и присваивает их переданным переменным.
func (qr *queryRows) Scan(dest ...any) error {
	if err := qr.rows.Scan(dest...); err != nil {
		return wrapErrorFetchDataFailed(err)
	}

	return nil
}

// Err - возвращает ошибку в результате обработки списка переданного в запросе.
func (qr *queryRows) Err() error {
	if err := qr.rows.Err(); err != nil {
		return wrapErrorFetchDataFailed(err)
	}

	return nil
}

// Close - закрывает запрос вернувший список записей.
func (qr *queryRows) Close() {
	qr.rows.Close()
}
