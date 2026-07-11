package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

type GenerateExtractedInput struct {
	Content       string
	Value         string
	SearchFromEnd bool
}

// GenerateExtracted generates the prefixes and suffix for the provided value.
func GenerateExtracted(input GenerateExtractedInput) (Extracted, error) {
	if input.Content == "" {
		return Extracted{}, errors.New("content is empty")
	}
	if input.Value == "" {
		return Extracted{}, errors.New("value is empty")
	}

	startIdxs := getValueStartIndexes(input.Content, input.Value, input.SearchFromEnd)
	if len(startIdxs) == 0 {
		return Extracted{}, fmt.Errorf("value %q not found", input.Value)
	}

	for _, startIdx := range startIdxs {
		endIdx := startIdx + len(input.Value)
		suffix := getSuffix(input.Content, endIdx)
		extracted, ok := getExtracted(input, suffix, startIdx)
		if !ok {
			continue
		}
		return extracted, nil
	}
	return Extracted{}, errors.New("unable to generate a unique prefix/suffix for the value")
}

func getValueStartIndexes(content, value string, searchFromEnd bool) []int {
	var (
		startIdxs []int
		startIdx  = 0
	)
	for {
		idx := strings.Index(content[startIdx:], value)
		if idx == -1 {
			break
		}
		currStartIdx := startIdx + idx
		startIdxs = append(startIdxs, currStartIdx)
		startIdx = currStartIdx + len(value)
	}
	if searchFromEnd {
		slices.Reverse(startIdxs)
	}
	return startIdxs
}

func getSuffix(content string, startIdx int) string {
	if startIdx >= len(content) {
		return ""
	}
	r, size := utf8.DecodeRuneInString(content[startIdx:])
	if size == 0 {
		return ""
	}
	return string(r)
}

const maxPrefixLookback = 128

func getExtracted(input GenerateExtractedInput, suffix string, startIdx int) (Extracted, bool) {
	content := []rune(input.Content[:startIdx])
	maxLookback := min(len(content), maxPrefixLookback)
	for length := 0; length <= maxLookback; length++ {
		r := content[len(content)-length:]
		prefix := string(r)
		if res, ok := extractForPositionAt(input, prefix, suffix, startIdx); ok {
			return res, true
		}
	}
	return Extracted{}, false
}

func extractForPositionAt(input GenerateExtractedInput, prefix, suffix string, startIdx int) (Extracted, bool) {
	extractInput := ExtractInput{
		Content:       input.Content,
		Prefixes:      []string{prefix},
		Suffix:        suffix,
		SearchFromEnd: input.SearchFromEnd,
	}
	res, err := Extract(extractInput)
	if err != nil ||
		res.Match.Value != input.Value ||
		res.Match.StartIndex != startIdx {
		return Extracted{}, false
	}
	return res, true
}
