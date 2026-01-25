package mrstorage

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
)

type (
	// FileProviderPool - пул файловых провайдеров, позволяет
	// хранить файловые провайдеры предназначенные для различных целей.
	FileProviderPool struct {
		providers providerMap
	}

	providerMap map[string]FileProvider
)

// NewFileProviderPool - создаёт объект FileProviderPool.
func NewFileProviderPool() *FileProviderPool {
	return &FileProviderPool{
		providers: make(providerMap),
	}
}

// Register - регистрирует провайдера по его имени.
func (p *FileProviderPool) Register(name string, provider FileProvider) error {
	if _, ok := p.providers[name]; ok {
		return errors.NewInternalError(
			"file provider is already registered",
			"name", name,
		)
	}

	p.providers[name] = provider

	return nil
}

// ProviderAPI - возвращает API по имени провайдера или ошибку, если он не был найден.
func (p *FileProviderPool) ProviderAPI(name string) (FileProviderAPI, error) {
	if provider, ok := p.providers[name]; ok {
		return provider, nil
	}

	return nil, errors.NewInternalError(
		"file provider is not registered",
		"name", name,
	)
}

// Ping - сообщает, установлено ли соединение и является ли оно стабильным для всех зарегистрированных провайдеров.
func (p *FileProviderPool) Ping(ctx context.Context) error {
	for name, provider := range p.providers {
		if err := provider.Ping(ctx); err != nil {
			return ErrSystemFileProviderPingError.Wrap(err, "provider", name)
		}
	}

	return nil
}

// Close - закрывает текущие соединения всех зарегистрированных провайдеров.
func (p *FileProviderPool) Close() error {
	var errs []error

	for name, provider := range p.providers {
		if providerErr := provider.Close(); providerErr != nil {
			errs = append(
				errs,
				errors.ErrSystemStorageFailedToClose.Wrap(providerErr, "source_provider", name),
			)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
