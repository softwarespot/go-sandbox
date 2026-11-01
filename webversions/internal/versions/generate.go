package versions

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

type generateExtractionInput struct {
	Content       string
	Value         string
	SearchFromEnd bool
}

func generateExtraction(input generateExtractionInput) (ExtractionResult, error) {
	if input.Content == "" {
		return ExtractionResult{}, errors.New("content is empty")
	}
	if input.Value == "" {
		return ExtractionResult{}, errors.New("value is empty")
	}

	startIdxs := getValueStartIndexes(input.Content, input.Value, input.SearchFromEnd)
	if len(startIdxs) == 0 {
		return ExtractionResult{}, fmt.Errorf("value %q not found", input.Value)
	}

	for _, startIdx := range startIdxs {
		extracted, ok := getExtractedResults(input, startIdx)
		if !ok {
			continue
		}
		return extracted, nil
	}
	return ExtractionResult{}, errors.New("unable to generate a unique prefix/suffix for the value")
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

const maxPrefixLookback = 128

func getExtractedResults(input generateExtractionInput, startIdx int) (ExtractionResult, bool) {
	endIdx := startIdx + len(input.Value)
	suffix := getSuffix(input.Content, endIdx)
	content := []rune(input.Content[:startIdx])
	maxLookback := min(len(content), maxPrefixLookback)
	for length := 0; length <= maxLookback; length++ {
		r := content[len(content)-length:]
		prefix := string(r)
		if res, ok := extractForPositionAt(input, prefix, suffix, startIdx); ok {
			return res, true
		}
	}
	return ExtractionResult{}, false
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

func extractForPositionAt(input generateExtractionInput, prefix, suffix string, startIdx int) (ExtractionResult, bool) {
	extractInput := extractionInput{
		Content:       input.Content,
		Prefixes:      []string{prefix},
		Suffix:        suffix,
		SearchFromEnd: input.SearchFromEnd,
	}
	res, err := extract(extractInput)
	if err != nil ||
		res.Match.Value != input.Value ||
		res.Match.StartIndex != startIdx {
		return ExtractionResult{}, false
	}
	return res, true
}
