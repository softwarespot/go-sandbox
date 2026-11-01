package versions

import "strconv"

func quote(s string) string {
	return strconv.Quote(s)
}

func unquote(s string) string {
	us, err := strconv.Unquote(s)
	if err != nil {
		return s
	}
	return us
}
