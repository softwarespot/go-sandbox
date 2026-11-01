package main

import (
	"fmt"

	"webversions/internal/downloader"
	"webversions/internal/helpers"
	"webversions/internal/versions"
)

func runCLI(opts cliFlags) error {
	dlr := downloader.New(
		downloader.WithTimeout(opts.Timeout),
		downloader.WithUserAgent(opts.UserAgent),
	)

	managerInput := versions.ManagerInput{
		Path:       opts.ConfigPath,
		ConfigType: versions.ConfigTypeWebVersions,
	}
	m, err := versions.NewManager(managerInput)
	if err != nil {
		return fmt.Errorf("new manager: %w", err)
	}

	for _, cfg := range m.Configs() {
		content, err := dlr.Download(cfg.URL)
		if err != nil {
			fmt.Printf("Error downloading content for URL %s: %v\n", cfg.URL, err)
			continue
		}

		res, err := m.Extract(cfg, content)
		if err != nil {
			fmt.Printf("Error extracting version for URL %s: %v\n", cfg.URL, err)
			continue
		}

		if res.Match.Value == cfg.CurrVersion {
			fmt.Printf("%s is up to date: %s\n", cfg.Name, res.Match.Value)
		} else {
			fmt.Printf("%s is outdated: current=%q found=%q\n", cfg.Name, cfg.CurrVersion, res.Match.Value)
		}
		if opts.Extra {
			showAnnotatedContent(content, res)
			showGenerateExtracted(m, cfg, content, res.Match.Value)
		}
		fmt.Println("--------------------------------------------------")
	}
	return nil
}

func showAnnotatedContent(content string, extractedRes versions.ExtractionResult) {
	var ops []helpers.InsertAtOp
	for _, prefix := range extractedRes.Prefixes {
		if prefix.StartIndex >= 0 {
			ops = append(ops, helpers.InsertAtOp{
				Value:  "🐕",
				PosIdx: prefix.StartIndex,
			})
		}
		if prefix.EndIndex >= 0 {
			ops = append(ops, helpers.InsertAtOp{
				Value:  "🐕",
				PosIdx: prefix.EndIndex,
			})
		}
	}
	if extractedRes.Match.StartIndex >= 0 {
		ops = append(ops, helpers.InsertAtOp{
			Value:  "✮",
			PosIdx: extractedRes.Match.StartIndex,
		})
	}
	if extractedRes.Match.EndIndex >= 0 {
		ops = append(ops, helpers.InsertAtOp{
			Value:  "✮",
			PosIdx: extractedRes.Match.EndIndex,
		})
	}
	if extractedRes.Suffix.StartIndex >= 0 {
		ops = append(ops, helpers.InsertAtOp{
			Value:  "🐈",
			PosIdx: extractedRes.Suffix.StartIndex,
		})
	}
	if extractedRes.Suffix.EndIndex >= 0 {
		ops = append(ops, helpers.InsertAtOp{
			Value:  "🐈",
			PosIdx: extractedRes.Suffix.EndIndex,
		})
	}
	annotatedContent := helpers.InsertAt(content, ops...)
	fmt.Println(annotatedContent)
}

func showGenerateExtracted(m *versions.Manager, cfg versions.AppConfig, content, value string) {
	res, err := m.Generate(cfg, content, value)
	if err != nil {
		fmt.Printf("Error generating extracted value for URL %s: %v\n", cfg.URL, err)
		return
	}
	var prefixes []string
	for _, prefix := range res.Prefixes {
		prefixes = append(prefixes, fmt.Sprintf("%q", prefix.Value))
	}
	fmt.Printf("Generated prefixes: %s, suffix: %q for version %q\n",
		prefixes,
		res.Suffix.Value,
		res.Match.Value,
	)
}
