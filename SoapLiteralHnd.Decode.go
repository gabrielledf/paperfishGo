package paperfishGo

import (
	"encoding/json"
	"io"
)

func (Hand SoapLiteralHnd) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
