package paperfishGo

import (
   "io"
   "time"
   "bytes"
   "errors"
   "regexp"
   "reflect"
   "net/http"
   "crypto/tls"
   "crypto/x509"
   "encoding/xml"
   "github.com/luisfurquim/goose"
)

type XsdSymT struct {
   reflect.Type
   name string
   ns string
   xsdref interface{}
}

type XsdSymTabT map[string]*XsdSymT

type ParameterT struct {
   Name string
   Kind reflect.Kind
}

type TextHnd struct {
   charset string
}
type JsonHnd struct{}
type XmlHnd struct{}
type FormURLHnd struct{}
type FormDataHnd struct{}
type BinaryHnd struct{}
type Base64Hnd struct{}
type SoapLiteralHnd struct {
   ws *WSClientT
   symtab XsdSymTabT
   Conv func(string) string
}

type SoapData interface {
   SetName(nm string, attr []xml.Attr) (interface{})
}

type SoapBodyT struct {
   XMLName xml.Name
   Data interface{}
}

/*
type soapEnvelopeT struct {
   XMLName xml.Name `xml:"SOAP-ENV:Envelope"`
   Xmlns string `xml:"xmlns:SOAP-ENV,attr"`
   Header string `xml:"SOAP-ENV:Header"`
   Body interface{} `xml:"SOAP-ENV:Body"`
}
*/

type soapEnvelopeT struct {
   XMLName xml.Name `xml:"soapenv:Envelope"`
   Xmlns string `xml:"xmlns:soapenv,attr"`
   Header string `xml:"soapenv:Header"`
   Body interface{} `xml:"soapenv:Body"`
}

type Writer struct {
   buf *bytes.Buffer
}

type Encoder interface {
   Encode(output *bytes.Buffer, name string, val interface{}, isTail bool) error
}

type Decoder interface {
   Decode(io.Reader, interface{}) error
}

type OperationT struct {
   Path          string
   Schemes       []string
   Encoder       Encoder
   Decoder       Decoder
   inMesg        string
   outMesg       string
   PathParm      []*ParameterT
   HeaderParm    []*ParameterT
   QueryParm     []*ParameterT
   BodyParm      *ParameterT
   FormParm      []*ParameterT
   SubOperations map[string]*SubOperationT
}

type SubOperationT struct {
   Id    string
   Parms []*ParameterT
}

type WSClientT struct {
   Host             string
   BasePath         string
   Binding          string
   Schemes          []string
   symtab           XsdSymTabT
   TargetNamespace  string
   Xmlns            map[string]string
   Encoder          Encoder
   Decoder          Decoder
   Client           *http.Client
   GetOperation     map[string]*OperationT
   PostOperation    map[string]*OperationT
   PutOperation     map[string]*OperationT
   DeleteOperation  map[string]*OperationT
   OptionsOperation map[string]*OperationT
   HeadOperation    map[string]*OperationT
   PatchOperation   map[string]*OperationT
}

type WSockClientT struct {
   SubOperations map[string]*SubOperationT
   receiver      chan []interface{}
   cli2srvch     chan WSockRequest
   bindch        chan WSockRequest
}

type CallbackT struct {
   Callback     reflect.Value
   FailCallback func(int)
}

type WSockRequest struct {
   SubOperation string
   Params       []interface{}
   CallbackT
}

/*
type WSockRequestParam struct {
   SubOperation string
   Request interface{}
}
*/

//type Operation = func(map[string]interface{}) (interface{},error)
type WSResponse map[string]interface{}

/*
type WsdlT struct {
   Endpoint string `xml:"definitions>service>port>address>location,attr"`
}
*/

//Created only to cast the variables inside customized XML Unmarshal methods
type GoUnmarshal WSDLStruct
type GoUnmarshalSchema SchemaT
type GoUnmarshalMessage Message

//TODO: fazer algum processamento para preencher os vetores de parâmetros a partir do nome das mensagens
type Message struct {
   Name             string `xml:"message,attr"`
   Namespace        map[string]string
   NamespaceReverse map[string]string
   //InputParameters    []ParametersT
   //OutputParameters   []ParametersT
}

type ExtensionBaseT struct {
   Base     string     `xml:"base,attr"`
   Sequence []ElementT `xml:"sequence>element"`
}

