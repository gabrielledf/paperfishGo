package paperfishGo

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func (ws *WSClientT) prepareQuery(op *OperationT, opId string, scheme string, inputDef []*ParameterT, inputValues map[string]interface{}, targetURI string) (string, error) {
	var p *ParameterT
	var ok bool
	var val interface{}
	var parms []string

	if len(inputDef) == 0 {
		return targetURI, nil
	}

	for _, p = range inputDef {
		if val, ok = inputValues[p.Name]; !ok {
			Goose.Fetch.Logf(1, "%s %s for %s", ErrParmNotFound, p.Name, opId)
			return "", ErrParmNotFound
		}

		if reflect.ValueOf(val).Kind() != p.Kind {
			Goose.Fetch.Logf(1, "%s of %s for %s", ErrWrongParmType, p.Name, opId)
			return "", ErrWrongParmType
		}

		parms = append(parms, opId+"="+url.QueryEscape(fmt.Sprintf("%v", val)))
	}

	return fmt.Sprintf("%s?%s", targetURI, strings.Join(parms, "&")), nil
}
