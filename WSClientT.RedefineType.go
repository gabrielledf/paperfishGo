package paperfishGo

import (
   "strings"
)

func (ws WSClientT) RedefineType(typeName, typeDef string) error {
   var nm string
   var fld ElementT
   var sym *XsdSymT
   var attr AttributeT
   var xpath []string
   var i int
   var ok bool
   var typ *ComplexTypeT

   if ws.symtab != nil {
      xpath = strings.Split(typeName,".")

      if len(xpath) == 1 {
         for nm, sym = range ws.symtab {
            if typeName!=nm || sym.xsdref==nil {
               continue
            }

            if _, ok = sym.xsdref.(*SimpleTypeT); ok {
               sym.xsdref.(*SimpleTypeT).RestrictionBase.Base = typeDef
               return nil
            }
         }
      } else {
         for nm, sym = range ws.symtab {
            if xpath[0]!=nm || sym.xsdref==nil {
               continue
            }

            if typ, ok = sym.xsdref.(*ComplexTypeT); ok {
               for i, attr = range typ.Attribute {
                  if xpath[1]==attr.Name {
                     sym.xsdref.(*ComplexTypeT).Attribute[i].Type = typeDef
                     return nil
                  }
               }

               for i, fld = range typ.Sequence {
                  if xpath[1]==fld.Name {
                     sym.xsdref.(*ComplexTypeT).Sequence[i].Type = typeDef
                     return nil
                  }
               }
            }
         }
      }
   }

   return ErrParmNotFound
}
