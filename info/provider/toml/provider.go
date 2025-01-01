package toml

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/sync/singleflight"

	"github.com/snowmerak/configweb/info"
)

var _ info.Provider = &Provider{}

type Provider struct {
	path   string
	info   atomic.Pointer[info.Data]
	worker singleflight.Group
}

func New(path string) *Provider {
	return &Provider{
		path:   path,
		info:   atomic.Pointer[info.Data]{},
		worker: singleflight.Group{},
	}
}

func (p *Provider) Get(_ context.Context) (i *info.Data, err error) {
	if v := p.info.Load(); v != nil {
		return v, nil
	}

	r, err, _ := p.worker.Do("get", func() (interface{}, error) {
		f, err := os.Open(p.path)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		decoder := toml.NewDecoder(f)
		data := make(map[string]any)
		if err = decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to decode toml: %w", err)
		}

		i = info.With(data)

		p.info.Store(i)

		return i, nil
	})
	if err != nil {
		return nil, err
	}

	i = r.(*info.Data)

	return i, nil
}

func (p *Provider) Set(_ context.Context, i *info.Data) error {
	f, err := os.Create(p.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err = encoder.Encode(i.Get()); err != nil {
		return fmt.Errorf("failed to encode toml: %w", err)
	}

	return nil
}
