package paperfishGo

import (
	"reflect"
)

func (wsc *WSockClientT) Send(opName string, params map[string]interface{}, fn interface{}) error {
	var parms []interface{}
	var fnval reflect.Value

	fnval = reflect.ValueOf(fn)

	if fnval.Kind() != reflect.Func {
		Goose.Fetch.Logf(1, "Error %s for %s callback function", ErrWrongParmType, opName)
		return ErrWrongParmType
	}

	go func() {
		wsc.cli2srvch <- WSockRequest{
			SubOperation: opName,
			Params:       parms,
			Callback:     fnval,
		}
	}()

	return nil
}
