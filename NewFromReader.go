package paperfishGo

import (
   "io"
   "fmt"
   "bytes"
   "reflect"
   "strings"
   "net/http"
   "crypto/tls"
   "crypto/x509"
   "encoding/xml"
   "encoding/json"
   "github.com/luisfurquim/stonelizard"
)

func NewFromReader(contract io.Reader, client *http.Client) ([]WSClientT, error) {
   var err error
   var n int
   var ct stonelizard.SwaggerT
   var wsdl WSDLStruct
   var service Endpoint
   var port PortT
   var binding BindingT
   var oper ConcreteOperationT
   var e ElementT
   var t SchemaT
   var c ComplexTypeT
   var s SimpleTypeT
   var i, j int
   var ws []WSClientT
   var basepath string
   var pathname string
   var pathdef stonelizard.SwaggerPathT
   var peekChar []byte
   var op *OperationT
   var swaggerParm stonelizard.SwaggerParameterT
   var paperParm *ParameterT
   var k reflect.Kind
   var method string
   var operation *stonelizard.SwaggerOperationT
   var coder interface{}
   var subop *SubOperationT
   var subOpSpec *stonelizard.SwaggerWSOperationT
   var param stonelizard.SwaggerParameterT
   var pos int
   var xsdType reflect.Type
   var xsdSymTab XsdSymTabT
   var schemes []string
   var operName string
   var mName string
   var mesgName string
   var elemName string
   var operIndex int
   var mesgIndex int
   var typ string
   var soapenc *SoapLiteralHnd

   ws = []WSClientT{}
   if client == nil {
      client = &http.Client{
         Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
               InsecureSkipVerify: true,
               RootCAs:            x509.NewCertPool(),
               Certificates:       []tls.Certificate{tls.Certificate{}},
            },
            DisableCompression: true,
         },
      }
   }

   peekChar = make([]byte, 1)
   n, err = contract.Read(peekChar)
   if err != nil || n != 1 {
      Goose.New.Logf(1, "Error peeking first byte of service contract: %d/%c/%s", n, peekChar[0], err)
      return nil, err
   }

   if peekChar[0] == '<' {
      // This is a WSDL contract

      err = xml.NewDecoder(io.MultiReader(bytes.NewReader(peekChar), contract)).Decode(&wsdl)
      if err != nil {
         Goose.New.Logf(1, "Error decoding wsdl (%s)", err)
         return nil, err
      }

      for i = 0; i < len(wsdl.Types); i++ {
         for j = 0; j < len(wsdl.Types[i].Elements); j++ {
            if wsdl.Types[i].Elements[j].Type == "" && len(wsdl.Types[i].Elements[j].ComplexTypes) == 1 {
               wsdl.Types[i].Elements[j].ComplexTypes[0].Name = wsdl.Types[i].Elements[j].Name
               wsdl.Types[i].ComplexTypes = append(
                  wsdl.Types[i].ComplexTypes,
                  wsdl.Types[i].Elements[j].ComplexTypes...)
               wsdl.Types[i].Elements = append(wsdl.Types[i].Elements[:j], wsdl.Types[i].Elements[j+1:]...)
               j--
//          } else {
//             fmt.Printf("%#v\n\n", wsdl.Types[i].Elements[j])
            }
         }
      }

      xsdSymTab = XsdSymTabT{
         "iface":        &XsdSymT{Type: typeOfiface},
         "string":       &XsdSymT{Type: typeOfstring},
         "boolean":      &XsdSymT{Type: typeOfboolean},
         "int":          &XsdSymT{Type: typeOfint},
         "integer":      &XsdSymT{Type: typeOfint},
         "decimal":      &XsdSymT{Type: typeOffloat},
         "float":        &XsdSymT{Type: typeOffloat},
         "double":       &XsdSymT{Type: typeOffloat},
         "duration":     &XsdSymT{Type: typeOfduration},
         "dateTime":     &XsdSymT{Type: typeOftime},
         "time":         &XsdSymT{Type: typeOftime},
         "date":         &XsdSymT{Type: typeOftime},
         "gYearMonth":   &XsdSymT{Type: typeOftime},
         "gYear":        &XsdSymT{Type: typeOftime},
         "gMonthDay":    &XsdSymT{Type: typeOftime},
         "gDay":         &XsdSymT{Type: typeOftime},
         "gMonth":       &XsdSymT{Type: typeOftime},
         "hexBinary":    &XsdSymT{Type: typeOfBinary},
         "base64Binary": &XsdSymT{Type: typeOfBinary},
      }
         //         case "anyURI":
         //         case "QName":
         //         case "NOTATION":

      // Before parsing the types we must list all of them because they may be
      // defined in any order, despite of their dependences

      if wsdl.Namespace == nil {
         wsdl.Namespace = map[string]string{}
      }
      if wsdl.NamespaceReverse == nil {
         wsdl.NamespaceReverse = map[string]string{}
      }
      for _, ns := range wsdl.XMLAttr {
         if strings.ToLower(ns.Name.Space) == "xmlns" {
            wsdl.Namespace[ns.Name.Local] = ns.Value
            wsdl.NamespaceReverse[ns.Value] = ns.Name.Local
         }
         if ns.Name.Space=="" && strings.ToLower(ns.Name.Local)=="xmlns" {
            wsdl.TargetNamespace = ns.Value
         }
      }

      Goose.New.Logf(0,"_tns: %#v", wsdl.TargetNamespace)
      Goose.New.Logf(0,"_ns: %#v", wsdl.Namespace)
      Goose.New.Logf(0,"_nsr: %#v", wsdl.NamespaceReverse)

      for i, t = range wsdl.Types {
         for j, s = range t.SimpleTypes {
            xsdSymTab[s.Name] = &XsdSymT{
               xsdref: &wsdl.Types[i].SimpleTypes[j],
               name: s.Name,
               ns: t.TargetNamespace,
            }
         }
         for j, c = range t.ComplexTypes {
            xsdSymTab[c.Name] = &XsdSymT{
               xsdref: &wsdl.Types[i].ComplexTypes[j],
               name: c.Name,
               ns: t.TargetNamespace,
            }
         }
         Goose.New.Logf(0,"ns")
         if wsdl.Types[i].Namespace == nil {
            wsdl.Types[i].Namespace = map[string]string{}
         }
         if wsdl.Types[i].NamespaceReverse == nil {
            wsdl.Types[i].NamespaceReverse = map[string]string{}
         }
         for _, ns := range t.XMLAttr {
            if strings.ToLower(ns.Name.Space) == "xmlns" {
               wsdl.Types[i].Namespace[ns.Name.Local] = ns.Value
               wsdl.Types[i].NamespaceReverse[ns.Value] = ns.Name.Local
            }
            if ns.Name.Space=="" && strings.ToLower(ns.Name.Local)=="xmlns" {
               wsdl.Types[i].TargetNamespace = ns.Value
            }
         }
      }

      Goose.New.Logf(0,"tns: %#v", wsdl.Types[i].TargetNamespace)
      Goose.New.Logf(0,"ns: %#v", wsdl.Types[i].Namespace)
      Goose.New.Logf(0,"nsr: %#v", wsdl.Types[i].NamespaceReverse)

      for _, t = range wsdl.Types {
         for _, s = range t.SimpleTypes {
            //            Goose.New.Logf(1, "Iter of %s", s.Name)
            Goose.New.Logf(4, "+++ simple %s", s.Name)
            xsdType, err = convXsdToGo(&s, xsdSymTab, s.Name)
            if err != nil {
               Goose.New.Logf(1, "Error converting XSD definition of %s to reflect.Value (%s)", s.Name, err)
               return nil, err
            }

            Goose.New.Logf(4, "--- simple %s", s.Name)
            //            fmt.Printf("%T\n\n",reflect.New(xsdType).Elem().Interface())
            xsdSymTab[s.Name].Type = xsdType
         }
         for _, c = range t.ComplexTypes {
            //            Goose.New.Logf(1, "Iter of %s", c.Name)
            Goose.New.Logf(4, "+++ complex %s", c.Name)
            xsdType, err = convXsdToGo(&c, xsdSymTab, c.Name)
            if err != nil {
               Goose.New.Logf(1, "Error converting XSD definition of %s to reflect.Value (%s)", c.Name, err)
               return nil, err
            }

            Goose.New.Logf(4, "--- complex %s", c.Name)
            //            fmt.Printf("%T\n\n",reflect.New(xsdType).Elem().Interface())
            xsdSymTab[c.Name].Type = xsdType
         }
      }

      for _, service = range wsdl.Service {
         for _, port = range service.Port {
            pos = strings.Index(port.Address.Location, ":")
            if pos == -1 {
               pos = 0
               schemes = []string{}
            } else {
               schemes = []string{port.Address.Location[:pos]}
               pos += 3
               if pos >= len(port.Address.Location) {
                  Goose.New.Logf(1, "Error on %s (%s)", port.Address.Location, ErrBadAddressLocationOfService)
                  return nil, err
               }
            }
            ws = append(ws, WSClientT{
               TargetNamespace:  wsdl.TargetNamespace,
               Client:           client,
               Host:             port.Address.Location[pos:],
               Schemes:          schemes,
               Binding:          port.Binding,
               GetOperation:     map[string]*OperationT{},
               PostOperation:    map[string]*OperationT{},
               PutOperation:     map[string]*OperationT{},
               DeleteOperation:  map[string]*OperationT{},
               OptionsOperation: map[string]*OperationT{},
               HeadOperation:    map[string]*OperationT{},
               PatchOperation:   map[string]*OperationT{},
            })
         }
      }

      for _, binding = range wsdl.Binding {
         for _, oper = range binding.ConcreteOperation {
            for i = 0; i < len(ws); i++ {
               if oper.Operation.Location != "" {
                  pos = strings.Index(oper.Operation.Location, ":")
                  if pos == -1 {
                     pos = 0
                     schemes = []string{}
                  } else {
                     schemes = []string{oper.Operation.Location[:pos]}
                     pos += 3
                     if pos >= len(port.Address.Location) {
                        Goose.New.Logf(1, "Error on %s (%s)", port.Address.Location, ErrBadAddressLocationOfService)
                        return nil, err
                     }
                  }
                  ws[i].PostOperation[oper.Name] = &OperationT{
                     Path:    oper.Operation.Location[pos:],
                     Schemes: schemes,
                  }
               } else {
                  ws[i].PostOperation[oper.Name] = &OperationT{
                     Path:    ws[i].Host,
                     Schemes: ws[i].Schemes,
                  }
               }

               //BodyParm   *ParameterT,
               for operIndex = 0; operIndex < len(wsdl.PortType); operIndex++ {
                  operName = bName(wsdl.PortType[operIndex].Name)
                  if operName == oper.Name {
                     mesgName = bName(wsdl.PortType[operIndex].Input.Name)
                     break
                  }
               }

               if operIndex == len(wsdl.PortType) {
                  Goose.New.Logf(1, "Error no operation %s found on portype", oper.Name)
                  return nil, err
               }

               for mesgIndex = 0; mesgIndex < len(wsdl.Message); mesgIndex++ {
                  mName = bName(wsdl.Message[mesgIndex].Name)
                  if mName == mesgName {
                     elemName = bName(wsdl.Message[mesgIndex].Part.Element)
                     break
                  }
               }

               if mesgIndex == len(wsdl.Message) {
                  Goose.New.Logf(1, "Error no message %s found on messages", mesgName)
                  return nil, err
               }

               Goose.New.Logf(1, "found %s on message part elements ", elemName)
               /*
                  if err != nil {
                     Goose.New.Logf(1, "Ignoring operation %s.%s.%s: %s", method, operation.OperationId, swaggerParm.Name, err)
                  }
               */

               Goose.New.Logf(1, "-----------------> %d - %s", len(wsdl.Types), elemName)
               for _, t = range wsdl.Types {
                  for _, e = range t.Elements {
                     if e.Name == elemName {
                        Goose.New.Logf(1, "t.ElementName: %s", e.Name)
                        typ = bName(e.Type)
                        Goose.New.Logf(1, "elemName: %s - type: %s - xsdSymTab[e.Type]: %#v", elemName, e.Type, xsdSymTab[typ])
                        ws[i].PostOperation[oper.Name].BodyParm = &ParameterT{
                           Name: wsdl.PortType[operIndex].Name,
                           Kind: xsdSymTab[typ].Type.Kind(),
                        }
                        xsdSymTab[e.Name] = &XsdSymT{
                           Type: xsdSymTab[typ].Type,
                           xsdref: xsdSymTab[typ].xsdref,
                        }
                        break
                     }
                  }
/*
                  for _, c = range t.ComplexTypes {
                     Goose.New.Logf(1, "c.Name: %#v", c.Name)
                     if c.Name == elemName {
                        ws[i].PostOperation[oper.Name].BodyParm = &ParameterT{
                           Name: wsdl.PortType[operIndex].Name,
                           Kind: reflect.Struct,
                        }
                        break
                     }
                  }
*/
               }

               if strings.ToLower(oper.Operation.Style) == "document" {
                  soapenc = &SoapLiteralHnd{
                     ws: &ws[i],
                     symtab: xsdSymTab,
                  }
                  if strings.ToLower(oper.InputSOAP.Use) == "literal" {
                     ws[i].PostOperation[oper.Name].Encoder = soapenc
                  }
                  if strings.ToLower(oper.OutputSOAP.Use) == "literal" {
                     ws[i].PostOperation[oper.Name].Decoder = soapenc
                  }
               }

               ws[i].symtab = xsdSymTab
            }
         }
      }


      return ws, nil
   } else {
      // This is a Swagger / OpenAPI contract
      err = json.NewDecoder(io.MultiReader(bytes.NewReader(peekChar), contract)).Decode(&ct)
      if err != nil {
         Goose.New.Logf(1, "Error decoding service contract: %s", err)
         return nil, err
      }

      ws = append(ws, WSClientT{
         Client:           client,
         GetOperation:     map[string]*OperationT{},
         PostOperation:    map[string]*OperationT{},
         PutOperation:     map[string]*OperationT{},
         DeleteOperation:  map[string]*OperationT{},
         OptionsOperation: map[string]*OperationT{},
         HeadOperation:    map[string]*OperationT{},
         PatchOperation:   map[string]*OperationT{},
      })

      ws[0].Host = ct.Host
      basepath = ct.BasePath
      if basepath[0] == '/' {
         basepath = basepath[1:]
      }
      if basepath[len(basepath)-1] == '/' {
         basepath = basepath[:len(basepath)-1]
      }
      ws[0].BasePath = basepath
      ws[0].Schemes = ct.Schemes

      // consumes
      coder, err = getCoder(ct.Consumes)
      if err != nil {
         Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
         return nil, err
      }

      ws[0].Encoder = coder.(Encoder)

      // Produces
      coder, err = getCoder(ct.Produces)
      if err != nil {
         Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
         return nil, err
      }

      ws[0].Decoder = coder.(Decoder)

      for pathname, pathdef = range ct.Paths {
      OperLoop:
         for method, operation = range pathdef {
            op = new(OperationT)

            if pathname[0] == '/' {
               pathname = pathname[1:]
            }
            op.Path = fmt.Sprintf("%s/%s/%s", ws[0].Host, ws[0].BasePath, pathname)

            // consumes
            coder, err = getCoder(operation.Consumes)
            if err != nil {
               Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
               return nil, err
            }
            op.Encoder = coder.(Encoder)

            // Produces
            coder, err = getCoder(operation.Produces)
            if err != nil {
               Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
               return nil, err
            }
            op.Decoder = coder.(Decoder)

            if operation.Schemes == nil {
               op.Schemes = ws[0].Schemes
            } else {
               op.Schemes = operation.Schemes
            }

            Goose.New.Logf(4, " %s %#v -> %s+%s -> %s %#v\n", operation.Consumes, op.Encoder, method, operation.OperationId, operation.Produces, op.Decoder)

            for _, subOpSpec = range operation.XWSOperations {
               Goose.New.Logf(2, "Registering sub-operation %s.%s.%s", method, operation.OperationId, subOpSpec.SuboperationId)
               subop = &SubOperationT{Id: subOpSpec.SuboperationId}
               for _, param = range subOpSpec.Parameters {
                  if param.Type != "" {
                     k, err = getKind(param.Type)
                  } else {
                     k = kindString
                  }
                  if err != nil {
                     Goose.New.Logf(1, "Ignoring sub-operation %s.%s.%s.%s: %s", method, operation.OperationId, subOpSpec.SuboperationId, param.Name, err)
                     continue
                  }

                  paperParm = &ParameterT{
                     Name: param.Name,
                     Kind: k,
                  }

                  subop.Parms = append(subop.Parms, paperParm)
               }
               if op.SubOperations == nil {
                  op.SubOperations = map[string]*SubOperationT{}
               }
               op.SubOperations[subOpSpec.SuboperationId] = subop
               Goose.New.Logf(5, "Registered sub-operation %s.%s.%s: %#v", method, operation.OperationId, subOpSpec.SuboperationId, op.SubOperations)
            }

            switch strings.ToLower(method) {
            case "get":
               ws[0].GetOperation[operation.OperationId] = op
            case "post":
               ws[0].PostOperation[operation.OperationId] = op
            case "put":
               ws[0].PutOperation[operation.OperationId] = op
            case "delete":
               ws[0].DeleteOperation[operation.OperationId] = op
            case "options":
               ws[0].OptionsOperation[operation.OperationId] = op
            case "head":
               ws[0].HeadOperation[operation.OperationId] = op
            case "patch":
               ws[0].PatchOperation[operation.OperationId] = op
            default:
               Goose.New.Logf(1, "Ignoring operation %s.%s: %s", method, operation.OperationId, ErrUnknownMethod)
               continue OperLoop
            }

            for _, swaggerParm = range operation.Parameters {
               if swaggerParm.Type != "" {
                  k, err = getKind(swaggerParm.Type)
               } else if (swaggerParm.Schema != nil) && (swaggerParm.Schema.Type != "") {
                  k, err = getKind(swaggerParm.Schema.Type)
               } else {
                  k = kindString
               }
               if err != nil {
                  Goose.New.Logf(1, "Ignoring operation %s.%s.%s: %s", method, operation.OperationId, swaggerParm.Name, err)
                  continue
               }

               paperParm = &ParameterT{
                  Name: swaggerParm.Name,
                  Kind: k,
               }

               switch strings.ToLower(swaggerParm.In) {
               case "path":
                  op.PathParm = append(op.PathParm, paperParm)
               case "header":
                  op.HeaderParm = append(op.HeaderParm, paperParm)
               case "query":
                  op.QueryParm = append(op.QueryParm, paperParm)
               case "body":
                  op.BodyParm = paperParm
               case "form":
                  op.FormParm = append(op.FormParm, paperParm)
               }
            }
         }
      }
   }

   return ws, nil
}
