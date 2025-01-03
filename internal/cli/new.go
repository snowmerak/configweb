package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/snowmerak/configweb/info"
	"github.com/snowmerak/configweb/info/provider"
	"github.com/snowmerak/configweb/info/provider/json"
	"github.com/snowmerak/configweb/info/provider/toml"
	"github.com/snowmerak/configweb/info/provider/yaml"
)

var newCommand = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "Create a new config",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "project",
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "json-provider",
			Aliases: []string{"j"},
		},
		&cli.BoolFlag{
			Name:    "yaml-provider",
			Aliases: []string{"y"},
		},
		&cli.BoolFlag{
			Name:    "toml-provider",
			Aliases: []string{"t"},
		},
		&cli.BoolFlag{
			Name:    "package",
			Aliases: []string{"k"},
		},
	},
	Args:      true,
	ArgsUsage: "<dir>",
	Action:    newProjectAction,
}

var newProjectAction = func(context *cli.Context) error {
	args := context.Args().First()
	firstFlag := context.FlagNames()[0]
	log.Printf("Creating new %s", firstFlag)
	switch firstFlag {
	case "project", "p":
		log.Printf("Creating project in %s", args)
		if err := os.MkdirAll(args, 0777); err != nil {
			return err
		}

		ps := new(provider.Set)
		ps.Members = []*provider.Member{
			{
				Name:     "json_example",
				Type:     "json",
				Location: "./infra/json_example.json",
			},
			{
				Name:     "yaml_example",
				Type:     "yaml",
				Location: "./infra/yaml_example.yaml",
			},
			{
				Name:     "toml_example",
				Type:     "toml",
				Location: "./infra/toml_example.toml",
			},
		}

		configSetFilePath := filepath.Join(args, ConfigSetFile)
		log.Printf("Writing config set file to %s", configSetFilePath)
		if err := ps.To(configSetFilePath); err != nil {
			return fmt.Errorf("failed to write config set file: %w", err)
		}

		log.Printf("Creating infra directory")
		infraDir := filepath.Join(args, InfraConfigDir)
		if err := os.MkdirAll(infraDir, 0777); err != nil {
			return fmt.Errorf("failed to create infra directory: %w", err)
		}

		log.Printf("Creating package directory")
		packageDir := filepath.Join(args, PackageConfigDir)
		if err := os.MkdirAll(packageDir, 0777); err != nil {
			return fmt.Errorf("failed to create package directory: %w", err)
		}

		log.Printf("Creating joined directory")
		joinedDir := filepath.Join(args, JoinedConfigDir)
		if err := os.MkdirAll(joinedDir, 0777); err != nil {
			return fmt.Errorf("failed to create joined directory: %w", err)
		}

		log.Printf("Created new project in %s", args)
	case "json-provider", "j":
		log.Printf("Creating json provider in %s", args)
		if !strings.HasSuffix(strings.ToLower(args), ".json") {
			args = args + ".json"
		}
		path := filepath.Join(".", InfraConfigDir)
		if err := os.MkdirAll(path, 0777); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		path = filepath.Join(path, args)
		pv := json.New(path)
		if err := pv.Set(context.Context, info.With(map[string]any{
			"KEY":  "value",
			"LIST": []string{"a", "b", "c"},
		})); err != nil {
			return fmt.Errorf("failed to set: %w", err)
		}
		log.Printf("Created json provider in %s", path)
	case "yaml-provider", "y":
		log.Printf("Creating yaml provider in %s", args)
		if !strings.HasSuffix(strings.ToLower(args), ".yaml") || !strings.HasSuffix(strings.ToLower(args), ".yml") {
			args = args + ".yaml"
		}
		path := filepath.Join(".", InfraConfigDir)
		if err := os.MkdirAll(path, 0777); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		path = filepath.Join(path, args)
		pv := yaml.New(path)
		if err := pv.Set(context.Context, info.With(map[string]any{
			"KEY":  "value",
			"LIST": []string{"a", "b", "c"},
		})); err != nil {
			return fmt.Errorf("failed to set: %w", err)
		}
		log.Printf("Created yaml provider in %s", path)
	case "toml-provider", "t":
		log.Printf("Creating toml provider in %s", args)
		if !strings.HasSuffix(strings.ToLower(args), ".toml") {
			args = args + ".toml"
		}
		path := filepath.Join(".", InfraConfigDir)
		if err := os.MkdirAll(path, 0777); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		path = filepath.Join(path, args)
		pv := toml.New(path)
		if err := pv.Set(context.Context, info.With(map[string]any{
			"KEY":  "value",
			"LIST": []string{"a", "b", "c"},
		})); err != nil {
			return fmt.Errorf("failed to set: %w", err)
		}
		log.Printf("Created toml provider in %s", path)
	case "package", "k":
		log.Printf("Creating package in %s", args)
		if !strings.HasSuffix(strings.ToLower(args), ".yaml") || !strings.HasSuffix(strings.ToLower(args), ".yml") {
			args = args + ".yaml"
		}
		path := filepath.Join(".", PackageConfigDir)
		if err := os.MkdirAll(path, 0777); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		path = filepath.Join(path, args)
		pk := yaml.New(path)
		if err := pk.Set(context.Context, info.With(map[string]any{
			"APP_NAME":      "qwerty",
			"LISTEN_PORT":   3030,
			"LOG_LEVEL":     "info",
			"VALKEY_SERVER": "$SHARED_VALKEY_CLUSTER",
		})); err != nil {
			return fmt.Errorf("failed to set: %w", err)
		}
		log.Printf("Created package in %s", args)
	default:
		return fmt.Errorf("unknown flag: %s", firstFlag)
	}

	return nil
}
