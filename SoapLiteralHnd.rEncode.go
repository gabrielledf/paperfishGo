package paperfishGo

import (
   "fmt"
)

func (Hand SoapLiteralHnd) rEncode(nm string, v interface{}) (string, error) {
   var ok bool

   switch Hand.symtab[nm].Type {
      case typeOfiface:
      case typeOfstring:
         if _, ok = v.(string); !ok {
            return fmt.Sprintf("<%s>%s</%s>",nm,v,nm), nil
         }
      case typeOfboolean:
      case typeOfint:
      case typeOffloat:
      case typeOfduration:
      case typeOftime:
      case typeOfBinary:
      default:
   }




   return "", ErrWrongParmType
}
