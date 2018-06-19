package paperfishGo

import (
	"fmt"
	"net/url"
	"reflect"
)

func (wsc *WSockClientT) Send(opName string, params map[string]interface{}, fn interface{}, fnfail func(int)) error {
	var parms []interface{}
	var fnval reflect.Value
	var p *ParameterT
	var ok bool
	var val interface{}
	var subop *SubOperationT

	fnval = reflect.ValueOf(fn)

	if fnval.Kind() != reflect.Func {
		Goose.Fetch.Logf(1, "Error %s for %s callback function", ErrWrongParmType, opName)
		return ErrWrongParmType
	}

	Goose.Fetch.Logf(4, "Will send data to %s through websocket", opName)

	if wsc==nil || wsc.SubOperations == nil {
		Goose.Fetch.Logf(1, "Error %s for %s suboperation", ErrNilHandle, opName)
		return ErrNilHandle
	}

	if subop, ok = wsc.SubOperations[opName]; !ok {
		Goose.Fetch.Logf(1, "Error %s for %s suboperation", ErrUnknownOperation, opName)
		return ErrUnknownOperation
	}

	for _, p = range subop.Parms {
		if val, ok = params[p.Name]; !ok {
			Goose.Fetch.Logf(1, "%s %s of %s", ErrParmNotFound, p.Name, opName)
			return ErrParmNotFound
		}

		if reflect.ValueOf(val).Kind() != p.Kind {
			Goose.Fetch.Logf(1, "%s %s of %s", ErrWrongParmType, p.Name, opName)
			return ErrWrongParmType
		}

		Goose.Fetch.Logf(5, "Adding parm %s.%s=%v", opName, p.Name, val)
		parms = append(parms, url.QueryEscape(fmt.Sprintf("%v", val)))
	}

	Goose.Fetch.Logf(4, "Params added")

	go func() {
		Goose.Fetch.Logf(4, "Params will be sent")
		wsc.cli2srvch <- WSockRequest{
			SubOperation: opName,
			Params:       parms,
			CallbackT: CallbackT{
				Callback:     fnval,
				FailCallback: fnfail,
			},
		}
		Goose.Fetch.Logf(4, "Params sent")
	}()

	return nil
}
