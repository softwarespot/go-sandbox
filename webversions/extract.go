package main

import (
	"errors"
	"fmt"
	"strings"
)

type ExtractInput struct {
	Content       string
	Prefixes      []string
	Suffix        string
	SearchFromEnd bool
}

type Extracted struct {
	Match    ExtractedMatch
	Prefixes []ExtractedMatch
	Suffix   ExtractedMatch
}

type ExtractedMatch struct {
	Value      string
	StartIndex int
	EndIndex   int
}

// Extract extracts the value using the provided prefixes and suffix.
func Extract(input ExtractInput) (Extracted, error) {
	if input.Content == "" {
		return Extracted{}, errors.New("content is empty")
	}
	if len(input.Prefixes) == 0 {
		return Extracted{}, errors.New("no prefixes provided")
	}

	var (
		prefix   = input.Prefixes[0]
		startIdx int
	)
	switch {
	case prefix == "":
		startIdx = 0
	case input.SearchFromEnd:
		startIdx = strings.LastIndex(input.Content, prefix)
	default:
		startIdx = strings.Index(input.Content, prefix)
	}
	if startIdx == -1 {
		return Extracted{}, fmt.Errorf("prefix %q not found", prefix)
	}

	out := Extracted{
		Prefixes: []ExtractedMatch{
			newExtractedMatch(prefix, startIdx),
		},
	}
	absEndIdx := startIdx + len(prefix)
	content := input.Content[absEndIdx:]

	for _, prefix := range input.Prefixes[1:] {
		if prefix == "" {
			out.Prefixes = append(out.Prefixes, newEmptyExtractedMatch(prefix))
			continue
		}
		idx := strings.Index(content, prefix)
		if idx == -1 {
			return Extracted{}, fmt.Errorf("prefix %q not found", prefix)
		}

		out.Prefixes = append(out.Prefixes, newExtractedMatch(prefix, absEndIdx+idx))
		nextStartIdx := idx + len(prefix)
		absEndIdx += nextStartIdx
		content = content[nextStartIdx:]
	}

	var value string
	if input.Suffix == "" {
		value = content
		out.Suffix = newEmptyExtractedMatch(input.Suffix)
	} else {
		idx := strings.Index(content, input.Suffix)
		if idx == -1 {
			return Extracted{}, fmt.Errorf("suffix %q not found", input.Suffix)
		}
		value = content[:idx]
		out.Suffix = newExtractedMatch(input.Suffix, absEndIdx+idx)
	}

	if value == "" {
		return Extracted{}, errors.New("value is empty")
	}

	out.Match = newExtractedMatch(value, absEndIdx)

	return out, nil
}

func newExtractedMatch(value string, startIndex int) ExtractedMatch {
	return ExtractedMatch{
		Value:      value,
		StartIndex: startIndex,
		EndIndex:   startIndex + len(value),
	}
}

func newEmptyExtractedMatch(value string) ExtractedMatch {
	return ExtractedMatch{
		Value:      value,
		StartIndex: -1,
		EndIndex:   -1,
	}
}
