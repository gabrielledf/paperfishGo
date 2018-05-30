package paperfishGo

import (
	"bytes"
	"fmt"
)

func (Hand TextHnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {
	var err error
	var n, l int

	switch target := v.(type) {
	case *string:
		n, err = w.Write([]byte(*target))
		l = len(*target)
	case []byte:
		n, err = w.Write(target)
		l = len(target)
	default:
		fmt.Printf("[%#v\n", v)
		err = ErrWrongParmType
	}

	if err == nil && n != l {
		err = ErrWrite
	}

	return err
}
