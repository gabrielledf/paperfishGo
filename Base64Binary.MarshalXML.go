package paperfishGo

import (
   "encoding/xml"
   "encoding/base64"
)

func (bb Base64Binary) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
   e.EncodeElement(base64.StdEncoding.EncodeToString([]byte(bb)), start)
	return nil
}

