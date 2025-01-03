package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/snowmerak/configweb/info"
	"github.com/snowmerak/configweb/info/config"
	"github.com/snowmerak/configweb/info/provider"
	jsonProvider "github.com/snowmerak/configweb/info/provider/json"
	tomlProvider "github.com/snowmerak/configweb/info/provider/toml"
	yamlProvider "github.com/snowmerak/configweb/info/provider/yaml"
)

var joinCommand = &cli.Command{
	Name:      "join",
	Aliases:   []string{"j"},
	Usage:     "Join infra and package config",
	Args:      true,
	ArgsUsage: "<name>",
	Action:    joinCommandAction,
}

var joinCommandAction = func(context *cli.Context) error {
	cfgSet, err := provider.From(filepath.Join(".", ConfigSetFile))
	if err != nil {
		return fmt.Errorf("failed to get config set: %w", err)
	}

	args := context.Args().First()
	if args == "" {
		return errors.New("package name is required")
	}

	pkgPv := info.Provider(nil)
	pkgPath := filepath.Join(".", PackageConfigDir, args)
	switch filepath.Ext(strings.ToLower(pkgPath)) {
	case ".json":
		pkgPv = jsonProvider.New(pkgPath)
	case ".yaml", ".yml":
		pkgPv = yamlProvider.New(pkgPath)
	case ".toml":
		pkgPv = tomlProvider.New(pkgPath)
	default:
		return errors.New("unknown package type")
	}

	pkgData, err := pkgPv.Get(context.Context)
	if err != nil {
		return fmt.Errorf("failed to get package data: %w", err)
	}

	cfg := config.New(pkgData, cfgSet)
	data, err := cfg.Build(context.Context, config.TargetYAML)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	base := filepath.Base(args)
	ext := filepath.Ext(base)
	joined := strings.TrimSuffix(base, ext)
	dest, err := os.Create(filepath.Join(".", JoinedConfigDir, joined+".yaml"))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dest.Close()

	if _, err = dest.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}
