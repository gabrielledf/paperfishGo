package paperfishGo

import (
	"reflect"
)

func (wsc *WSockClientT) On(evtName string, fn interface{}) error {
	var fnval reflect.Value

	fnval = reflect.ValueOf(fn)

	if fnval.Kind() != reflect.Func {
		Goose.Fetch.Logf(1, "Error %s for %s callback function", ErrWrongParmType, evtName)
		return ErrWrongParmType
	}

	go func() {
		wsc.bindch <- WSockRequest{
			SubOperation: evtName,
			CallbackT: CallbackT{
				Callback: fnval,
			},
		}
	}()

	return nil
}
