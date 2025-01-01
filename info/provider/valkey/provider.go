package valkey

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/snowmerak/configweb/info"
)

const (
	clientSideCacheTTL = 5 * time.Minute
)

func makeKey(key string) string {
	return "configweb:siteinfo:" + key
}

var _ info.Provider = &Provider{}

type Provider struct {
	client valkey.Client
	key    string
	info   atomic.Pointer[info.Data]
}

func New(client valkey.Client, key string) *Provider {
	return &Provider{
		client: client,
		key:    makeKey(key),
		info:   atomic.Pointer[info.Data]{},
	}
}

func (p *Provider) Get(ctx context.Context) (*info.Data, error) {
	if v := p.info.Load(); v != nil {
		return v, nil
	}

	data, err := p.client.DoCache(ctx, p.client.B().Get().Key(p.key).Cache(), clientSideCacheTTL).AsBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	value := make(map[string]any)
	if err = json.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	i := info.With(value)

	p.info.Store(i)

	return i, nil
}

func (p *Provider) Set(ctx context.Context, data *info.Data) error {
	value := data.Get()
	parsed, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err = p.client.Do(ctx, p.client.B().Set().Key(p.key).Value(valkey.BinaryString(parsed)).Build()).Error(); err != nil {
		return fmt.Errorf("failed to set data: %w", err)
	}

	p.info.Store(data)

	return nil
}
