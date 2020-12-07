package paperfishGo

import (
   "reflect"
   "strings"
   "github.com/luisfurquim/stonelizard"
)

func getType(t stonelizard.SwaggerSchemaT) (string, reflect.Type, error) {
   Goose.Fetch.Logf(4, "Type is %s", strings.ToLower(t.Type))

   switch strings.ToLower(t.Type) {
   case "string", "file":
      return t.Title, reflect.TypeOf(""), nil
   case "number":
      return t.Title, reflect.TypeOf(float64(.1)), nil
   case "integer":
      return t.Title, reflect.TypeOf(int64(1)), nil
   case "boolean":
      return t.Title, reflect.TypeOf(true), nil
   case "object":
      return getStruct(t.Title, t.Properties)
   case "array":
      if t.XKeyType!="" {
         return getMap(t.XKeyType, t.Items)
      } else {
         return getArray(t.Items)
      }
   }

   return "", reflect.TypeOf(""), ErrUnknownKind
}

func getStruct(title string, fields map[string]stonelizard.SwaggerSchemaT) (string, reflect.Type, error) {
   var t reflect.Type
   var tit string
   var err error
   var k string
   var v stonelizard.SwaggerSchemaT
   var flds []reflect.StructField
   var jsonName string

   for k, v = range fields {
      jsonName = ""
      if k[0:1] == strings.ToLower(k[0:1]) {
         jsonName = `json:"` + k + `"`
         k = strings.ToUpper(k[0:1]) + k[1:]
      }
      tit, t, err = getType(v)
      if err!=nil {
         return title, t, err
      }

      flds = append(flds, reflect.StructField{
         Name: k,
         Type: t,
         Tag: reflect.StructTag(`typename:"` + tit + `"` + jsonName),
      })
   }

   return title, reflect.StructOf(flds), nil
}

func getArray(items *stonelizard.SwaggerSchemaT) (string, reflect.Type, error) {
   var t reflect.Type
   var err error
   var title string

   if items == nil {
      return "", reflect.ArrayOf(0,reflect.TypeOf(true)), nil
   }

   title, t, err = getType(*items)
   if err!=nil {
      return title, t, err
   }

   return items.Title, reflect.ArrayOf(0,t), nil
}

func getMap(key string, items *stonelizard.SwaggerSchemaT) (string, reflect.Type, error) {
   var t reflect.Type
   var err error
   var ktype reflect.Type
   var title string

   title, ktype, err = getType(stonelizard.SwaggerSchemaT{
      Type: key,
   })
   if err!=nil {
      return title, t, err
   }

   title, t, err = getType(*items)
   if err!=nil {
      return title, t, err
   }

   return items.Title, reflect.MapOf(ktype,t), nil
}
