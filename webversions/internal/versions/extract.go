package versions

import (
	"errors"
	"fmt"
	"strings"
)

type extractionInput struct {
	Content       string
	Prefixes      []string
	Suffix        string
	SearchFromEnd bool
}

// ExtractionResult represents the result of an extraction operation.
type ExtractionResult struct {
	Match    ExtractionMatch
	Prefixes []ExtractionMatch
	Suffix   ExtractionMatch
}

// ExtractionMatch represents a matched value along with its start and end indexes in the content.
type ExtractionMatch struct {
	Value      string
	StartIndex int
	EndIndex   int
}

func extract(input extractionInput) (ExtractionResult, error) {
	if input.Content == "" {
		return ExtractionResult{}, errors.New("content is empty")
	}
	if len(input.Prefixes) == 0 {
		return ExtractionResult{}, errors.New("no prefixes provided")
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
		return ExtractionResult{}, fmt.Errorf("prefix %q not found", prefix)
	}

	out := ExtractionResult{
		Prefixes: []ExtractionMatch{
			newExtractionMatch(prefix, startIdx),
		},
	}
	absEndIdx := startIdx + len(prefix)
	content := input.Content[absEndIdx:]

	for _, prefix := range input.Prefixes[1:] {
		if prefix == "" {
			out.Prefixes = append(out.Prefixes, newEmptyExtractionMatch(prefix))
			continue
		}
		idx := strings.Index(content, prefix)
		if idx == -1 {
			return ExtractionResult{}, fmt.Errorf("prefix %q not found", prefix)
		}

		out.Prefixes = append(out.Prefixes, newExtractionMatch(prefix, absEndIdx+idx))
		nextStartIdx := idx + len(prefix)
		absEndIdx += nextStartIdx
		content = content[nextStartIdx:]
	}

	var value string
	if input.Suffix == "" {
		value = content
		out.Suffix = newEmptyExtractionMatch(input.Suffix)
	} else {
		idx := strings.Index(content, input.Suffix)
		if idx == -1 {
			return ExtractionResult{}, fmt.Errorf("suffix %q not found", input.Suffix)
		}
		value = content[:idx]
		out.Suffix = newExtractionMatch(input.Suffix, absEndIdx+idx)
	}

	if value == "" {
		return ExtractionResult{}, errors.New("value is empty")
	}

	out.Match = newExtractionMatch(value, absEndIdx)

	return out, nil
}

func newExtractionMatch(value string, startIndex int) ExtractionMatch {
	return ExtractionMatch{
		Value:      value,
		StartIndex: startIndex,
		EndIndex:   startIndex + len(value),
	}
}

func newEmptyExtractionMatch(value string) ExtractionMatch {
	return ExtractionMatch{
		Value:      value,
		StartIndex: -1,
		EndIndex:   -1,
	}
}
