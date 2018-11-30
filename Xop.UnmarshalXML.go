package paperfishGo

import (
   "encoding/xml"
)

func (x *Xop) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
   var v string
   var attr xml.Attr
   var tk xml.Token
   var err error
   var ok bool
   var iface interface{}

   tk, err = d.Token()
   if err != nil {
      Goose.Type.Logf(1,"Unmarshal error extracting token: %s", err)
      return err
   }

   if start, ok = tk.(xml.StartElement); ok {
      if start.Name.Local == "Include" {
         for _, attr = range start.Attr {
            if attr.Name.Local == "href" {
               iface, ok = xopParts.Load(attr.Value[4:])
               if !ok {
                  Goose.Type.Logf(1,"Unmarshal error part referenced as %s not found: %s", attr.Value[4:])
                  return err
               }

               *x = Xop(iface.(string))
               xopParts.Delete(attr.Value[4:])

               tk, err = d.Token()
               if err != nil {
                  Goose.Type.Logf(1,"Unmarshal error extracting ending token: %s", err)
                  return err
               }

               tk, err = d.Token()
               if err != nil {
                  Goose.Type.Logf(1,"Unmarshal error extracting last ending token: %s", err)
                  return err
               }
               return nil
            }
         }
      }
   }

   d.DecodeElement(&v, &start)
   *x = Xop(v)

   return nil
}

