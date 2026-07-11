package versions

import (
	"fmt"
	"slices"
)

type ConfigType string

const (
	ConfigTypeWebVersions ConfigType = "webversions"
)

type Manager struct {
	path    string
	cfgType ConfigType
	cfgs    []AppConfig
}

type ManagerInput struct {
	Path       string
	ConfigType ConfigType
}

// NewManager creates a new [Manager] instance based on the provided [ManagerInput].
func NewManager(input ManagerInput) (*Manager, error) {
	m := &Manager{
		path:    input.Path,
		cfgType: input.ConfigType,
		cfgs:    nil,
	}
	if err := m.Load(); err != nil {
		return nil, fmt.Errorf("load application configurations: %w", err)
	}
	return m, nil
}

// Load loads the application configurations managed by the [Manager].
func (m *Manager) Load() error {
	switch m.cfgType {
	case ConfigTypeWebVersions:
		cfgs, err := loadWebVersions(m.path)
		if err != nil {
			return fmt.Errorf("load webversions config file: %w", err)
		}
		m.cfgs = cfgs
		return nil
	}
	return fmt.Errorf("unsupported config type: %q", m.cfgType)
}

// Configs returns a copy of the application configurations managed by the [Manager].
func (m *Manager) Configs() []AppConfig {
	return slices.Clone(m.cfgs)
}

// Extract extracts a value from the provided content using the specified [AppConfig].
func (m *Manager) Extract(cfg AppConfig, content string) (ExtractionResult, error) {
	input := extractionInput{
		Content:       content,
		Prefixes:      cfg.Prefixes,
		Suffix:        cfg.Suffix,
		SearchFromEnd: cfg.SearchFromEnd,
	}
	res, err := extract(input)
	if err != nil {
		return ExtractionResult{}, fmt.Errorf("extract version: %w", err)
	}
	return res, nil
}

// Generate generates the prefixes and suffix for the provided value using the specified [AppConfig].
func (m *Manager) Generate(cfg AppConfig, content, version string) (ExtractionResult, error) {
	input := generateExtractionInput{
		Content:       content,
		Value:         version,
		SearchFromEnd: cfg.SearchFromEnd,
	}
	res, err := generateExtraction(input)
	if err != nil {
		return ExtractionResult{}, fmt.Errorf("generate extraction version: %w", err)
	}
	return res, nil
}

// Store stores the application configurations managed by the [Manager].
func (m *Manager) Store() error {
	switch m.cfgType {
	case ConfigTypeWebVersions:
		if err := storeWebVersions(m.path, m.cfgs); err != nil {
			return fmt.Errorf("store webversions config file: %w", err)
		}
		return nil
	}
	return fmt.Errorf("unsupported config type: %q", m.cfgType)
}

// UpdateConfigs replaces the current manager configurations with the provided slice.
func (m *Manager) UpdateConfigs(cfgs []AppConfig) {
	m.cfgs = slices.Clone(cfgs)
}
