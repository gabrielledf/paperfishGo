package paperfishGo

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func (Hand FormURLHnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {
	var val reflect.Value
	var out string
	var k reflect.Kind
	var s []string
	var err error
	var n, i int
	var keys []reflect.Value
	var t reflect.Type

	if isTail {
		out = "&"
	}

	val = reflect.ValueOf(v).Elem()
	k = val.Kind()
	for (k == reflect.Ptr) || (k == reflect.Interface) { // Let's dereference it
		val = val.Elem()
		k = val.Kind()
	}

	out += nm + "="

	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		out += url.QueryEscape(fmt.Sprintf("%v", val.Interface()))
	case reflect.Array, reflect.Slice:
		s = make([]string, val.Len())
		for i = 0; i < len(s); i++ {
			s[i] = fmt.Sprintf("%v", val.Index(i).Interface())
		}
		out += url.QueryEscape(strings.Join(s, ","))
	case reflect.Map:
		keys = val.MapKeys()
		s = make([]string, len(keys))
		for i = 0; i < len(s); i++ {
			s[i] = fmt.Sprintf("%s:%v", keys[i].Interface(), val.MapIndex(keys[i]).Interface())
		}
		out += url.QueryEscape(strings.Join(s, ","))
	case reflect.Struct:
		t = val.Type()
		s = make([]string, t.NumField())
		for i = 0; i < len(s); i++ {
			s[i] = fmt.Sprintf("%s:%v", t.Field(i).Name, val.Field(i).Interface())
		}
		out += url.QueryEscape(strings.Join(s, ","))
	}

	n, err = w.Write([]byte(out))
	if err != nil || n < len(out) {
		Goose.Fetch.Logf(1, "%s (%s)", ErrWrite, err)
	}

	if n < len(out) {
		Goose.Fetch.Logf(1, "%s (Too few bytes written)", ErrWrite)
	}

	return err
}
