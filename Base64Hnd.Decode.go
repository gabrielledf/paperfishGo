package paperfishGo

import (
	"encoding/base64"
	"io"
)

func (Hand Base64Hnd) Decode(r io.Reader, v interface{}) error {
	return BinaryHnd{}.Decode(base64.NewDecoder(base64.StdEncoding, r), v)
}
