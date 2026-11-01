package versions

import (
	"testing"

	testhelpers "webversions/test-helpers"
)

func Test_extract(t *testing.T) {
	tests := []struct {
		name    string
		input   extractionInput
		want    ExtractionResult
		wantErr bool
	}{
		{
			name: "standard HTML extraction",
			input: extractionInput{
				Content: `<div class="downloads"><span class="v">1.42.38-beta</span></div>`,
				Prefixes: []string{
					`<div class="downloads">`,
					`<span class="v">`,
					"",
					"",
				},
				Suffix: `</span>`,
			},
			want: ExtractionResult{
				Match: ExtractionMatch{Value: "1.42.38-beta", StartIndex: 39, EndIndex: 51},
				Prefixes: []ExtractionMatch{
					{Value: `<div class="downloads">`, StartIndex: 0, EndIndex: 23},
					{Value: `<span class="v">`, StartIndex: 23, EndIndex: 39},
					{Value: "", StartIndex: -1, EndIndex: -1}, // Prefix3 empty
					{Value: "", StartIndex: -1, EndIndex: -1}, // Prefix4 empty
				},
				Suffix: ExtractionMatch{Value: "</span>", StartIndex: 51, EndIndex: 58},
			},
			wantErr: false,
		},
		{
			name: "JSON configuration format",
			input: extractionInput{
				Content: `{"status": "success", "data": {"current_version": "v2.0.1-rc1", "build": 99}}`,
				Prefixes: []string{
					`"current_version":`,
					`"v`,
					"",
					"",
				},
				Suffix: `"`,
			},
			want: ExtractionResult{
				Match: ExtractionMatch{Value: "2.0.1-rc1", StartIndex: 52, EndIndex: 61},
				Prefixes: []ExtractionMatch{
					{Value: `"current_version":`, StartIndex: 31, EndIndex: 49},
					{Value: `"v`, StartIndex: 50, EndIndex: 52},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractionMatch{Value: `"`, StartIndex: 61, EndIndex: 62},
			},
			wantErr: false,
		},
		{
			name: "raw unstructured server log line",
			input: extractionInput{
				Content: `INFO 2026-07-20 12:00:00 [release_pipeline] Deploying artifact version:12.4.0-patch2 to prod`,
				Prefixes: []string{
					`version:`,
					"",
					"",
					"",
				},
				Suffix: ` `,
			},
			want: ExtractionResult{
				Match: ExtractionMatch{Value: "12.4.0-patch2", StartIndex: 71, EndIndex: 84},
				Prefixes: []ExtractionMatch{
					{Value: `version:`, StartIndex: 63, EndIndex: 71},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractionMatch{Value: ` `, StartIndex: 84, EndIndex: 85},
			},
			wantErr: false,
		},
		{
			name: "search from end bottom-up scan parity",
			input: extractionInput{
				Content: "release: v1.0.0\nrelease: v1.1.0\nrelease: v1.2.5\n",
				Prefixes: []string{
					"release: v",
					"",
					"",
					"",
				},
				SearchFromEnd: true,
			},
			want: ExtractionResult{
				Match: ExtractionMatch{Value: "1.2.5\n", StartIndex: 42, EndIndex: 48},
				Prefixes: []ExtractionMatch{
					{Value: "release: v", StartIndex: 32, EndIndex: 42},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractionMatch{Value: "", StartIndex: -1, EndIndex: -1},
			},
			wantErr: false,
		},
		{
			name: "empty initial prefix at content start",
			input: extractionInput{
				Content: "1.0.0-rc1 is ready for deployment",
				Prefixes: []string{
					"",
					"",
					"",
					"",
				},
				Suffix: " ",
			},
			want: ExtractionResult{
				Match: ExtractionMatch{Value: "1.0.0-rc1", StartIndex: 0, EndIndex: 9},
				Prefixes: []ExtractionMatch{
					{Value: "", StartIndex: 0, EndIndex: 0},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
					{Value: "", StartIndex: -1, EndIndex: -1},
				},
				Suffix: ExtractionMatch{Value: " ", StartIndex: 9, EndIndex: 10},
			},
			wantErr: false,
		},
		{
			name: "missing structural prefix failure",
			input: extractionInput{
				Content: `<div>Version: 1.0.0</div>`,
				Prefixes: []string{
					`<span class="missing">`,
					"",
					"",
					"",
				},
			},
			want:    ExtractionResult{},
			wantErr: true,
		},
		{
			name: "missing trailing suffix failure",
			input: extractionInput{
				Content: `<div>Version: 1.0.0</div>`,
				Prefixes: []string{
					`Version: `,
					"",
					"",
					"",
				},
				Suffix: `</span-mismatch>`,
			},
			want:    ExtractionResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extract(tt.input)
			testhelpers.AssertEqual(t, err != nil, tt.wantErr)
			testhelpers.AssertEqual(t, got, tt.want)
		})
	}
}

func Fuzz_extract(f *testing.F) {
	// Command: go test -fuzz=Fuzz_extract -fuzztime=30s
	f.Add("version:1.2.3-alpha", "version:", "", "", "", "-", false)
	f.Add("<div class=\"v\">2.0.0</div>", "<div class=\"v\">", "", "", "", "</div>", false)
	f.Add("v1 v2 v3", "v", "", "", "", " ", true)

	f.Fuzz(func(t *testing.T, content, p1, p2, p3, p4, suffix string, searchFromEnd bool) {
		// If extract has an index tracking bug, it will panic here.
		// Go's fuzzer treats panics as a failure automatically.
		extractInput := extractionInput{
			Content:       content,
			Prefixes:      []string{p1, p2, p3, p4},
			Suffix:        suffix,
			SearchFromEnd: searchFromEnd,
		}
		_, _ = extract(extractInput)
	})
}
