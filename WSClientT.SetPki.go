package paperfishGo

import (
   "fmt"
   "github.com/luisfurquim/paracentric"
)

func (ws *WSClientT) SetPki(pki *paracentric.PkiT) {
   ws.Pki = pki
   if pki != nil && pki.Cert != nil {
      ws.KeyId = fmt.Sprintf("%2X", pki.Cert.SubjectKeyId)
   }
}
