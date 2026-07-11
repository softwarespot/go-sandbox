package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	configPath := flag.String("config", "WebVersions.txt", "path to WebVersions configuration file")
	userAgent := flag.String("user-agent", "WebVersions/0.0", "User-Agent header for downloads")
	timeout := flag.Duration("timeout", 10*time.Second, "HTTP request timeout")
	annotate := flag.Bool("annotate", false, "annotate extracted matches in downloaded content")
	flag.Parse()

	b, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(bytes.NewReader(b))
	db, err := NewWebVersionsDB(r)
	if err != nil {
		panic(err)
	}

	downloader := NewDownloader(
		WithTimeout(*timeout),
		WithUserAgent(*userAgent),
	)
	for _, cfg := range db.Configs() {
		content, err := downloader.Download(cfg.URL)
		if err != nil {
			fmt.Printf("Error downloading content for URL %s: %v\n", cfg.URL, err)
			continue
		}
		input := ExtractInput{
			Content: content,
			Prefixes: []string{
				cfg.Prefix1,
				cfg.Prefix2,
				cfg.Prefix3,
				cfg.Prefix4,
			},
			Suffix:        cfg.Suffix,
			SearchFromEnd: cfg.SearchFromEnd,
		}
		res, err := Extract(input)
		if err != nil {
			fmt.Printf("Error extracting value for URL %s: %v\n", cfg.URL, err)
			continue
		}

		if res.Match.Value == cfg.CurrVersion {
			fmt.Printf("%s is up to date: %s\n", cfg.AppName, res.Match.Value)
		} else {
			fmt.Printf("%s is outdated: current=%q found=%q\n", cfg.AppName, cfg.CurrVersion, res.Match.Value)
		}

		if *annotate {
			annotateMatches(content, res)
		}

		generateInput := GenerateExtractedInput{
			Content:       content,
			Value:         res.Match.Value,
			SearchFromEnd: cfg.SearchFromEnd,
		}
		generateExtractedRes, err := GenerateExtracted(generateInput)
		if err != nil {
			fmt.Printf("Error generating extracted value for URL %s: %v\n", cfg.URL, err)
			continue
		}

		var prefixes []string
		for _, prefix := range generateExtractedRes.Prefixes {
			prefixes = append(prefixes, fmt.Sprintf("%q", prefix.Value))
		}
		fmt.Printf("Generated prefixes: %s, suffix: %q for version %q\n", prefixes, generateExtractedRes.Suffix.Value, generateExtractedRes.Match.Value)
		fmt.Println("--------------------------------------------------")
	}
}

func annotateMatches(content string, res Extracted) {
	var insertAtOps []InsertManyAtOp
	for _, prefix := range res.Prefixes {
		if prefix.StartIndex >= 0 {
			insertAtOps = append(insertAtOps, InsertManyAtOp{
				Value:  "🐕",
				PosIdx: prefix.StartIndex,
			})
		}
		if prefix.EndIndex >= 0 {
			insertAtOps = append(insertAtOps, InsertManyAtOp{
				Value:  "🐕",
				PosIdx: prefix.EndIndex,
			})
		}
	}
	if res.Match.StartIndex >= 0 {
		insertAtOps = append(insertAtOps, InsertManyAtOp{
			Value:  "✮",
			PosIdx: res.Match.StartIndex,
		})
	}
	if res.Match.EndIndex >= 0 {
		insertAtOps = append(insertAtOps, InsertManyAtOp{
			Value:  "✮",
			PosIdx: res.Match.EndIndex,
		})
	}
	if res.Suffix.StartIndex >= 0 {
		insertAtOps = append(insertAtOps, InsertManyAtOp{
			Value:  "🐈",
			PosIdx: res.Suffix.StartIndex,
		})
	}
	if res.Suffix.EndIndex >= 0 {
		insertAtOps = append(insertAtOps, InsertManyAtOp{
			Value:  "🐈",
			PosIdx: res.Suffix.EndIndex,
		})
	}
	content = InsertManyAt(content, insertAtOps)
	fmt.Println(content)
}
