package paperfishGo

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"github.com/luisfurquim/goose"
	"io"
	"net/http"
	"reflect"
)

//type Callback func(...interface{}) // status, <param, ...>

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
	Path       string
	Schemes    []string
	Encoder    Encoder
	Decoder    Decoder
	PathParm   []*ParameterT
	HeaderParm []*ParameterT
	QueryParm  []*ParameterT
	BodyParm   *ParameterT
	FormParm   []*ParameterT
}

type WSClientT struct {
	Host             string
	BasePath         string
	Schemes          []string
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
	receiver  chan []interface{}
	cli2srvch chan WSockRequest
	bindch    chan WSockRequest
}

type WSockRequest struct {
	SubOperation string
	Params       []interface{}
	Callback     reflect.Value
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

//Created only to casting the variables inside customized Unmarshal methods to go on with xml.enconding Unmarshal
type GoUnmarshal WSDLStruct
type GoUnmarshalSchema SchemaT
type GoUnmarshalMessage Message

//TODO: fazer algum processamento para preencher os vetores de parâmetros a partir do nome das mensagens
type Message struct {
	Name string `xml:"message,attr"`
	Namespace        map[string]string
	NamespaceReverse map[string]string
	//InputParameters    []ParametersT
	//OutputParameters   []ParametersT
}

type ExtensionBaseT struct {
	Base     string     `xml:"base,attr"`
	Sequence []ElementT `xml:"sequence>element"`
}

type ComplexTypeT struct {
	Name          string           `xml:"name,attr"`
	Sequence      []ElementT       `xml:"sequence>element"`
	ExtensionBase []ExtensionBaseT `xml:"complexContent>extension"`
}

type EnumerationT struct {
	Value string `xml:"value,attr"`
}

type RestrictionBaseT struct {
	Base        string         `xml:"base,attr"`
	Enumeration []EnumerationT `xml:"enumeration"`
}

type SimpleTypeT struct {
	Name            string           `xml:"name,attr"`
	RestrictionBase RestrictionBaseT `xml:"restriction"`
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
	AttributeFormDefault string         `xml:"attributeFormDefault,attr"`
	ElementFormDefault   string         `xml:"elementFormDefault,attr"`
	TargetNamespace      string         `xml:"targetNamespace,attr"`
	XMLAttr              []xml.Attr     `xml:",any"`
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
	TargetName    string      `xml:"targetNamespace,attr"`
	Namespace        map[string]string
	NamespaceReverse map[string]string
	Documentation string      `xml:"documentation"`
	Types         []SchemaT   `xml:"types>schema"`
	Message       []MessageT  `xml:"message"`
	PortType      []Operation `xml:"portType>operation"`
	Binding       []BindingT  `xml:"binding"`
	Service       []Endpoint  `xml:"service"`
}

type GooseG struct {
	New   goose.Alert
	Fetch goose.Alert
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
var ErrUnknownKind error = errors.New("Err unknown kind")
var ErrUnknownOperation error = errors.New("Err unknown operation")
var ErrUnknownMimeType error = errors.New("Err unknown mimetype")
var ErrWrite error = errors.New("Err writing stream")
var ErrBuffer error = errors.New("Err writing buffer")

var Fake bool // Se true, não vai acessar web service, vai usar arquivos XML locais

//Whenever you add a new struct or new field to handle a xml tag, you must add the tag name in TagsT slice
var TagsT []string
