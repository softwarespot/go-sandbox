package versions

import (
	"testing"

	testhelpers "webversions/test-helpers"
)

func Test_generateExtraction(t *testing.T) {
	tests := []struct {
		name    string
		input   generateExtractionInput
		want    ExtractionResult
		wantErr bool
	}{
		{
			name: "resiliency to prepended letters and trailing punctuation noise",
			input: generateExtractionInput{
				Content: "New release!\nUpdated version: v1.2.3, please check the test string for more content.",
				Value:   "1.2.3",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "1.2.3", StartIndex: 31, EndIndex: 36},
				Prefixes: []ExtractionMatch{{Value: ": v", StartIndex: 28, EndIndex: 31}},
				Suffix:   ExtractionMatch{Value: ",", StartIndex: 36, EndIndex: 37},
			},
			wantErr: false,
		},
		{
			name: "HTML element attribute extraction",
			input: generateExtractionInput{
				Content: `<meta name="software_version" content="4.15.2-alpha">`,
				Value:   "4.15.2-alpha",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "4.15.2-alpha", StartIndex: 39, EndIndex: 51},
				Prefixes: []ExtractionMatch{{Value: `t="`, StartIndex: 36, EndIndex: 39}},
				Suffix:   ExtractionMatch{Value: `"`, StartIndex: 51, EndIndex: 52},
			},
			wantErr: false,
		},
		{
			name: "JSON configuration map structure",
			input: generateExtractionInput{
				Content: `{"api":{"status":"active","build":"v0.9.84"}}`,
				Value:   "0.9.84",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "0.9.84", StartIndex: 36, EndIndex: 42},
				Prefixes: []ExtractionMatch{{Value: `"v`, StartIndex: 34, EndIndex: 36}},
				Suffix:   ExtractionMatch{Value: `"`, StartIndex: 42, EndIndex: 43},
			},
			wantErr: false,
		},
		{
			name: "target version is located at the absolute beginning of text",
			input: generateExtractionInput{
				Content: `1.0.0-rc1 is ready for deployment`,
				Value:   "1.0.0-rc1",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "1.0.0-rc1", StartIndex: 0, EndIndex: 9},
				Prefixes: []ExtractionMatch{{Value: "", StartIndex: 0, EndIndex: 0}},
				Suffix:   ExtractionMatch{Value: " ", StartIndex: 9, EndIndex: 10},
			},
			wantErr: false,
		},
		{
			name: "target version is located at the absolute end of text",
			input: generateExtractionInput{
				Content: `Current system version is v3.4.1`,
				Value:   "3.4.1",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "3.4.1", StartIndex: 27, EndIndex: 32},
				Prefixes: []ExtractionMatch{{Value: "s v", StartIndex: 24, EndIndex: 27}},
				Suffix:   ExtractionMatch{Value: "", StartIndex: -1, EndIndex: -1},
			},
			wantErr: false,
		},
		{
			name: "isolates identical duplicate strings by growing lookback window",
			input: generateExtractionInput{
				Content: `service-a version: v1.2.0, service-b version: v1.2.0`,
				Value:   "1.2.0",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "1.2.0", StartIndex: 20, EndIndex: 25},
				Prefixes: []ExtractionMatch{{Value: ": v", StartIndex: 17, EndIndex: 20}},
				Suffix:   ExtractionMatch{Value: ",", StartIndex: 25, EndIndex: 26},
			},
			wantErr: false,
		},
		{
			name: "duplicate version with SearchFromEnd selects last occurrence",
			input: generateExtractionInput{
				Content:       `service-a version: v1.2.0, service-b version: v1.2.0`,
				Value:         "1.2.0",
				SearchFromEnd: true,
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "1.2.0", StartIndex: 47, EndIndex: 52},
				Prefixes: []ExtractionMatch{{Value: "v", StartIndex: 46, EndIndex: 47}},
				Suffix:   ExtractionMatch{Value: "", StartIndex: -1, EndIndex: -1},
			},
			wantErr: false,
		},
		{
			name: "multi-byte UTF-8 character boundary safety",
			input: generateExtractionInput{
				Content: `🚀 Release Version: ✨2.11.0✨ available now`,
				Value:   "2.11.0",
			},
			want: ExtractionResult{
				Match:    ExtractionMatch{Value: "2.11.0", StartIndex: 25, EndIndex: 31},
				Prefixes: []ExtractionMatch{{Value: "✨", StartIndex: 22, EndIndex: 25}},
				Suffix:   ExtractionMatch{Value: "✨", StartIndex: 31, EndIndex: 34},
			},
			wantErr: false,
		},
		{
			name: "error path: target version missing entirely from payload",
			input: generateExtractionInput{
				Content: `Log line without any semantic numbers here`,
				Value:   "1.0.0",
			},
			want:    ExtractionResult{},
			wantErr: true,
		},
		{
			name: "error path: empty version input validation failure",
			input: generateExtractionInput{
				Content: `Version: 1.0.0`,
				Value:   "",
			},
			want:    ExtractionResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateExtraction(tt.input)
			testhelpers.AssertEqual(t, err != nil, tt.wantErr)
			if !tt.wantErr {
				testhelpers.AssertEqual(t, got, tt.want)
			}
		})
	}
}

func Fuzz_generateExtraction(f *testing.F) {
	// Command: go test -fuzz=FuzzGenerateExtracted -fuzztime=30s
	f.Add("release: v1.2.3", "1.2.3", false)
	f.Add("release: v2.0.0\nrelease: v2.1.0", "2.1.0", true)
	f.Add("service-a version: v1.2.0, service-b version: v1.2.0", "1.2.0", false)

	f.Fuzz(func(t *testing.T, content, value string, searchFromEnd bool) {
		if value == "" {
			t.Skip()
		}
		generateInput := generateExtractionInput{
			Content:       content,
			Value:         value,
			SearchFromEnd: searchFromEnd,
		}
		res, err := generateExtraction(generateInput)
		if err != nil {
			return
		}
		var prefixes []string
		for _, prefix := range res.Prefixes {
			prefixes = append(prefixes, prefix.Value)
		}
		extractInput := extractionInput{
			Content:       content,
			Prefixes:      prefixes,
			Suffix:        res.Suffix.Value,
			SearchFromEnd: searchFromEnd,
		}
		got, err := extract(extractInput)
		testhelpers.AssertEqual(t, err == nil, true)
		testhelpers.AssertEqual(t, got.Match.Value, value)
	})
}
