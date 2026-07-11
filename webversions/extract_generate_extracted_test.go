package main

import (
	"testing"

	testhelpers "webversions/test-helpers"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name    string
		input   ExtractInput
		want    Extracted
		wantErr bool
	}{
		{
			name: "standard HTML extraction",
			input: ExtractInput{
				Content: `<div class="downloads"><span class="v">1.42.38-beta</span></div>`,
				Prefixes: []string{
					`<div class="downloads">`,
					`<span class="v">`,
					"",
					"",
				},
				Suffix: `</span>`,
			},
			want: Extracted{
				Match: ExtractedMatch{Value: "1.42.38-beta", StartIndex: 39, EndIndex: 51},
				Prefixes: []ExtractedMatch{
					{Value: `<div class="downloads">`, StartIndex: 0, EndIndex: 23},
					{Value: `<span class="v">`, StartIndex: 23, EndIndex: 39},
					{Value: "", StartIndex: -1, EndIndex: -1}, // Prefix3 empty
					{Value: "", StartIndex: -1, EndIndex: -1}, // Prefix4 empty
				},
				Suffix: ExtractedMatch{Value: "</span>", StartIndex: 51, EndIndex: 58},
			},
			wantErr: false,
		},
		{
			name: "JSON configuration format",
			input: ExtractInput{
				Content: `{"status": "success", "data": {"current_version": "v2.0.1-rc1", "build": 99}}`,
				Prefixes: []string{
					`"current_version":`,
					`"v`,
					"",
					"",
				},
				Suffix: `"`,
			},
			want: Extracted{
				Match: ExtractedMatch{Value: "2.0.1-rc1", StartIndex: 52, EndIndex: 61},
				Prefixes: []ExtractedMatch{
					{Value: `"current_version":`, StartIndex: 31, EndIndex: 49},
					{Value: `"v`, StartIndex: 50, EndIndex: 52},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractedMatch{Value: `"`, StartIndex: 61, EndIndex: 62},
			},
			wantErr: false,
		},
		{
			name: "raw unstructured server log line",
			input: ExtractInput{
				Content: `INFO 2026-07-20 12:00:00 [release_pipeline] Deploying artifact version:12.4.0-patch2 to prod`,
				Prefixes: []string{
					`version:`,
					"",
					"",
					"",
				},
				Suffix: ` `,
			},
			want: Extracted{
				Match: ExtractedMatch{Value: "12.4.0-patch2", StartIndex: 71, EndIndex: 84},
				Prefixes: []ExtractedMatch{
					{Value: `version:`, StartIndex: 63, EndIndex: 71},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractedMatch{Value: ` `, StartIndex: 84, EndIndex: 85},
			},
			wantErr: false,
		},
		{
			name: "search from end bottom-up scan parity",
			input: ExtractInput{
				Content: "release: v1.0.0\nrelease: v1.1.0\nrelease: v1.2.5\n",
				Prefixes: []string{
					"release: v",
					"",
					"",
					"",
				},
				SearchFromEnd: true,
			},
			want: Extracted{
				Match: ExtractedMatch{Value: "1.2.5\n", StartIndex: 42, EndIndex: 48},
				Prefixes: []ExtractedMatch{
					{Value: "release: v", StartIndex: 32, EndIndex: 42},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractedMatch{Value: "", StartIndex: -1, EndIndex: -1},
			},
			wantErr: false,
		},
		{
			name: "empty initial prefix at content start",
			input: ExtractInput{
				Content: "1.0.0-rc1 is ready for deployment",
				Prefixes: []string{
					"",
					"",
					"",
					"",
				},
				Suffix: " ",
			},
			want: Extracted{
				Match: ExtractedMatch{Value: "1.0.0-rc1", StartIndex: 0, EndIndex: 9},
				Prefixes: []ExtractedMatch{
					{Value: "", StartIndex: 0, EndIndex: 0},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractedMatch{Value: " ", StartIndex: 9, EndIndex: 10},
			},
			wantErr: false,
		},
		{
			name: "missing structural prefix failure",
			input: ExtractInput{
				Content: `<div>Version: 1.0.0</div>`,
				Prefixes: []string{
					`<span class="missing">`,
					"",
					"",
					"",
				},
			},
			want:    Extracted{},
			wantErr: true,
		},
		{
			name: "missing trailing suffix failure",
			input: ExtractInput{
				Content: `<div>Version: 1.0.0</div>`,
				Prefixes: []string{
					`Version: `,
					"",
					"",
					"",
				},
				Suffix: `</span-mismatch>`,
			},
			want:    Extracted{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Extract(tt.input)
			testhelpers.AssertEqual(t, err != nil, tt.wantErr)
			testhelpers.AssertEqual(t, got, tt.want)
		})
	}
}

func TestGenerateExtracted(t *testing.T) {
	tests := []struct {
		name    string
		input   GenerateExtractedInput
		want    Extracted
		wantErr bool
	}{
		{
			name: "resiliency to prepended letters and trailing punctuation noise",
			input: GenerateExtractedInput{
				Content: "New release!\nUpdated version: v1.2.3, please check the test string for more content.",
				Value:   "1.2.3",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "1.2.3", StartIndex: 31, EndIndex: 36},
				Prefixes: []ExtractedMatch{{Value: ": v", StartIndex: 28, EndIndex: 31}},
				Suffix:   ExtractedMatch{Value: ",", StartIndex: 36, EndIndex: 37},
			},
			wantErr: false,
		},
		{
			name: "HTML element attribute extraction",
			input: GenerateExtractedInput{
				Content: `<meta name="software_version" content="4.15.2-alpha">`,
				Value:   "4.15.2-alpha",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "4.15.2-alpha", StartIndex: 39, EndIndex: 51},
				Prefixes: []ExtractedMatch{{Value: `t="`, StartIndex: 36, EndIndex: 39}},
				Suffix:   ExtractedMatch{Value: `"`, StartIndex: 51, EndIndex: 52},
			},
			wantErr: false,
		},
		{
			name: "JSON configuration map structure",
			input: GenerateExtractedInput{
				Content: `{"api":{"status":"active","build":"v0.9.84"}}`,
				Value:   "0.9.84",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "0.9.84", StartIndex: 36, EndIndex: 42},
				Prefixes: []ExtractedMatch{{Value: `"v`, StartIndex: 34, EndIndex: 36}},
				Suffix:   ExtractedMatch{Value: `"`, StartIndex: 42, EndIndex: 43},
			},
			wantErr: false,
		},
		{
			name: "target version is located at the absolute beginning of text",
			input: GenerateExtractedInput{
				Content: `1.0.0-rc1 is ready for deployment`,
				Value:   "1.0.0-rc1",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "1.0.0-rc1", StartIndex: 0, EndIndex: 9},
				Prefixes: []ExtractedMatch{{Value: "", StartIndex: 0, EndIndex: 0}},
				Suffix:   ExtractedMatch{Value: " ", StartIndex: 9, EndIndex: 10},
			},
			wantErr: false,
		},
		{
			name: "target version is located at the absolute end of text",
			input: GenerateExtractedInput{
				Content: `Current system version is v3.4.1`,
				Value:   "3.4.1",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "3.4.1", StartIndex: 27, EndIndex: 32},
				Prefixes: []ExtractedMatch{{Value: "s v", StartIndex: 24, EndIndex: 27}},
				Suffix:   ExtractedMatch{Value: "", StartIndex: -1, EndIndex: -1},
			},
			wantErr: false,
		},
		{
			name: "isolates identical duplicate strings by growing lookback window",
			input: GenerateExtractedInput{
				Content: `service-a version: v1.2.0, service-b version: v1.2.0`,
				Value:   "1.2.0",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "1.2.0", StartIndex: 20, EndIndex: 25},
				Prefixes: []ExtractedMatch{{Value: ": v", StartIndex: 17, EndIndex: 20}},
				Suffix:   ExtractedMatch{Value: ",", StartIndex: 25, EndIndex: 26},
			},
			wantErr: false,
		},
		{
			name: "duplicate version with SearchFromEnd selects last occurrence",
			input: GenerateExtractedInput{
				Content:       `service-a version: v1.2.0, service-b version: v1.2.0`,
				Value:         "1.2.0",
				SearchFromEnd: true,
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "1.2.0", StartIndex: 47, EndIndex: 52},
				Prefixes: []ExtractedMatch{{Value: "v", StartIndex: 46, EndIndex: 47}},
				Suffix:   ExtractedMatch{Value: "", StartIndex: -1, EndIndex: -1},
			},
			wantErr: false,
		},
		{
			name: "multi-byte UTF-8 character boundary safety",
			input: GenerateExtractedInput{
				Content: `🚀 Release Version: ✨2.11.0✨ available now`,
				Value:   "2.11.0",
			},
			want: Extracted{
				Match:    ExtractedMatch{Value: "2.11.0", StartIndex: 25, EndIndex: 31},
				Prefixes: []ExtractedMatch{{Value: "✨", StartIndex: 22, EndIndex: 25}},
				Suffix:   ExtractedMatch{Value: "✨", StartIndex: 31, EndIndex: 34},
			},
			wantErr: false,
		},
		{
			name: "error path: target version missing entirely from payload",
			input: GenerateExtractedInput{
				Content: `Log line without any semantic numbers here`,
				Value:   "1.0.0",
			},
			want:    Extracted{},
			wantErr: true,
		},
		{
			name: "error path: empty version input validation failure",
			input: GenerateExtractedInput{
				Content: `Version: 1.0.0`,
				Value:   "",
			},
			want:    Extracted{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateExtracted(tt.input)
			testhelpers.AssertEqual(t, err != nil, tt.wantErr)
			if !tt.wantErr {
				testhelpers.AssertEqual(t, got, tt.want)
			}
		})
	}
}

