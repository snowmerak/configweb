package info

import "context"

type Data struct {
	data map[string]any
}

func New() *Data {
	return &Data{data: make(map[string]any)}
}

func With(data map[string]any) *Data {
	return &Data{data: data}
}

func (s *Data) Get() map[string]any {
	return s.data
}

type Provider interface {
	Get(ctx context.Context) (*Data, error)
	Set(ctx context.Context, data *Data) error
}
