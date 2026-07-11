package versions

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

const (
	webVersionsHeader          = "Doc version 3.1;;"
	defaultWebVersionsFileMode = 0o644
)

var webversionsFieldMapping = struct {
	AppName       int
	CurrVersion   int
	WebVersion    int
	URL           int
	SearchFromEnd int
	Suffix        int
	Prefix1       int
	Prefix2       int
	Prefix3       int
	Prefix4       int
	KnownVersion  int
	Info          int
	TabNames      int
	FieldsCount   int
}{
	AppName:       0,
	CurrVersion:   1,
	WebVersion:    2,
	URL:           3,
	SearchFromEnd: 4,
	Suffix:        5,
	Prefix1:       6,
	Prefix2:       7,
	Prefix3:       8,
	Prefix4:       9,
	KnownVersion:  10, // Unused, as it's unclear as what this field is used for.
	Info:          11,
	TabNames:      12,
	FieldsCount:   13,
}

func loadWebVersions(path string) ([]AppConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read webversions config file: %w", err)
	}
	r := bytes.NewReader(b)
	cfgs, err := parseWebVersionsLines(r)
	if err != nil {
		return nil, fmt.Errorf("parse webversions config file: %w", err)
	}
	return cfgs, nil
}

func parseWebVersionsLines(r io.Reader) ([]AppConfig, error) {
	var (
		scanner = bufio.NewScanner(r)
		cfgs    []AppConfig
		lineNo  = 0
		id      = 1
	)
	for scanner.Scan() {
		lineNo += 1
		line := scanner.Text()
		if lineNo == 1 {
			if strings.TrimSpace(line) != webVersionsHeader {
				return nil, fmt.Errorf("invalid webversions header, expected %q", webVersionsHeader)
			}
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.SplitN(line, ";", webversionsFieldMapping.FieldsCount)
		if len(fields) != webversionsFieldMapping.FieldsCount {
			return nil, fmt.Errorf("parse webversions config line %d: malformed line, expected %d fields, got %d", lineNo, webversionsFieldMapping.FieldsCount, len(fields))
		}
		cfg := AppConfig{
			ID:          strconv.Itoa(id),
			Name:        fields[webversionsFieldMapping.AppName],
			URL:         fields[webversionsFieldMapping.URL],
			Info:        fields[webversionsFieldMapping.Info],
			CurrVersion: fields[webversionsFieldMapping.CurrVersion],
			WebVersion:  fields[webversionsFieldMapping.WebVersion],
			Prefixes: []string{
				unquote(fields[webversionsFieldMapping.Prefix1]),
				unquote(fields[webversionsFieldMapping.Prefix2]),
				unquote(fields[webversionsFieldMapping.Prefix3]),
				unquote(fields[webversionsFieldMapping.Prefix4]),
			},
			Suffix:        unquote(fields[webversionsFieldMapping.Suffix]),
			SearchFromEnd: strings.EqualFold(fields[webversionsFieldMapping.SearchFromEnd], "true"),
			TabNames:      strings.Split(fields[webversionsFieldMapping.TabNames], ","),
		}
		cfgs = append(cfgs, cfg)
		id += 1
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parse webversions config: %w", err)
	}
	return cfgs, nil
}

func storeWebVersions(path string, cfgs []AppConfig) error {
	var w strings.Builder
	if err := writeWebVersionsLines(&w, cfgs); err != nil {
		return err
	}

	var mode fs.FileMode
	fi, err := os.Stat(path)
	switch {
	case err == nil:
		mode = fi.Mode().Perm()
	case errors.Is(err, os.ErrNotExist):
		mode = defaultWebVersionsFileMode
	default:
		return fmt.Errorf("stat webversions config file: %w", err)
	}
	if err := os.WriteFile(path, []byte(w.String()), mode); err != nil {
		return fmt.Errorf("write webversions config file: %w", err)
	}
	return nil
}

func writeWebVersionsLines(w io.Writer, cfgs []AppConfig) error {
	if _, err := fmt.Fprintln(w, webVersionsHeader); err != nil {
		return fmt.Errorf("write webversions config header: %w", err)
	}

	buf := make([]string, webversionsFieldMapping.FieldsCount)
	for _, cfg := range cfgs {
		buf[webversionsFieldMapping.AppName] = cfg.Name
		buf[webversionsFieldMapping.CurrVersion] = cfg.CurrVersion
		buf[webversionsFieldMapping.WebVersion] = cfg.WebVersion
		buf[webversionsFieldMapping.URL] = cfg.URL
		buf[webversionsFieldMapping.SearchFromEnd] = strconv.FormatBool(cfg.SearchFromEnd)
		buf[webversionsFieldMapping.Suffix] = quote(cfg.Suffix)
		prefixes := make([]string, 4)
		copy(prefixes, cfg.Prefixes)
		buf[webversionsFieldMapping.Prefix1] = quote(prefixes[0])
		buf[webversionsFieldMapping.Prefix2] = quote(prefixes[1])
		buf[webversionsFieldMapping.Prefix3] = quote(prefixes[2])
		buf[webversionsFieldMapping.Prefix4] = quote(prefixes[3])
		buf[webversionsFieldMapping.KnownVersion] = ""
		buf[webversionsFieldMapping.Info] = cfg.Info
		buf[webversionsFieldMapping.TabNames] = strings.Join(cfg.TabNames, ",")
		if _, err := fmt.Fprintln(w, strings.Join(buf, ";")); err != nil {
			return fmt.Errorf("write webversions config line: %w", err)
		}
	}
	return nil
}
