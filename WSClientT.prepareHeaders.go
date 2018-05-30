package paperfishGo

import (
	"fmt"
	"net/http"
	"reflect"
)

func (ws *WSClientT) prepareHeaders(op *OperationT, opId string, inputDef []*ParameterT, inputValues map[string]interface{}, req *http.Request) error {
	var p *ParameterT
	var ok bool
	var val interface{}

	for _, p = range inputDef {
		if val, ok = inputValues[p.Name]; !ok {
			Goose.Fetch.Logf(1, "%s %s for %s", ErrParmNotFound, p.Name, opId)
			return ErrParmNotFound
		}

		if reflect.ValueOf(val).Kind() != p.Kind {
			Goose.Fetch.Logf(1, "%s of %s for %s", ErrWrongParmType, p.Name, opId)
			return ErrWrongParmType
		}

		req.Header.Set(p.Name, fmt.Sprintf("%v", val))
	}

	return nil
}
