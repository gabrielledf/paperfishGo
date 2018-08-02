package paperfishGo

import (
	"fmt"
	"bytes"
	"reflect"
)

func (ws *WSClientT) prepareBody(op *OperationT, opId string, inputDef []*ParameterT, inputValues map[string]interface{}) ([]byte, error) {
	var postdata bytes.Buffer
	var p *ParameterT
	var ok bool
	var val interface{}
	var err error
	var isTail bool
	var kind reflect.Kind

	Goose.Fetch.Logf(4, "Preparing body")
	for _, p = range inputDef {
		Goose.Fetch.Logf(4, "Preparing input %#v", p)
		if val, ok = inputValues[p.Name]; !ok {
			Goose.Fetch.Logf(1, "%s %s for %s", ErrParmNotFound, p.Name, opId)
			return nil, ErrParmNotFound
		}

		if reflect.ValueOf(val).Kind() != p.Kind {
			Goose.Fetch.Logf(1, "%s of %s for %s", ErrWrongParmType, p.Name, opId)
			return nil, ErrWrongParmType
		}

		Goose.Fetch.Logf(4, "Encoding with %#v", op.Encoder)

		kind = reflect.TypeOf(val).Kind()
		if kind == reflect.Struct {
			err = op.Encoder.Encode(&postdata, p.Name, val, isTail)
		} else {
			err = op.Encoder.Encode(&postdata, p.Name, fmt.Sprintf("%v", val), isTail)
		}
		if err != nil {
			Goose.Fetch.Logf(1, "%s:%s of %s for %s", ErrBuffer, err, p.Name, opId)
			return nil, err
		}

		isTail = true
	}

	Goose.Fetch.Logf(6, "Preparing body done: %s", postdata.Bytes())
	return postdata.Bytes(), nil
}
