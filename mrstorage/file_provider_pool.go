package mrstorage

import (
	"fmt"
)

type (
	FileProviderPool struct {
		providers providerMap
	}

	providerMap map[string]FileProviderAPI
)

func NewFileProviderPool() *FileProviderPool {
	return &FileProviderPool{
		providers: make(providerMap, 0),
	}
}

func (p *FileProviderPool) Register(name string, provider FileProviderAPI) error {
	if _, ok := p.providers[name]; ok {
		return fmt.Errorf("file provider '%s' is already registered", name)
	}

	p.providers[name] = provider

	return nil
}

func (p *FileProviderPool) Provider(name string) (FileProviderAPI, error) {
	if provider, ok := p.providers[name]; ok {
		return provider, nil
	}

	return nil, fmt.Errorf("file provider '%s' is not registered", name)
}
