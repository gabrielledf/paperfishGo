package paperfishGo

import (
   "os"
   "fmt"
   "bytes"
   "encoding/xml"
)

func (Hand SoapLiteralHnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {
   var buf []byte
   var err error
   var vv SoapData
   var vvv interface{}
   var attr []xml.Attr
   var xmlns, alias string

   fmt.Printf(`<?xml version="1.0" encoding="UTF-8"?>`)
   fmt.Printf("\n\nxmlnsrev: %#v\n\n", Hand.ws.Xmlns)

   attr = make([]xml.Attr,0,len(Hand.ws.Xmlns) + 1)
   for xmlns, alias = range Hand.ws.Xmlns {
      attr = append(attr,xml.Attr{
         Name: xml.Name{Local: "xmlns:" + alias},
         Value: xmlns,
      })
   }
   attr = append(attr,xml.Attr{
      Name: xml.Name{Local: "xmlns:tns"},
      Value: Hand.ws.TargetNamespace,
   })

   //ns = Hand.ws.Xmlns[Hand.ws.TargetNamespace]
   vv = v.(SoapData)
   vvv = vv.SetName("tns:" + nm, attr)
   envelope.Body = SoapBodyT{
      Data: vvv,
   }
   buf, err = xml.Marshal(envelope)
   if err != nil {
      Goose.Fetch.Logf(0,"Error marshaling v: %s",err)
      os.Exit(0)
   }

   fmt.Printf("%s",buf)


   os.Exit(0)
   return xml.NewEncoder(Writer{buf: w}).Encode(v)
}
