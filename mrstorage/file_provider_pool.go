package mrstorage

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
)

type (
	// FileProviderPool - пул файловых провайдеров.
	// Позволяет хранить и управлять файловыми провайдерами для различных целей.
	// Каждый провайдер регистрируется под уникальным именем.
	FileProviderPool struct {
		providers providerMap
	}

	// providerMap - внутренняя карта провайдеров (имя -> провайдер).
	providerMap map[string]FileProvider
)

// NewFileProviderPool - создаёт пустой пул файловых провайдеров.
func NewFileProviderPool() *FileProviderPool {
	return &FileProviderPool{
		providers: make(providerMap),
	}
}

// Register - регистрирует файловый провайдер под указанным именем.
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

// ProviderAPI - возвращает API файлового провайдера по его имени.
func (p *FileProviderPool) ProviderAPI(name string) (FileProviderAPI, error) {
	if provider, ok := p.providers[name]; ok {
		return provider, nil
	}

	return nil, errors.NewInternalError(
		"file provider is not registered",
		"name", name,
	)
}

// Ping - проверяет работоспособность всех зарегистрированных файловых провайдеров.
func (p *FileProviderPool) Ping(ctx context.Context) error {
	for name, provider := range p.providers {
		if err := provider.Ping(ctx); err != nil {
			return ErrSystemFileProviderPingError.Wrap(err, "provider", name)
		}
	}

	return nil
}

// Close - закрывает соединения всех зарегистрированных файловых провайдеров.
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
