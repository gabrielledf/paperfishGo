package paperfishGo

import (
	"io"
	"io/ioutil"
)

func (Hand BinaryHnd) Decode(r io.Reader, v interface{}) error {
	var buf []byte
	var err error

	switch target := v.(type) {
	case *string:
		buf, err = ioutil.ReadAll(r)
		if err == nil {
			*target = string(buf)
		}
	case []byte:
		target, err = ioutil.ReadAll(r)
	default:
		err = ErrWrongParmType
	}

	return err
}
