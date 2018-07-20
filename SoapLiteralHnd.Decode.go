package paperfishGo

import (
   "io"
   "encoding/xml"
)

func (Hand SoapLiteralHnd) Decode(r io.Reader, v interface{}) error {
   var vv SoapData
   var vvv interface{}
   var envelope soapEnvelopeT = soapEnvelopeT{Xmlns:"http://schemas.xmlsoap.org/soap/envelope/"}

   vv = v.(SoapData)
   vvv = vv.SetName(nm, nil)
   envelope.Body = SoapBodyT{
      Data: v,
   }

   return xml.NewDecoder(r).Decode(&envelope)
}