type AttributeT struct {
   Name string `xml:"name,attr"`
   Type string `xml:"type,attr"`
   Use  string `xml:"use,attr"`
}

type ComplexTypeT struct {
   Name          string           `xml:"name,attr"`
   Sequence      []ElementT       `xml:"sequence>element"`
   ExtensionBase []ExtensionBaseT `xml:"complexContent>extension"`
   Attribute     []AttributeT     `xml:"attribute"`
}

type EnumerationT struct {
   Value string `xml:"value,attr"`
}

type RestrictionBaseT struct {
   Base        string         `xml:"base,attr"`
   Enumeration []EnumerationT `xml:"enumeration"`
}

type ListT struct {
   ItemType string `xml:"itemType,attr"`
}

type SimpleTypeT struct {
   Name            string           `xml:"name,attr"`
   RestrictionBase RestrictionBaseT `xml:"restriction"`
   List            ListT            `xml:"list"`
}

//TODO: fazer algum processamento para a partir da string nillable (true ou false) para boolean go
type ElementT struct {
   Name          string         `xml:"name,attr"`
   Type          string         `xml:"type,attr"`
   Documentation string         `xml:"annotation>documentation"`
   Nillable      string         `xml:"nillable,attr"`
   MaxOccurs     string         `xml:"maxOccurs,attr"`
   MinOccurs     string         `xml:"minOccurs,attr"`
   ComplexTypes  []ComplexTypeT `xml:"complexType"`
}

type ImportT struct {
   SchemaLocation string `xml:"schemaLocation,attr"`
   NameSpace      string `xml:"namespace,attr"`
}

//TODO: fazer algum processamento para quebrar a string e pegar os valores corretos para o alias e para o value
type XMLnsT struct {
   Alias string
   Value string `xml:"xmlns,attr"`
}

type SchemaT struct {
   AttributeFormDefault string     `xml:"attributeFormDefault,attr"`
   ElementFormDefault   string     `xml:"elementFormDefault,attr"`
   TargetNamespace      string     `xml:"targetNamespace,attr"`
   XMLAttr              []xml.Attr `xml:",any,attr"`
   Namespace            map[string]string
   NamespaceReverse     map[string]string
   Import               []ImportT      `xml:"import"`
   Elements             []ElementT     `xml:"element"`
   SimpleTypes          []SimpleTypeT  `xml:"simpleType"`
   ComplexTypes         []ComplexTypeT `xml:"complexType"`
}

type Operation struct {
   Name   string  `xml:"name,attr"`
   Input  Message `xml:"input"`
   Output Message `xml:"output"`
}

type ProtocolBinding struct {
   Style     string `xml:"style,attr"`
   Transport string `xml:"transport,attr"`

   Verb string `xml:"verb,attr"`
}

type HTTPContent struct {
   Type string `xml:"type,attr"`
}

type SoapBody struct {
   Use  string `xml:"use,attr"`
   Type string `xml:"type,attr"`
}

type OperationBinding struct {
   SoapAction string `xml:"soapAction,attr"`
   Style      string `xml:"style,attr"`

   Location string `xml:"location,attr"`
}

type ConcreteOperationT struct {
   Name      string           `xml:"name,attr"`
   Operation OperationBinding `xml:"operation"`

   InputSOAP  SoapBody `xml:"input>body"`
   OutputSOAP SoapBody `xml:"output>body"`

   InputHTTP  HTTPContent `xml:"input>content"`
   OutputHTTP HTTPContent `xml:"output>content"`
}

type BindingT struct {
   Name              string               `xml:"name,attr"`
   Type              string               `xml:"type,attr"`
   Protocol          ProtocolBinding      `xml:"binding"`
   ConcreteOperation []ConcreteOperationT `xml:"operation"`
}

type Address struct {
   Location string `xml:"location,attr"`
}

type PortT struct {
   Binding string  `xml:"binding,attr"`
   Name    string  `xml:"name,attr"`
   Address Address `xml:"address"`
}

type Endpoint struct {
   Name string  `xml:"name,attr"`
   Port []PortT `xml:"port"`
}

type PartT struct {
   Element string `xml:"element,attr"`
   Name    string `xml:"name,attr"`
}

type MessageT struct {
   Name string `xml:"name,attr"`
   Part PartT  `xml:"part"`
}

