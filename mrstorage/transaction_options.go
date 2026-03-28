package mrstorage

import (
	"github.com/mondegor/go-storage/mrstorage/txisolevel"
)

type (
	// TxOption - настройка объекта TxOptions.
	TxOption func(o *TxOptions)

	// TxOptions - настройки для создания транзакции.
	TxOptions struct {
		IsoLevel txisolevel.Enum
	}
)

// WithTxIsoLevel - устанавливает уровень изоляции для транзакции.
func WithTxIsoLevel(value txisolevel.Enum) TxOption {
	return func(o *TxOptions) {
		o.IsoLevel = value
	}
}

// WithTxIsoLevelRepeatableRead - устанавливает уровень изоляции RepeatableRead для транзакции.
func WithTxIsoLevelRepeatableRead() TxOption {
	return func(o *TxOptions) {
		o.IsoLevel = txisolevel.RepeatableRead
	}
}

// WithTxIsoLevelSerializable - устанавливает уровень изоляции Serializable для транзакции.
func WithTxIsoLevelSerializable() TxOption {
	return func(o *TxOptions) {
		o.IsoLevel = txisolevel.Serializable
	}
}
