package paperfishGo

import (
   "fmt"
)


func (ws WSClientT) Types() string {
   var ret string
   var err error
   var nm, name string
   var tag string
   var fld ElementT
   var sym *XsdSymT
   var fldType string
   var fldname string
   var attr AttributeT

   if ws.symtab == nil {
      fmt.Printf("ws==nil\n")
   } else {
      for nm, sym = range ws.symtab {
         if sym.xsdref == nil {
            continue
         }
         Goose.Type.Logf(1,"sym.Name: %s", nm)
         name, err = Exported(nm)
         if err != nil {
            Goose.Type.Logf(1,"Error exporting %s: %s", nm, err)
            continue
         }
         name += "T"
         switch typ := sym.xsdref.(type) {
//            case reflect.Array, reflect.Slice:
            case *SimpleTypeT:
               fldType = bName(typ.RestrictionBase.Base)
               if ws.symtab[fldType].xsdref != nil {
                  fldType, err = Exported(bName(attr.Type))
                  if err != nil {
                     Goose.Type.Logf(1,"Error exporting fieldtype %s: %s", tag, err)
                     continue
                  }
                  fldType += "T"
               }

               ret += "type " + name + " " + fldType + "\n\n"
            case *ComplexTypeT:
               ret += "type " + name + " struct{\n"

               for _, attr = range typ.Attribute {
                  fldname, err = Exported(attr.Name)
                  if err != nil {
                     Goose.Type.Logf(1,"Error exporting field name %s: %s", fldname, err)
                     continue
                  }

                  fldType = bName(attr.Type)
                  if ws.symtab[fldType].xsdref != nil {
                     fldType, err = Exported(bName(attr.Type))
                     if err != nil {
                        Goose.Type.Logf(1,"Error exporting fieldtype %s: %s", tag, err)
                        continue
                     }
                     fldType += "T"
                  }

                  tag = bName(attr.Name)
                  tag = " `xml:\"" + tag + ",attr\" json:\"" + tag + "\"`"

                  ret += IndentPrefix +
                         fldname + " " +
                         fldType +
                         tag + "\n"
               }

               for _, fld = range typ.Sequence {
                  fldname, err = Exported(fld.Name)
                  if err != nil {
                     Goose.Type.Logf(1,"Error exporting field name %s: %s", fldname, err)
                     continue
                  }

                  fldType = bName(fld.Type)
                  if ws.symtab[fldType].xsdref != nil {
                     fldType, err = Exported(bName(fld.Type))
                     if err != nil {
                        Goose.Type.Logf(1,"Error exporting fieldtype %s: %s", tag, err)
                        continue
                     }
                     fldType += "T"
                  }

                  tag = bName(fld.Name)
                  if fld.Nillable == "true" {
                     tag +=  ",omitempty"
                  }
                  tag = " `xml:\"" + tag + "\" json:\"" + tag + "\"`"

                  if fld.MaxOccurs != "" {
                     fldType = "[]" + fldType
                  } else {
                     fldType = "*" + fldType
                  }

                  ret += IndentPrefix +
                         fldname + " " +
                         fldType +
                         tag + "\n"
               }
               ret += "}\n\n"
         }
      }
   }

   return ret
}