type WSDLStruct struct {
   TargetNamespace  string      `xml:"targetNamespace,attr"`
   XMLAttr          []xml.Attr  `xml:",any,attr"`
   Namespace        map[string]string
   NamespaceReverse map[string]string
   Documentation    string      `xml:"documentation"`
   Message          []MessageT  `xml:"message"`
   PortType         []Operation `xml:"portType>operation"`
   Binding          []BindingT  `xml:"binding"`
   Service          []Endpoint  `xml:"service"`
   Types            []SchemaT   `xml:"types>schema"`
}

type GooseG struct {
   New   goose.Alert
   Fetch goose.Alert
   Set   goose.Alert
   Type  goose.Alert
}

/*
type PortT struct {
   Endpoint string `xml:"location,attr"`
}

type WsdlT struct {
//   Add PortT `xml:"service>port>address"`
   clientName string
   xmlData    WSDLStruct
}
*/


var iface interface{}
var typeOfiface reflect.Type = reflect.TypeOf(&iface).Elem()
var typeOfstring reflect.Type = reflect.TypeOf("")
var typeOfboolean reflect.Type = reflect.TypeOf(true)
var typeOfint reflect.Type = reflect.TypeOf(1)
var typeOffloat reflect.Type = reflect.TypeOf(1.0)
var typeOfduration reflect.Type = reflect.TypeOf(time.Second)
var typeOftime reflect.Type = reflect.TypeOf(time.Time{})
var typeOfBinary reflect.Type = reflect.TypeOf([]byte{})

type Base64Binary []byte


var RootCAs *x509.CertPool
var CliCerts []tls.Certificate

var Goose GooseG

var ErrEmptyResponse error = errors.New("Empty response")
var ErrParmNotFound error = errors.New("Parameter not found")
var ErrWrongParmCount error = errors.New("Wrong parameter count")
var ErrWrongParmType error = errors.New("Wrong parameter type")
var ErrWrongReturnParmType error = errors.New("Wrong return parameter type")
var ErrFetchingContract error = errors.New("Error fetching contract")
var ErrUnknownMethod error = errors.New("Error unknown method")
var ErrUnknownKind error = errors.New("Error unknown kind")
var ErrUnknownOperation error = errors.New("Error unknown operation")
var ErrNilHandle error = errors.New("Err nil handle")
var ErrProtocol error = errors.New("Err protocol syntax error")
var ErrServer error = errors.New("Err on server")
var ErrUnknownMimeType error = errors.New("Error unknown mimetype")
var ErrWrite error = errors.New("Error writing stream")
var ErrBuffer error = errors.New("Error writing buffer")
var ErrUndetectableType error = errors.New("Error undetectable type")
var ErrEmptyString error = errors.New("Error empty string")
var ErrBadAddressLocationOfService error = errors.New("Bad address location of service")
var ErrNoElementFoundOnMessage error = errors.New("no element found on message")
var Fake bool // Se true, não vai acessar web service, vai usar arquivos XML locais

//Whenever you add a new struct or new field to handle a xml tag, you must add the tag name in TagsT slice
var TagsT []string

var IndentPrefix string = "   "

var xsd2go map[string]string = map[string]string{
   "anyURI"        : "string",
   "base64Binary"  : "paperfishGo.Base64Binary",
   "boolean"       : "bool",
   "byte"          : "int8",
   "date"          : "time.Time",
   "dateTime"      : "time.Time",
   "decimal"       : "float64",
   "double"        : "float64",
   "duration"      : "string",
   "float"         : "float32",
   "gDay"          : "string",
   "gMonth"        : "string",
   "gMonthDay"     : "string",
   "gYear"         : "string",
   "gYearMonth"    : "string",
   "hexBinary"     : "string",
   "ID"            : "string",
   "int"           : "int",
   "integer"       : "int",
   "language"      : "string",
   "long"          : "int64",
   "Name"          : "string",
   "short"         : "int16",
   "string"        : "string",
   "time"          : "time.Time",
   "unsignedByte"  : "uint8",
   "unsignedShort" : "uint16",
   "unsignedInt"   : "uint",
   "unsignedLong"  : "uint64",
}

var xopxmlEnvelopRE *regexp.Regexp = regexp.MustCompile(`(?m:^Content-Type:)`)
