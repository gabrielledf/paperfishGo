package paperfishGo

import (
	"bytes"
	"encoding/json"
)

func (Hand JsonHnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {
	return json.NewEncoder(Writer{buf: w}).Encode(v)
}
