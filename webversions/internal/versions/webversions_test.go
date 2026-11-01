package versions

import (
	"bytes"
	"strings"
	"testing"

	testhelpers "webversions/test-helpers"
)

func Test_parseWebVersionsLines(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		want    []AppConfig
		wantErr bool
	}{
		{
			name: "valid file",
			src: strings.Join(
				[]string{
					webVersionsHeader,
					"App1;0.0.0;0.1.0;https://www.app1.com;false;<;Ver ;;;;0.2.0;INFORMATION;TAG1",
					"App2;1.0.0;1.0.0;https://www.app2.com;false; ;Ver ;;;;0.2.0;INFORMATION;TAG1,TAG2",
				},
				"\n",
			),
			want: []AppConfig{
				{
					ID:          "1",
					Name:        "App1",
					URL:         "https://www.app1.com",
					Info:        "INFORMATION",
					CurrVersion: "0.0.0",
					WebVersion:  "0.1.0",
					Prefixes: []string{
						"Ver ",
						"",
						"",
						"",
					},
					Suffix:        "<",
					SearchFromEnd: false,
					TabNames:      []string{"TAG1"},
				},
				{
					ID:          "2",
					Name:        "App2",
					URL:         "https://www.app2.com",
					Info:        "INFORMATION",
					CurrVersion: "1.0.0",
					WebVersion:  "1.0.0",
					Prefixes: []string{
						"Ver ",
						"",
						"",
						"",
					},
					Suffix:        " ",
					SearchFromEnd: false,
					TabNames:      []string{"TAG1", "TAG2"},
				},
			},
		}, {
			name: "blank lines",
			src: strings.Join(
				[]string{
					webVersionsHeader,
					"",
					"App1;0.0.0;0.1.0;https://www.app1.com;false;<;Ver ;;;;0.2.0;INFORMATION;TAG1",
					"",
					"App2;1.0.0;1.0.0;https://www.app2.com;false; ;Ver ;;;;0.2.0;INFORMATION;TAG1,TAG2",
					"",
				},
				"\n",
			),
			want: []AppConfig{
				{
					ID:          "1",
					Name:        "App1",
					URL:         "https://www.app1.com",
					Info:        "INFORMATION",
					CurrVersion: "0.0.0",
					WebVersion:  "0.1.0",
					Prefixes: []string{
						"Ver ",
						"",
						"",
						"",
					},
					Suffix:        "<",
					SearchFromEnd: false,
					TabNames:      []string{"TAG1"},
				},
				{
					ID:          "2",
					Name:        "App2",
					URL:         "https://www.app2.com",
					Info:        "INFORMATION",
					CurrVersion: "1.0.0",
					WebVersion:  "1.0.0",
					Prefixes: []string{
						"Ver ",
						"",
						"",
						"",
					},
					Suffix:        " ",
					SearchFromEnd: false,
					TabNames:      []string{"TAG1", "TAG2"},
				},
			},
		},
		{
			name: "bad header",
			src: strings.Join(
				[]string{
					"Bad header;;",
					"App1;0.0.0;0.1.0;https://www.app1.com;false;<;Ver ;;;;0.2.0;INFORMATION;TAG1",
				},
				"\n",
			),
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfgs, err := parseWebVersionsLines(strings.NewReader(tt.src))
			if tt.wantErr {
				testhelpers.AssertError(t, err)
				return
			}
			testhelpers.AssertNoError(t, err)
			testhelpers.AssertEqual(t, cfgs, tt.want)
		})
	}
}

func Test_writeWebVersionsLines(t *testing.T) {
	cfgs := []AppConfig{
		{
			Name:          "app1App",
			URL:           "https://www.app1.com/",
			Info:          "INFORMATION",
			CurrVersion:   "1.0.0",
			WebVersion:    "1.1.0",
			Prefixes:      []string{"prefix:"},
			Suffix:        "</span>",
			SearchFromEnd: true,
			TabNames:      []string{"TAG1", "TAG2"},
		},
	}

	var buf bytes.Buffer
	testhelpers.AssertNoError(t, writeWebVersionsLines(&buf, cfgs))

	got, err := parseWebVersionsLines(&buf)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, len(got), 1)
	testhelpers.AssertEqual(t, got[0].Name, cfgs[0].Name)
	testhelpers.AssertEqual(t, got[0].URL, cfgs[0].URL)
	testhelpers.AssertEqual(t, got[0].Info, cfgs[0].Info)
	testhelpers.AssertEqual(t, got[0].CurrVersion, cfgs[0].CurrVersion)
	testhelpers.AssertEqual(t, got[0].WebVersion, cfgs[0].WebVersion)
	testhelpers.AssertEqual(t, got[0].Prefixes, []string{"prefix:", "", "", ""})
	testhelpers.AssertEqual(t, got[0].Suffix, cfgs[0].Suffix)
	testhelpers.AssertEqual(t, got[0].SearchFromEnd, cfgs[0].SearchFromEnd)
	testhelpers.AssertEqual(t, got[0].TabNames, cfgs[0].TabNames)
}
