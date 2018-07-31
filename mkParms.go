package paperfishGo

import (
	"reflect"
)

func mkParms(fn reflect.Value, parms []interface{}) ([]reflect.Value, error) {
	var parmValues []reflect.Value
	var parmValue reflect.Value
	var fnType reflect.Type
	var i int

	fnType = fn.Type()

	if len(parms) != fnType.NumIn() {
		Goose.Fetch.Logf(1, "%s: parms: %#v /// fn: %#v", ErrWrongParmCount, parms, fnType.NumIn())
		return nil, ErrWrongParmCount
	}

	parmValues = make([]reflect.Value, len(parms))

	for i = 0; i < len(parms); i++ {
		parmValue = reflect.New(fnType.In(i))
		Goose.Fetch.Logf(0, "parm: %d = %#v", i, parmValue)
		Set(parmValue, reflect.ValueOf(parms[i]))
		parmValues[i] = parmValue.Elem()
	}

	return parmValues, nil
}
