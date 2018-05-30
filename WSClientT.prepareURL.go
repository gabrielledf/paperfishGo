package paperfishGo

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func (ws *WSClientT) prepareURL(op *OperationT, opId string, scheme string, inputDef []*ParameterT, inputValues map[string]interface{}) (string, error) {
	var targetURI string
	var p *ParameterT
	var ok bool
	var val interface{}

	targetURI = fmt.Sprintf("%s://%s/%s/%s", scheme, ws.Host, ws.BasePath, op.Path)

	for _, p = range inputDef {
		if val, ok = inputValues[p.Name]; !ok {
			Goose.Fetch.Logf(1, "%s %s for %s", ErrParmNotFound, p.Name, opId)
			return "", ErrParmNotFound
		}

		if reflect.ValueOf(val).Kind() != p.Kind {
			Goose.Fetch.Logf(1, "%s of %s for %s", ErrWrongParmType, p.Name, opId)
			return "", ErrWrongParmType
		}

		targetURI = strings.Replace(targetURI, "{"+p.Name+"}", url.PathEscape(fmt.Sprintf("%v", val)), -1)
	}

	return targetURI, nil
}
