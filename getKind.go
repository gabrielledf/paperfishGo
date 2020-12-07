package paperfishGo

import (
   "reflect"
   "strings"
)

func getKind(t string) (reflect.Kind, error) {
   Goose.Fetch.Logf(4, "Type is %s", strings.ToLower(t))

   switch strings.ToLower(t) {
   case "string", "file":
      return reflect.String, nil
   case "number":
      return reflect.Float64, nil
   case "integer":
      return reflect.Int64, nil
   case "boolean":
      return reflect.Bool, nil
   case "array":
      return reflect.Array, nil
   case "object":
      return reflect.Struct, nil
   }

   return reflect.String, ErrUnknownKind
}
