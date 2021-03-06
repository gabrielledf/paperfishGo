package paperfishGo

import (
   "fmt"
)


func (ws WSClientT) Types(pkg string) string {
   var ret string
   var err error
   var nm, name string
   var tag string
   var fld ElementT
   var sym *XsdSymT
   var fldType string
   var newFldType string
   var fldname string
   var attr AttributeT
   var ns string
   var ok bool
   var nsrev map[string]string
   var isMainMesg bool
   var isOutputMesg bool
   var oper *OperationT
   var impXML, impTime, impPaper bool
   var hdr string

   if ws.symtab != nil {
      hdr = "package " + pkg + "\n\n"
      nsrev = map[string]string{}
      for nm, sym = range ws.symtab {
         if sym.xsdref == nil {
            continue
         }
         Goose.Type.Logf(1,"sym.Name: %s", nm)

         isMainMesg = false
         isOutputMesg = false
         for _, oper = range ws.PostOperation {
            Goose.Type.Logf(1,"in: %s, out: %s", oper.inMesg, oper.outMesg)
            if nm == oper.inMesg {
               isMainMesg = true
               break
            }
            if nm == oper.outMesg {
               isOutputMesg = true
               break
            }
         }

         for _, oper = range ws.PostOperation {
            Goose.Type.Logf(1,"in: %s, out: %s", oper.inMesg, oper.outMesg)
         }

         name, err = Exported(nm)
         if err != nil {
            Goose.Type.Logf(1,"Error exporting %s: %s", nm, err)
            continue
         }
         name += "T"
         if sym.ns == "" {
            ns = ""
         } else {
            ns = sym.ns + " "
         }
         switch typ := sym.xsdref.(type) {
//            case reflect.Array, reflect.Slice:
            case *SimpleTypeT:
               fldType = bName(typ.RestrictionBase.Base)
               if ws.symtab[fldType].xsdref != nil {
                  fldType, err = Exported(bName(attr.Type)) // bug: unitialized var, cut/paste problem
                  if err != nil {
                     Goose.Type.Logf(1,"Error exporting fieldtype %s: %s", tag, err)
                     continue
                  }
                  fldType += "T"
               } else if _, ok = xsd2go[fldType]; ok {
                  fldType = xsd2go[fldType]
                  if (fldType == "paperfishGo.Base64Binary") || (fldType == "paperfishGo.Xop")  {
                     impPaper = true
                  } else if fldType == "time.Time" {
                     impTime = true
                  }
               }

               ret += "type " + name + " " + fldType + "\n\n"
            case *ComplexTypeT:
               ret += "type " + name + " struct{\n"
               if isMainMesg {
                  impXML = true
                  ret += "   XMLName xml.Name\n"
                  ret += "   XMLAttr []xml.Attr `xml:\",attr,any\"`\n"
               }

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
                  } else if _, ok = xsd2go[fldType]; ok {
                     fldType = xsd2go[fldType]
                     if (fldType == "paperfishGo.Base64Binary") || (fldType == "paperfishGo.Xop")  {
                        impPaper = true
                     } else if fldType == "time.Time" {
                        impTime = true
                     }
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
                  Goose.Type.Logf(1,"Checking fldType %s on symtab", fldType)
                  if ws.symtab[fldType].xsdref != nil {
                     newFldType, err = Exported(bName(fld.Type))
                     if err != nil {
                        Goose.Type.Logf(1,"Error exporting fieldtype %s: %s", tag, err)
                        continue
                     }
                     newFldType += "T"

                     if fld.MaxOccurs != "" {
                        newFldType = "[]" + newFldType
                     } else if _, ok = ws.symtab[fldType].xsdref.(*ComplexTypeT) ; ok {
                        newFldType = "*" + newFldType
                     }
                     fldType = newFldType
                  } else if _, ok = xsd2go[fldType]; ok {
                     fldType = xsd2go[fldType]
                     if (fldType == "paperfishGo.Base64Binary") || (fldType == "paperfishGo.Xop")  {
                        impPaper = true
                     } else if fldType == "time.Time" {
                        impTime = true
                     }
                  }

                  tag = bName(fld.Name)
                  if fld.Nillable == "true" {
                     tag +=  ",omitempty"
                  }
                  if isOutputMesg {
                     tag = " `xml:\"Body>" + nm + ">" + tag + "\" json:\"" + tag + "\"`"
                  } else {
                     tag = " `xml:\"" + ns + tag + "\" json:\"" + tag + "\"`"
                  }

                  ret += IndentPrefix +
                         fldname + " " +
                         fldType +
                         tag + "\n"
               }
               ret += "}\n\n"

               if isMainMesg {
                  ret += "func (soapdata " + name + ") SetName(nm string, attr []xml.Attr) (interface{}) {\n"
                  ret += "   soapdata.XMLName.Local = nm\n"
                  ret += "   soapdata.XMLAttr = attr\n"
                  ret += "   return soapdata\n"
                  ret += "}\n\n"
               }
         }
      }

      ret += fmt.Sprintf("var Xmlns map[string]string = %#v\n\n",nsrev)
      ret += fmt.Sprintf("var Tns string = %#v\n\n",ws.TargetNamespace)
   }


   if impXML  {
      hdr += `import "encoding/xml"` + "\n\n"
   }
   if impTime {
      hdr += `import "time"` + "\n\n"
   }
   if impPaper {
      hdr += `import "github.com/gabrielledf/paperfishGo"` + "\n\n"
   }


   return hdr + ret
}

