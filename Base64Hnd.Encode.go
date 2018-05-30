package paperfishGo

import (
	"bytes"
	"encoding/base64"
	"io"
)

func (Hand Base64Hnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {
	var err error
	var n, l int
	var out io.WriteCloser

	out = base64.NewEncoder(base64.StdEncoding, Writer{buf: w})
	defer out.Close()

	switch target := v.(type) {
	case *string:
		n, err = out.Write([]byte(*target))
		l = len(*target)
	case string:
		n, err = out.Write([]byte(target))
		l = len(target)
	case []byte:
		n, err = out.Write(target)
		l = len(target)
	default:
		Goose.Fetch.Logf(1, "%s: %T", ErrWrongParmType, v)
		return ErrWrongParmType
	}

	if n != l {
		return ErrWrite
	}

	return err
}
