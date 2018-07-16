package paperfishGo

import (
	"strings"
)

func bName(nm string) string {
	var pair []string
	pair = strings.Split(nm, ":")
	return pair[len(pair)-1]
}
