package paperfishGo

import (
	"fmt"
	"strings"
)

func (ws *WSClientT) SetEndPoint(endpoint string, targets ...interface{}) error {
   var opName string
   var oldEndpoint string
   var endpointParts []string

   oldEndpoint = fmt.Sprintf("%s/%s",ws.Host,ws.BasePath)
   if oldEndpoint[0] == '/' {
      oldEndpoint = oldEndpoint[1:]
   }
   if oldEndpoint[len(oldEndpoint)-1] == '/' {
      oldEndpoint = oldEndpoint[:len(oldEndpoint)-1]
   }

   if endpoint[0] == '/' {
      endpoint = endpoint[1:]
   }
   if endpoint[len(endpoint)-1] == '/' {
      endpoint = endpoint[:len(endpoint)-1]
   }

   endpointParts = strings.Split(endpoint,"/",)

   if len(targets) == 0 {
      ws.Host = strings.Join(endpointParts[:1],"/")
      ws.BasePath = strings.Join(endpointParts[1:],"/")
      for opName, _ = range ws.GetOperation {
         ws.GetOperation[opName].Path = strings.Replace(ws.GetOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

      for opName, _ = range ws.PostOperation {
         ws.PostOperation[opName].Path = strings.Replace(ws.PostOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

      for opName, _ = range ws.PutOperation {
         ws.PutOperation[opName].Path = strings.Replace(ws.PutOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

      for opName, _ = range ws.DeleteOperation {
         ws.DeleteOperation[opName].Path = strings.Replace(ws.DeleteOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

      for opName, _ = range ws.OptionsOperation {
         ws.OptionsOperation[opName].Path = strings.Replace(ws.OptionsOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

      for opName, _ = range ws.HeadOperation {
         ws.HeadOperation[opName].Path = strings.Replace(ws.HeadOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

      for opName, _ = range ws.PatchOperation {
         ws.PatchOperation[opName].Path = strings.Replace(ws.PatchOperation[opName].Path, oldEndpoint, endpoint, 1)
      }

   } else {
      _ = endpoint
   }
   return nil
}

