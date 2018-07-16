package paperfishGo

import (
	"bytes"
	"encoding/xml"
)

func (Hand SoapLiteralHnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {

	return xml.NewEncoder(Writer{buf: w}).Encode(v)
}
