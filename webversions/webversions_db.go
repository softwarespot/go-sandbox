package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type WebVersionsDB struct {
	cfgs []WebVersionsConfig
}

type WebVersionsConfig struct {
	AppName       string
	URL           string
	CurrVersion   string
	WebVersion    string
	KnownVersion  string // Unclear as to why this is needed?
	SearchFromEnd bool
	Prefix1       string
	Prefix2       string
	Prefix3       string
	Prefix4       string
	Suffix        string
	Info          string
	Tabs          []string
}

func NewWebVersionsDB(r io.Reader) (*WebVersionsDB, error) {
	var (
		scanner        = bufio.NewScanner(r)
		cfgs           []WebVersionsConfig
		consumedHeader bool
	)
	for scanner.Scan() {
		// Check if the initial header has been consumed
		// i.e. this should be the first line.
		if !consumedHeader {
			consumedHeader = true
			continue
		}
		cfg, err := parseToAppConfig(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("parse to application config: %w", err)
		}
		cfgs = append(cfgs, cfg)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner.Err: %w", err)
	}
	db := &WebVersionsDB{
		cfgs: cfgs,
	}
	return db, nil
}

func (db *WebVersionsDB) Configs() []WebVersionsConfig {
	return slices.Clone(db.cfgs)
}

type fieldToIndexMapping struct {
	AppName       int
	URL           int
	CurrVersion   int
	WebVersion    int
	KnownVersion  int
	SearchFromEnd int
	Prefix1       int
	Prefix2       int
	Prefix3       int
	Prefix4       int
	Suffix        int
	Info          int
	Tabs          int
}

var fieldIndex = fieldToIndexMapping{
	AppName:       0,
	URL:           3,
	CurrVersion:   1,
	WebVersion:    2,
	KnownVersion:  10,
	SearchFromEnd: 4,
	Prefix1:       6,
	Prefix2:       7,
	Prefix3:       8,
	Prefix4:       9,
	Suffix:        5,
	Info:          11,
	Tabs:          12,
}

const requiredFieldsCount = 13

func parseToAppConfig(line string) (WebVersionsConfig, error) {
	fields := strings.SplitN(line, ";", requiredFieldsCount)
	if len(fields) != requiredFieldsCount {
		return WebVersionsConfig{}, fmt.Errorf("malformed line, expected %d fields, got %d", requiredFieldsCount, len(fields))
	}
	cfg := WebVersionsConfig{
		AppName:       fields[fieldIndex.AppName],
		URL:           fields[fieldIndex.URL],
		CurrVersion:   fields[fieldIndex.CurrVersion],
		WebVersion:    fields[fieldIndex.WebVersion],
		KnownVersion:  fields[fieldIndex.KnownVersion],
		SearchFromEnd: strings.EqualFold(fields[fieldIndex.SearchFromEnd], "true"),
		Prefix1:       unquote(fields[fieldIndex.Prefix1]),
		Prefix2:       unquote(fields[fieldIndex.Prefix2]),
		Prefix3:       unquote(fields[fieldIndex.Prefix3]),
		Prefix4:       unquote(fields[fieldIndex.Prefix4]),
		Suffix:        unquote(fields[fieldIndex.Suffix]),
		Info:          fields[fieldIndex.Info],
		Tabs:          strings.Split(fields[fieldIndex.Tabs], ","),
	}
	return cfg, nil
}

func unquote(s string) string {
	us, err := strconv.Unquote(s)
	if err != nil {
		return s
	}
	return us
}
