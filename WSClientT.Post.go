package paperfishGo

import (
   "io"
//   "fmt"
   "bytes"
   "strings"
   "net/http"
   "io/ioutil"
//   "encoding/xml"
   "mime/multipart"
)

func (ws *WSClientT) Post(opName string, input map[string]interface{}, output interface{}) (int, error) {
   var err error
   var ok bool
   var op *OperationT
   var trackId string
   var targetURI string
   var sch, scheme string
   var lenBodyParm int
   var req *http.Request
   var resp *http.Response
   var postdata []byte
   var postparm []*ParameterT
   var buf []byte
//   var soapbuf []string
   var mr *multipart.Reader
   var boundary [][]byte
   var p *multipart.Part
   var slurp []byte
   var ctype string
   var hd string
//   var xopResp string

   if op, ok = ws.PostOperation[opName]; !ok {
      Goose.Fetch.Logf(1, "%s", ErrUnknownOperation)
      return 0, ErrUnknownOperation
   }

   if op.BodyParm != nil {
      lenBodyParm = 1
   }

   if (len(op.PathParm) + len(op.HeaderParm) + len(op.QueryParm) + lenBodyParm + len(op.FormParm)) != len(input) {
      Goose.Fetch.Logf(1, "%s for %s", ErrWrongParmCount, opName)
      return 0, ErrWrongParmCount
   }

   trackId, err = NewTrackId()
   if err != nil {
      Goose.Fetch.Logf(1, "Error generating new Track ID: %s", err)
      return 0, err
   }

   if op.Path[:4] != "http" {
      scheme = "http"

      for _, sch = range op.Schemes {
         if strings.ToLower(sch) == "https" {
            scheme += "s"
         }
      }
   }

   targetURI, err = ws.prepareURL(op, opName, scheme, op.PathParm, input)
   if err != nil {
      Goose.Fetch.Logf(1, "Error preparing URL for %s: %s", opName, err)
      return 0, err
   }

   targetURI, err = ws.prepareQuery(op, opName, scheme, op.QueryParm, input, targetURI)
   if err != nil {
      Goose.Fetch.Logf(1, "Error preparing URL query for %s: %s", opName, err)
      return 0, err
   }

   if op.BodyParm != nil {
      postparm = []*ParameterT{op.BodyParm}
   } else {
      postparm = op.FormParm
   }

   postdata, err = ws.prepareBody(op, opName, postparm, input)
   if err != nil {
      Goose.Fetch.Logf(1, "Error preparing body for %s: %s", opName, err)
      return 0, err
   }

   Goose.Fetch.Logf(5, "Prepared body for %s @ %s: %s", opName, targetURI, postdata)

   req, err = http.NewRequest("POST", targetURI, bytes.NewReader(postdata))
   if err != nil {
      return 0, err
   }

   Goose.Fetch.Logf(6, "TID:[%s] Request:%#v", trackId, req)

   req.Header.Set("Connection", "keep-alive")

   err = ws.prepareHeaders(op, opName, op.HeaderParm, input, req)
   if err != nil {
      Goose.Fetch.Logf(1, "Error preparing headers for %s: %s", opName, err)
      return 0, err
   }

   Goose.Fetch.Logf(5, "TID:[%s] request URL path %s", trackId, req.URL.Path)
   resp, err = ws.Client.Do(req)
   if err != nil {
      Goose.Fetch.Logf(6, "TID:[%s] Error fetching %s:%s", trackId, opName, err)
      return 0, err
   }
   defer resp.Body.Close()

   buf, err = ioutil.ReadAll(resp.Body)
   if err != nil {
      Goose.Fetch.Logf(6, "TID:[%s] Error fetching response for %s:%s", trackId, opName, err)
      return 0, err
   }

   Goose.Fetch.Logf(8, "RESP: %s ===================================================", bytes.Replace(buf[:500],[]byte("\r\n"),[]byte("//"),-1))

  // Goose.Fetch.Fatalf(6, "TID:[%s] Regex: %#v", trackId, xopxmlEnvelopRE.Match(buf[:500]))

   if xopxmlEnvelopRE.Match(buf[:500]) {

      boundary = reBoundary.FindSubmatch(buf)
      if len(boundary) == 0 {
         Goose.Fetch.Logf(6, "TID:[%s] Error fetching multipart boundary for %s", trackId, opName, ErrEmptyString)
         return 0, ErrEmptyString
      }

      Goose.Fetch.Logf(6, "TID:[%s] Boundary [%s]\n", trackId, boundary[0])

      mr = multipart.NewReader(bytes.NewReader(buf), string(boundary[0][2:]))
//      xopParts = map[string]string{}
      for {
         p, err = mr.NextPart()
         if err == io.EOF {
            Goose.Fetch.Logf(6, "TID:[%s] fetching part for %s: %s", trackId, opName, io.EOF)
            break
         }
         if err != nil {
            break
         }
         slurp, err = ioutil.ReadAll(p)
         if err != nil {
            break
         }

         ctype = p.Header.Get("Content-Type")

         if len(ctype)>=19 && ctype[:19]=="application/xop+xml" {
            buf = slurp
         } else {
            hd = p.Header.Get("Content-ID")
            xopParts.Store(hd[1:len(hd)-1], string(slurp))
         }
   //      fmt.Printf("Part %q: %q\n", p.Header.Get("Content-Type"), slurp)
      }

      //soapbuf = strings.Split(string(buf),"\n")
      //buf = []byte(soapbuf[len(soapbuf)-2])
      Goose.Fetch.Logf(8, "REBUF: %s", buf)
   }

   Goose.Fetch.Logf(5, "Using decoder: %#v", ws.PostOperation[opName].Decoder)
   err = ws.PostOperation[opName].Decoder.Decode(bytes.NewReader(buf), output)
   if err != nil && err != io.EOF {
      Goose.Fetch.Logf(6, "TID:[%s] Error decoding response for %s:%s", trackId, opName, err)
      return 0, err
   }

   return resp.StatusCode, nil
}
