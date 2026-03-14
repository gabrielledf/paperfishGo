package paperfishGo

import (
   "net/http"
   "encoding/base64"
)

func (ws *WSClientT) signRequest(req *http.Request, method string, body []byte) {
   var err error
   var sig []byte
   var msgToSign string

   if ws.Pki == nil {
      return
   }

   msgToSign = method + "+" + req.URL.Path
   if len(body) > 0 {
      msgToSign += "\n" + string(body)
   }

   sig, err = ws.Pki.Sign(msgToSign)
   if err != nil {
      Goose.Fetch.Logf(1, "Error signing request: %s", err)
      return
   }

   req.Header.Set("X-Request-Signature", base64.StdEncoding.EncodeToString(sig))
   req.Header.Set("X-Request-Signer", ws.KeyId)
}
