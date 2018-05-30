package paperfishGo

import (
	"os"
	"reflect"
)

func Set(v reflect.Value, val reflect.Value) {
	var fld reflect.StructField
	var i int
	var length int
	var last int
	var elemType reflect.Type
	var valdata reflect.Value
	var mapKey reflect.Value
	var typ reflect.Type
	var vptr reflect.Value

	typ = v.Type()
	Goose.Fetch.Logf(5, "Type before indirect %v", typ)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if !v.IsValid() || v.IsNil() {
			v.Set(reflect.New(typ))
			Goose.Fetch.Logf(5, "v nil: %#v", v)
		}
	}
	Goose.Fetch.Logf(5, "Type after indirect %v", typ)

	if val.Kind() == reflect.Interface {
		Goose.Fetch.Logf(5, "val deinterfaced: %v => %v", val, val.Elem())
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		Goose.Fetch.Logf(5, "struct")
		for i = 0; i < typ.NumField(); i++ {
			fld = typ.Field(i)
			Goose.Fetch.Logf(5, "val %s: %T -> %#v", fld.Name, val.Interface(), val.Interface())
			if fld.Anonymous {
				Goose.Fetch.Logf(5, "Going to Struct for %#v", val)
				Set(reflect.Indirect(v).Field(i).Addr(), val)
				//            Set(reflect.Indirect(v).Field(i).Addr(),val)
			} else {
				if val.Kind() == reflect.Interface {
					val = val.Elem()
				}
				Goose.Fetch.Logf(5, "val %s: %T -> %#v", fld.Name, val.Interface(), val.Interface())
				valdata = val.MapIndex(reflect.ValueOf(fld.Name))
				Goose.Fetch.Logf(5, "Verifying Struct field %s = %#v", fld.Name, valdata)
				if valdata.Kind() != reflect.Invalid {
					Goose.Fetch.Logf(5, "Initializing Struct for key %s = %#v", fld.Name, valdata)
					Set(reflect.Indirect(v).Field(i).Addr(), valdata)
				}
			}
		}
		Goose.Fetch.Logf(4, "Struct value set to %#v %T <= %#v %T", v.Elem(), v.Elem().Interface(), val, val.Interface())
	} else if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
		elemType = typ.Elem()
		Goose.Fetch.Logf(6, "Array/slice:%s", typ.Kind())
		Goose.Fetch.Logf(6, "val:%v", val)

		if val.IsValid() {
			if reflect.Indirect(v).Len() < val.Len() {
				last = reflect.Indirect(v).Len()
				length = val.Len() - last
				reflect.Indirect(v).Set(reflect.AppendSlice(reflect.Indirect(v), reflect.MakeSlice(typ, length, length)))
			}

			for i = 0; i < val.Len(); i++ {
				Goose.Fetch.Logf(6, "Initializing Array/Slice set for key index %d = %#v", i, val.Index(i))
				// Providing pointer to element
				if reflect.Indirect(v).Index(i).Kind() == reflect.Ptr {
					Set(reflect.Indirect(v).Index(i), val.Index(i))
				} else {
					Set(reflect.Indirect(v).Index(i).Addr(), val.Index(i))
				}
			}
			Goose.Fetch.Logf(4, "Array value set to %#v %T <= %#v %T", v.Elem(), v.Elem().Interface(), val, val.Interface())
		} else {
			Goose.Fetch.Logf(3, "Invalid val:%v", val)
		}
	} else if typ.Kind() == reflect.Map {
		Goose.Fetch.Logf(5, "map")
		elemType = typ.Elem()
		if reflect.Indirect(v).IsNil() && !val.IsNil() {
			reflect.Indirect(v).Set(reflect.MakeMap(typ))
		}

		Goose.Fetch.Logf(6, "Map:%s~%s", typ, elemType)

		for _, mapKey = range val.MapKeys() {
			Goose.Fetch.Logf(6, "Initializing Map set for key %#v on %#v", val.MapIndex(mapKey), reflect.Indirect(v).MapIndex(mapKey))
			if !reflect.Indirect(v).MapIndex(mapKey).IsValid() {
				vptr = reflect.New(elemType.Elem())
				reflect.Indirect(v).SetMapIndex(mapKey, vptr)
			}
			Set(reflect.Indirect(v).MapIndex(mapKey), val.MapIndex(mapKey))
		}
	} else if typ.Kind() == reflect.Func || typ.Kind() == reflect.Chan || typ.Kind() == reflect.UnsafePointer {
		Goose.Fetch.Logf(6, "Skipping configuration of func/channel/unsafe pointer for key")
	} else {
		//      Goose.Fetch.Logf(4,"Scalar value setting: %#v <= %s", v, val)
		for val.Kind() == reflect.Ptr {
			Goose.Fetch.Logf(4, "Scalar value indirect setting: %#v <= %#v", v, val.Elem())
			val = val.Elem()
		}

		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}

		Goose.Fetch.Logf(4, "Scalar value setting: %#v %T <= %#v %T", v, v.Interface(), val, val.Interface())
		if val.Type().AssignableTo(typ) {
			reflect.Indirect(v).Set(val)
			Goose.Fetch.Logf(4, "Scalar value set to %#v %T <= %#v %T", v.Elem(), v.Elem().Interface(), val, val.Interface())
		} else if val.Type().ConvertibleTo(typ) {
			reflect.Indirect(v).Set(val.Convert(typ))
			Goose.Fetch.Logf(4, "Scalar value converted to %#v %T <= %#v %T", v.Elem(), v.Elem().Interface(), val, val.Interface())
		} else {
			Goose.Fetch.Logf(4, "Scalar value not assignable nor convertible: %#v %T <= %#v %T", v, v.Interface(), val, val.Interface())
			os.Exit(0)
		}
	}
}
