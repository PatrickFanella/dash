package importer

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseFile(path string) (*DashyConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

func Parse(r io.Reader) (*DashyConfig, error) {
	var cfg DashyConfig
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode YAML: %w", err)
	}
	if len(cfg.Sections) == 0 {
		return nil, fmt.Errorf("config contains no sections")
	}
	return &cfg, nil
}