func FuzzExtract(f *testing.F) {
	// Command: go test -fuzz=FuzzExtract -fuzztime=30s
	f.Add("version:1.2.3-alpha", "version:", "", "", "", "-", false)
	f.Add("<div class=\"v\">2.0.0</div>", "<div class=\"v\">", "", "", "", "</div>", false)
	f.Add("v1 v2 v3", "v", "", "", "", " ", true)

	f.Fuzz(func(t *testing.T, content, p1, p2, p3, p4, suffix string, searchFromEnd bool) {
		// If Extract has an index tracking bug, it will panic here.
		// Go's fuzzer treats panics as a failure automatically.
		extractInput := ExtractInput{
			Content:       content,
			Prefixes:      []string{p1, p2, p3, p4},
			Suffix:        suffix,
			SearchFromEnd: searchFromEnd,
		}
		_, _ = Extract(extractInput)
	})
}

func FuzzGenerateExtracted(f *testing.F) {
	// Command: go test -fuzz=FuzzGenerateExtracted -fuzztime=30s
	f.Add("release: v1.2.3", "1.2.3", false)
	f.Add("release: v2.0.0\nrelease: v2.1.0", "2.1.0", true)
	f.Add("service-a version: v1.2.0, service-b version: v1.2.0", "1.2.0", false)

	f.Fuzz(func(t *testing.T, content, value string, searchFromEnd bool) {
		if value == "" {
			t.Skip()
		}
		generateInput := GenerateExtractedInput{
			Content:       content,
			Value:         value,
			SearchFromEnd: searchFromEnd,
		}
		res, err := GenerateExtracted(generateInput)
		if err != nil {
			return
		}
		var prefixes []string
		for _, prefix := range res.Prefixes {
			prefixes = append(prefixes, prefix.Value)
		}
		extractInput := ExtractInput{
			Content:       content,
			Prefixes:      prefixes,
			Suffix:        res.Suffix.Value,
			SearchFromEnd: searchFromEnd,
		}
		got, err := Extract(extractInput)
		testhelpers.AssertEqual(t, err == nil, true)
		testhelpers.AssertEqual(t, got.Match.Value, value)
	})
}
