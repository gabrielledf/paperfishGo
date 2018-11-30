package paperfishGo

import (
   "io"
   "strings"
   "io/ioutil"
   "encoding/xml"
)

func (Hand SoapLiteralHnd) Decode(r io.Reader, v interface{}) error {
   var err error
   var buf []byte

   if Hand.Conv != nil {
      buf, err = ioutil.ReadAll(r)
      if err != nil {
         Goose.Fetch.Logf(1,"Error fetching server response: err", err)
         return err
      }
      Goose.Fetch.Logf(1,"converting %s",buf)
      Goose.Fetch.Logf(1,"converted to %s",Hand.Conv(string(buf)))
      Goose.Fetch.Logf(1,"unmarshaling to %#v",v)
      err = xml.NewDecoder(strings.NewReader(Hand.Conv(string(buf)))).Decode(v)
   } else {
      Goose.Fetch.Logf(1,"not converting")
      err = xml.NewDecoder(r).Decode(v)
   }

//   Goose.Fetch.Logf(1,"vvv=%#v",v)

   return err
}
