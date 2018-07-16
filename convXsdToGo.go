package paperfishGo

import (
   "fmt"
   "os"
   "reflect"
   "strings"
)

func convXsdToGo(ct interface{}, xsdSymTab XsdSymTabT, nmStack string) (reflect.Type, error) {
   var ret, t reflect.Type
   var e ElementT
   var fields []reflect.StructField
   var err error
   var n string
   var sym *XsdSymT
   var ok bool
   var pos int
   var tag reflect.StructTag
   var a AttributeT

   //   Goose.New.Logf(1, "... %s", nmStack)

   switch typ := ct.(type) {
   case (*ComplexTypeT):
      for _, e = range typ.Sequence {
         if e.Type != "" {
            n = bName(e.Type)
            if strings.Contains(nmStack, n) {
               t = xsdSymTab["iface"].Type
               tag = reflect.StructTag(`xml:"` + e.Name + `" self:"` + n + `"`)
               //               Goose.New.Logf(1, "\n\n------------\n\n%s ... %s\n\n------------\n\n", n, nmStack)
            } else {
               Goose.New.Logf(4, "+++ complex %s", n)
               t, err = convXsdToGo(n, xsdSymTab, fmt.Sprintf("%s%c%s", nmStack, 0, n))
               if err == nil {
                  Goose.New.Logf(4, "--- complex %s", n)
                  xsdSymTab[n].Type = t
               }
               tag = reflect.StructTag(`xml:"` + e.Name + `"`)
            }
         } else if len(e.ComplexTypes) == 1 {
            t, err = convXsdToGo(&e.ComplexTypes[0], xsdSymTab, fmt.Sprintf("%s%c%s", nmStack, 0, bName(e.ComplexTypes[0].Name)))
            tag = reflect.StructTag(`xml:"` + e.Name + `"`)
         } else {
            err = ErrUndetectableType
         }

         if e.MaxOccurs != "" && e.MaxOccurs != "1" {
            t = reflect.SliceOf(t)
         }
         if err != nil {
            Goose.New.Logf(1, "Error converting sequence element to struct field: %s", err)
            return nil, err
         }

         n, err = Exported(e.Name)
         if err != nil {
            Goose.New.Logf(1, "Error converting sequence element name to struct field: %s", err)
            return nil, err
         }

         fields = append(fields, reflect.StructField{Name: n, Type: t, Tag: tag})
      }

      for _, a = range typ.Attribute {
         if a.Type != "" {
            n = bName(a.Type)
            if strings.Contains(nmStack, n) {
               t = xsdSymTab["iface"].Type
               tag = reflect.StructTag(`xml:"` + a.Name + `" self:"` + n + `"`)
               //               Goose.New.Logf(1, "\n\n------------\n\n%s ... %s\n\n------------\n\n", n, nmStack)
            } else {
               Goose.New.Logf(4, "+++ complex %s", n)
               t, err = convXsdToGo(n, xsdSymTab, fmt.Sprintf("%s%c%s", nmStack, 0, n))
               if err == nil {
                  Goose.New.Logf(4, "--- complex %s", n)
                  xsdSymTab[n].Type = t
               }
               tag = reflect.StructTag(`xml:"` + a.Name + `,attr"`)
            }
         } else {
            err = ErrUndetectableType
         }

         if err != nil {
            Goose.New.Logf(1, "Error converting sequence element to struct field: %s", err)
            return nil, err
         }

         n, err = Exported(a.Name)
         if err != nil {
            Goose.New.Logf(1, "Error converting sequence element name to struct field: %s", err)
            return nil, err
         }

         fields = append(fields, reflect.StructField{Name: n, Type: t, Tag: tag})
      }

      ret = reflect.StructOf(fields)
   case (*SimpleTypeT):
      if typ.List.ItemType != "" {
         n = bName(typ.List.ItemType)
         t, err = convXsdToGo(n, xsdSymTab, fmt.Sprintf("%s%c%s", nmStack, 0, n))
         if err != nil {
            Goose.New.Logf(1, "Error converting simple type itemtype to reflect.type: %s", err)
            return nil, err
         }
         t = reflect.SliceOf(t)
      } else {
         n = bName(typ.RestrictionBase.Base)
         t, err = convXsdToGo(n, xsdSymTab, fmt.Sprintf("%s%c%s", nmStack, 0, n))
         if err != nil {
            Goose.New.Logf(1, "Error converting simple type base to reflect.type: %s", err)
            return nil, err
         }
      }

      ret = t
   case string:
      n = bName(typ)

      sym, ok = xsdSymTab[n]
      if !ok {
         Goose.New.Logf(1, "Error unknown type reference '%s' in xsd", typ)
         return nil, err
      }

      pos = strings.LastIndex(nmStack, "\x00")
      if pos == -1 {
         pos = len(nmStack)
      }

      if strings.Contains(nmStack[:pos], n) {
         ret = xsdSymTab["iface"].Type
         tag = reflect.StructTag(`xml:"` + e.Name + `" self:"` + n + `"`)
         //Goose.New.Logf(1, "\n\n------------\n\n%s ... %s\n\n------------\n\n", n, nmStack)
         os.Exit(0)
      } else {
         if sym.Type != nil {
            ret = sym.Type
         } else if sym.xsdref != nil {
            ret, err = convXsdToGo(sym.xsdref, xsdSymTab, fmt.Sprintf("%s%c%s", nmStack, 0, n))
            if err != nil {
               Goose.New.Logf(1, "Error following xsd type definition: %s", err)
               return nil, err
            }
         } else {
            Goose.New.Logf(1, "Error unknown type reference '%s' in xsd", typ)
            return nil, err
         }

      }

   }

   return ret, nil
}
