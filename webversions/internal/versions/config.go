package versions

type AppConfig struct {
	ID            string
	Name          string
	URL           string
	Info          string
	CurrVersion   string
	WebVersion    string
	Prefixes      []string
	Suffix        string
	SearchFromEnd bool
	TabNames      []string
}
