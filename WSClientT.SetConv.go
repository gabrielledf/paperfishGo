package paperfishGo

func (ws *WSClientT) SetConv(fn func(string) string, targets ...interface{}) error {
   var opName string

   if len(targets) == 0 {
      for opName, _ = range ws.PostOperation {
         ws.PostOperation[opName].Decoder.(*SoapLiteralHnd).Conv = fn
      }
   } else {
      _ = fn
   }
   return nil
}

