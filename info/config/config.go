package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"

	"github.com/snowmerak/configweb/info"
	"github.com/snowmerak/configweb/info/provider"
	"github.com/snowmerak/configweb/pair"
)

type Config struct {
	data map[string]any
	keys []*pair.Pair[[]string, string]
	set  *provider.Set
}

func iterateData(data map[string]any) (map[string]any, []*pair.Pair[[]string, string]) {
	queue := make([]*pair.Pair[[]string, map[string]any], 0, 1)
	queue = append(queue, pair.New([]string{}, data))

	keys := make([]*pair.Pair[[]string, string], 0)
	cloned := make(map[string]any)

	for len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]

		m := key.Second()
		for k, v := range m {
			insert := true
			kl := append(key.First(), k)
			switch p := v.(type) {
			case map[string]any:
				queue = append(queue, pair.New(append(key.First(), k), p))
				insert = false
			case string:
				if strings.HasPrefix(p, "$") {
					keys = append(keys, pair.New(kl, p))
				}
			}

			if insert {
				c := cloned
				for _, k := range kl[:len(kl)-1] {
					if _, ok := c[k]; !ok {
						c[k] = make(map[string]any)
					}

					c = c[k].(map[string]any)
				}
				c[kl[len(kl)-1]] = v
			}
		}
	}

	return cloned, keys
}

func New(data *info.Data, providers *provider.Set) *Config {
	d, k := iterateData(data.Get())
	return &Config{
		data: d,
		keys: k,
		set:  providers,
	}
}

type BuildTarget string

const (
	TargetYAML BuildTarget = "yaml"
	TargetJSON BuildTarget = "json"
	TargetTOML BuildTarget = "toml"
)

func (c *Config) Build(ctx context.Context, target BuildTarget) ([]byte, error) {
	for _, k := range c.keys {
		pv, err := c.set.Get(k.Second())
		if err != nil {
			return nil, fmt.Errorf("failed to get value: %w", err)
		}

		cm := c.data
		for _, k := range k.First()[:len(k.First())-1] {
			cm = cm[k].(map[string]any)
		}

		v, err := pv.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get value: %w", err)
		}

		cm[k.First()[len(k.First())-1]] = v
	}

	switch target {
	case TargetYAML:
		return yaml.Marshal(c.data)
	case TargetJSON:
		return json.Marshal(c.data)
	case TargetTOML:
		return toml.Marshal(c.data)
	}

	return nil, fmt.Errorf("unknown target: %s", target)
}
