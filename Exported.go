package paperfishGo

import (
	"strings"
)

func Exported(s string) (string, error) {
	if len(s) == 0 {
		return "", ErrEmptyString
	}

	return strings.ToUpper(s[:1]) + s[1:], nil
}
