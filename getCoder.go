package paperfishGo

func getCoder(mime []string) (interface{}, error) {
   var useJson, useXml, useText, useFormURL, useFormData, useBinary bool
   var s string

   for _, s = range mime {
      if s == "application/json" {
         useJson = true
      } else if s == "application/xml" {
         useXml = true
      } else if len(s) > 5 && s[:5] == "text/" {
         useText = true
      } else if s == "application/x-www-form-urlencoded" {
         useFormURL = true
      } else if s == "multipart/form-data" {
         useFormData = true
      } else if len(s) > 24 && s[:24] == "application/octet-stream" {
         useBinary = true
      }
   }

   if useJson {
      return JsonHnd{}, nil
   } else if useXml {
      // TODO
   } else if useFormURL {
      return FormURLHnd{}, nil
   } else if useFormData {
      return TextHnd{}, nil
   } else if useText {
      return TextHnd{}, nil
   } else if useBinary {
      if len(s) >= 31 && s[24:31] == ";base64" {
         return Base64Hnd{}, nil
      } else {
         return BinaryHnd{}, nil
      }
   }
   return nil, ErrUnknownMimeType
}
