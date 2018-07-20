package paperfishGo

import (
   "bytes"
   "encoding/xml"
)

func (Hand SoapLiteralHnd) Encode(w *bytes.Buffer, nm string, v interface{}, isTail bool) error {
   var vv SoapData
   var vvv interface{}
   var attr []xml.Attr
   var xmlns, alias string
   var envelope soapEnvelopeT = soapEnvelopeT{Xmlns:"http://schemas.xmlsoap.org/soap/envelope/"}

   w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))

   attr = make([]xml.Attr,0,len(Hand.ws.Xmlns) + 2)
   for xmlns, alias = range Hand.ws.Xmlns {
      attr = append(attr,xml.Attr{
         Name: xml.Name{Local: "xmlns:" + alias},
         Value: xmlns,
      })
   }

   attr = append(attr,xml.Attr{
      Name: xml.Name{Local: "xmlns"},
      Value: Hand.ws.TargetNamespace,
   })

   attr = append(attr,xml.Attr{
      Name: xml.Name{Local: "TargetNamespace"},
      Value: Hand.ws.TargetNamespace,
   })

   //ns = Hand.ws.Xmlns[Hand.ws.TargetNamespace]
   vv = v.(SoapData)
   vvv = vv.SetName(nm, attr)
   envelope.Body = SoapBodyT{
      Data: vvv,
   }

   return xml.NewEncoder(Writer{buf: w}).Encode(envelope)
}
