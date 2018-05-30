package paperfishGo

import (
	"reflect"
	"strings"
)

var kindString reflect.Kind = reflect.ValueOf("").Kind()
var kindNumber reflect.Kind = reflect.ValueOf(0.1).Kind()
var kindInteger reflect.Kind = reflect.ValueOf(1).Kind()
var kindBoolean reflect.Kind = reflect.ValueOf(true).Kind()
var kindArray reflect.Kind = reflect.ValueOf([]int{}).Kind()
var kindObject reflect.Kind = reflect.ValueOf(struct{}{}).Kind()

func getKind(t string) (reflect.Kind, error) {
	Goose.Fetch.Logf(4, "Type is %s", strings.ToLower(t))

	switch strings.ToLower(t) {
	case "string":
		return kindString, nil
	case "number":
		return kindNumber, nil
	case "integer":
		return kindInteger, nil
	case "boolean":
		return kindBoolean, nil
	case "array":
		return kindArray, nil
	case "object":
		return kindObject, nil
	}

	return kindString, ErrUnknownKind
}
