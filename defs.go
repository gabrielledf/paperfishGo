package paperfishGo

import (
	"io"
	"github.com/luisfurquim/goose"
)



type ParametersT struct {
	
	
}

//Created only to casting the variables inside customized Unmarshal methods to go on with xml.enconding Unmarshal
type GoUnmarshal			WSDLStruct
type GoUnmarshalSchema 		SchemaT
type GoUnmarshalMessage		Message

//TODO: na criação do cliente fazer algum processamento para preencher os vetores de parâmetros a partir do nome das mensagens
type Message struct {
	Name 				string				`xml:"message,attr"`
	Namespace			map[string]string
	NamespaceReverse	map[string]string
	//InputParameters		[]ParametersT
	//OutputParameters	[]ParametersT
}

type ExtensionBaseT struct {
	Base		string		`xml:"base,attr"`
	Sequence	[]ElementT	`xml:"sequence>element"`
}

type ComplexTypeT struct {
	Name			string				`xml:"name,attr"`
	Sequence		[]ElementT			`xml:"sequence>element"`
	ExtensionBase	[]ExtensionBaseT	`xml:"complexContent>extension"`
}

type EnumerationT struct {
	Value	string		`xml:"value,attr"`
}

type RestrictionBaseT struct {
	Base			string			`xml:"base,attr"`
	Enumeration		[]EnumerationT	`xml:"enumeration"`
}

type SimpleTypeT struct {
	Name				string				`xml:"name,attr"`
	RestrictionBase		RestrictionBaseT	`xml:"restriction"`
}


type ElementT struct {
	Name			string				`xml:"name,attr"`
	Type			string				`xml:"type,attr"`
	Documentation	string				`xml:"annotation>documentation"`
	Nillable		string				`xml:"nillable,attr"`
	MaxOccurs		string				`xml:"maxOccurs,attr"`
	MinOccurs		string				`xml:"minOccurs,attr"`
	ComplexTypes	[]ComplexTypeT		`xml:"complexType"`
}

type ImportT struct {
	SchemaLocation	string	`xml:"schemaLocation,attr"`
	NameSpace		string	`xml:"namespace,attr"`
}



type SchemaT struct {
	AttributeFormDefault	string				`xml:"attributeFormDefault,attr"`
	ElementFormDefault		string				`xml:"elementFormDefault,attr"`
	TargetNamespace			string				`xml:"targetNamespace,attr"`
	Namespace				map[string]string 	
	NamespaceReverse		map[string]string
	Import					[]ImportT			`xml:"import"`
	Elements				[]ElementT			`xml:"element"`
	SimpleTypes				[]SimpleTypeT		`xml:"simpleType"`
	ComplexTypes			[]ComplexTypeT		`xml:"complexType"`
}

type Operation struct {
	Name	string		`xml:"name,attr"`
	Input	Message		`xml:"input"`
	Output	Message		`xml:"output"`
}


type ProtocolBinding struct {
	Style		string	`xml:"style,attr"`
	Transport	string	`xml:"transport,attr"`
	
	Verb		string	`xml:"verb,attr"`
}

type HTTPContent struct {
	Type	string		`xml:"type,attr"`
}

type SoapBody struct {
	Use		string	`xml:"use,attr"`
	Type	string	`xml:"type,attr"`
}

type OperationBinding struct {
	SoapAction	string	`xml:"soapAction,attr"`
	Style		string	`xml:"style,attr"`
	
	Location	string	`xml:"location,attr"`
}

type ConcreteOperationT struct {
	Name			string				`xml:"name,attr"`
	Operation	 	OperationBinding	`xml:"operation"`
	
	InputSOAP		SoapBody 			`xml:"input>body"`
	OutputSOAP		SoapBody			`xml:"output>body"`
	
	InputHTTP		HTTPContent			`xml:"input>content"`
	OutputHTTP		HTTPContent			`xml:"output>content"`
}



type BindingT struct {
	Name					string						`xml:"name,attr"`
	Type					string						`xml:"type,attr"`
	Protocol				ProtocolBinding				`xml:"binding"`
	ConcreteOperation		[]ConcreteOperationT		`xml:"operation"`	
}


type Address struct {
	Location	string	`xml:"location,attr"`
}


type PortT struct {
	Binding		string	`xml:"binding,attr"`
	Name		string	`xml:"name,attr"`
	Address		Address	`xml:"address"`
}


type Endpoint struct {
	Name	string		`xml:"name,attr"`
	Port	[]PortT		`xml:"port"`
}

type PartT struct {
	Element string	`xml:"element,attr"`
	Name	string	`xml:"name,attr"`	
}

type MessageT struct {
	Name		string	`xml:"name,attr"`
	Part		PartT	`xml:"part"`
}


type WSDLStruct struct {
	TargetName			string				`xml:"targetNamespace,attr"`
	Namespace			map[string]string
	NamespaceReverse	map[string]string
	Documentation		string				`xml:"documentation"`
	Types				[]SchemaT			`xml:"types>schema"`
	Message				[]MessageT			`xml:"message"`
	PortType			[]Operation			`xml:"portType>operation"`
	Binding				[]BindingT			`xml:"binding"`
	Service				[]Endpoint			`xml:"service"`
}


type Goose struct {
	New			goose.Alert
	Read		goose.Alert
	Print		goose.Alert
}

var ClientName	string
var XMLFile		io.Reader
var xmlData		WSDLStruct
var debug		Goose

