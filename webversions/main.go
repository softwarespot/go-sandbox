package main

import (
	"flag"
	"time"
)

type cliFlags struct {
	ConfigPath string
	UserAgent  string
	Timeout    time.Duration
	CLI        bool
	Extra      bool
}

func main() {
	opts := cliFlags{
		ConfigPath: "WebVersions.txt",
		Timeout:    10 * time.Second,
		UserAgent:  "WebVersions/1.0",
		CLI:        true,
		Extra:      false,
	}
	flag.StringVar(&opts.ConfigPath, "config", opts.ConfigPath, "path to WebVersions configuration file")
	flag.DurationVar(&opts.Timeout, "timeout", opts.Timeout, "HTTP request timeout")
	flag.StringVar(&opts.UserAgent, "user-agent", opts.UserAgent, "User-Agent header for downloading content")
	flag.BoolVar(&opts.CLI, "cli", opts.CLI, "run in CLI mode instead of GUI")
	flag.BoolVar(&opts.Extra, "extra", opts.Extra, "extra information for debugging")
	flag.Parse()

	var err error
	if opts.CLI {
		err = runCLI(opts)
	} else {
		err = runGUI(opts)
	}
	if err != nil {
		panic(err)
	}
}
