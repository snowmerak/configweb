package provider

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/snowmerak/configweb/info"
	jsonProvider "github.com/snowmerak/configweb/info/provider/json"
	yamlProvider "github.com/snowmerak/configweb/info/provider/yaml"
)

type Member struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Location string `yaml:"location"`
}

type Set struct {
	Members []*Member `yaml:"providers"`

	path      string
	providers map[string]info.Provider
}

func makeProvider(m *Member) (info.Provider, error) {
	switch strings.ToLower(m.Type) {
	case "json":
		return jsonProvider.New(m.Location), nil
	case "yaml":
		return yamlProvider.New(m.Location), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", m.Type)
	}
}

func (s *Set) To(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	if err = encoder.Encode(s); err != nil {
		return fmt.Errorf("failed to encode yaml: %w", err)
	}

	return nil
}

func From(path string) (*Set, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	s := &Set{
		Members:   make([]*Member, 0),
		path:      path,
		providers: make(map[string]info.Provider),
	}

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(s); err != nil {
		return nil, fmt.Errorf("failed to decode yaml: %w", err)
	}

	for _, m := range s.Members {
		p, err := makeProvider(m)
		if err != nil {
			return nil, fmt.Errorf("failed to make provider: %w", err)
		}

		s.providers[m.Name] = p
	}

	return s, nil
}

func (s *Set) Get(name string) (info.Provider, error) {
	p, ok := s.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

	return p, nil
}
