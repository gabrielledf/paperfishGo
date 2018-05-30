package paperfishGo

/*

import (
	"bytes"
	"encoding/xml"
	"github.com/gabrielledf/paperfishGo/parser"
	"github.com/luisfurquim/goose"
	"io"
	"io/ioutil"
	"os"
	"strings"
	//"golang.org/x/net/html"
	//"fmt"
)

//TODO: Change goose.Printf to goose.Log where needed.
//TODO: Improve documentation
var schemaNS string

//Initializes the structs and variables to read and process the WSDL file in data structures
func New(name string, wsdlreader io.Reader) *WsdlT {
	var paperfish WsdlT

	XMLFile = wsdlreader
	debug.New = goose.Alert(5)
	debug.Read = goose.Alert(5)
	debug.Print = goose.Alert(5)

	paperfish = WsdlT{clientName: name}

	err := paperfish.readWSDL(wsdlreader)
	if err != nil {
		debug.New.Logf(1, "Failed init XML data struct from WSDL file %s", err)
		os.Exit(1)
	}

	//Input parameters are definitions, types, message, operations, bindings and service
	//paperfish.printWSDL(9,9,9,9,9,9)
	//paperfish.printWSDL(4,4,4,4,4,4)

	return &paperfish
}

//Treats the string contained in the xmlns attribute from schema tag to capture the key-value pair.
//For example, the atrribute xmlns:ax22="http://beans.tjrs/xsd", must be decoded in data strucutre XMLnsT follows:
// XMLns[alias] = value  -> XMLns[ax22] = "http: //beans.tjrs/xsd"
// Source: http://stackoverflow.com/questions/35044019/get-xml-namespace-prefix-in-go-using-unmarshal
func (ns *SchemaT) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Namespace = map[string]string{}
	ns.NamespaceReverse = map[string]string{}
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			ns.Namespace[attr.Name.Local] = attr.Value
			ns.NamespaceReverse[attr.Value] = attr.Name.Local
		}
	}

	// Go on with unmarshalling.
	decoder := (*GoUnmarshalSchema)(ns)
	return d.DecodeElement(decoder, &start)
}

//Treats the string contained in the xmlns attribute from definition tag to capture the key-value pair.
func (ns *WSDLStruct) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Namespace = map[string]string{}
	ns.NamespaceReverse = map[string]string{}
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			ns.Namespace[attr.Name.Local] = attr.Value
			ns.NamespaceReverse[attr.Value] = attr.Name.Local
		}
	}

	// Go on with unmarshalling.
	decoder := (*GoUnmarshal)(ns)
	return d.DecodeElement(decoder, &start)
}

//Treats the string contained in the xmlns attribute from operation>input/output tags to capture the key-value pair.
func (ns *Message) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Namespace = map[string]string{}
	ns.NamespaceReverse = map[string]string{}
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			ns.Namespace[attr.Name.Local] = attr.Value
			ns.NamespaceReverse[attr.Value] = attr.Name.Local
		}
	}

	// Go on with unmarshalling.
	decoder := (*GoUnmarshalMessage)(ns)
	return d.DecodeElement(decoder, &start)
}

//Reads WSDL file from local disk and unmarshal this file in WSDLStruct data structure
func (w *WsdlT) readWSDL(reader io.Reader) error {
	var wsdlbuf []byte
	var err error

	//Reads WSDL file
	wsdlbuf, err = ioutil.ReadAll(reader)
	if err != nil {
		debug.Read.Logf(1, "Failed reading WSDL file %s", err)
		return err
	}

	//Unmarshal WSDL file
	err = xml.Unmarshal(wsdlbuf, &w.xmlData)
	if err != nil {
		debug.Read.Logf(1, "Failed unmarshal xml: %v", err)
		return err
	}

	w.detectNS()
	debug.Read.Printf(9, "PaperfishGo: Schema namespace = %s\n", schemaNS)
	//Call the parser to check the schema tags of WSDL file.
	parser.CheckWSDL(bytes.NewReader(wsdlbuf), schemaNS)
	return nil
}

//Detect the schema tag namespace
func (w *WsdlT) detectNS() {
	for ns, value := range w.xmlData.Namespace {
		debug.Read.Printf(9, " [%s] = %s\n", ns, value)
		ok := strings.Contains(value, "XMLSchema")
		if ok {
			debug.Read.Printf(9, "SchemaNamespace = %v\n", ns)
			schemaNS = ns
			break

		}
	}
}

//Prints the XML structure in the command line. Only for debug purposes
func (w *WsdlT) printWSDL(definitions, types, message, operations, bindings, service int) {
	debug.Print.Printf(definitions, "-------------------------------------------------XML Data:--------------------------------------------------------------------------------\n\n")

	debug.Print.Printf(definitions, "Definitions \n")
	debug.Print.Printf(definitions, "TargetName: \t %s\n", w.xmlData.TargetName)
	debug.Print.Printf(definitions, "Namespaces: \t %s\n", w.xmlData.Namespace)
	debug.Print.Printf(definitions, "NamespaceReverse \t %s\n", w.xmlData.NamespaceReverse)
	debug.Print.Printf(definitions, "Documentation: \t %s\n", w.xmlData.Documentation)
	debug.Print.Printf(definitions, "---------------------------------------------------------------------------------------------\n")

	debug.Print.Printf(types, "Types Element \n")
	debug.Print.Printf(types, "\t Schemas \n")
	schemaNumber := 1
	for _, schema := range w.xmlData.Types {
		debug.Print.Printf(types, "\t Schema number %d\n", schemaNumber)

		debug.Print.Printf(types, "\t\t AttributeFormDefault: \t %s\n", schema.AttributeFormDefault)
		debug.Print.Printf(types, "\t\t ElementFormDefault: \t %s\n", schema.ElementFormDefault)
		debug.Print.Printf(types, "\t\t TargetNamespace: \t %s\n", schema.TargetNamespace)
		debug.Print.Printf(types, "\t\t XMLns: \t %v\n", schema.Namespace)
		debug.Print.Printf(types, "\t\t NamespaceReverse: \t %v\n", schema.NamespaceReverse)

		debug.Print.Printf(types, "\t \t Import: \t \n")
		for _, imp := range schema.Import {
			debug.Print.Printf(types, "\t\t\t Location: \t %s\n", imp.SchemaLocation)
			debug.Print.Printf(types, "\t\t\t NameSpace: \t %s\n", imp.NameSpace)
		}

		debug.Print.Printf(types, "\t\t Elements \n")
		for _, element := range schema.Elements {
			debug.Print.Printf(types, "\t\t\t Name: \t %s\n", element.Name)
			debug.Print.Printf(types, "\t\t\t Type: \t %s \n", element.Type)
			debug.Print.Printf(types, "\t\t\t Documentation: \t %s\n", element.Documentation)
			debug.Print.Printf(types, "\t\t\t Nillable: \t %s\n", element.Nillable)
			debug.Print.Printf(types, "\t\t\t MaxOccurs: \t %s\n", element.MaxOccurs)
			debug.Print.Printf(types, "\t\t\t MinOccurs: \t %s\n", element.MinOccurs)

			debug.Print.Printf(types, "\t\t\t Complex Types: \n")
			for _, complexTp := range element.ComplexTypes {
				debug.Print.Printf(types, "\t\t\t\t Name: \t %s\n", complexTp.Name)
				debug.Print.Printf(types, "\t\t\t\t Sequence:\n")
				for _, seq := range complexTp.Sequence {
					debug.Print.Printf(types, "\t\t\t\t\t Element Name: \t %s\n", seq.Name)
					debug.Print.Printf(types, "\t\t\t\t\t Element Type: \t %s\n", seq.Type)
					debug.Print.Printf(types, "\t\t\t\t\t Nillable: \t %s \n", seq.Nillable)
					debug.Print.Printf(types, "\t\t\t\t\t MaxOccurs: \t %s\n", seq.MaxOccurs)
					debug.Print.Printf(types, "\t\t\t\t\t MinOccurs: \t %s\n", seq.MinOccurs)
				}
			}
		}

		debug.Print.Printf(types, "\t\t Simple Types: \n")
		for _, simpleTp := range schema.SimpleTypes {
			debug.Print.Printf(types, "\t\t\t Name: \t %s\n", simpleTp.Name)
			debug.Print.Printf(types, "\t\t\t\t Base: \t %s\n", simpleTp.RestrictionBase.Base)
			debug.Print.Printf(types, "\t\t\t\t Enumeration:\n")
			for _, enum := range simpleTp.RestrictionBase.Enumeration {
				debug.Print.Printf(types, "\t\t\t\t\t Value: \t %s\n", enum.Value)
			}

		}

		debug.Print.Printf(types, "\t\t Complex Types: \n")
		for _, complexTp := range schema.ComplexTypes {
			debug.Print.Printf(types, "\t\t\t Name: \t %s\n", complexTp.Name)
			debug.Print.Printf(types, "\t\t\t Elements:\n")
			for _, seq := range complexTp.Sequence {
				debug.Print.Printf(types, "\t\t\t\t Name = %s", seq.Name)
				debug.Print.Printf(types, "\t Type =  %s", seq.Type)
				debug.Print.Printf(types, "\t Nillable =  %s", seq.Nillable)
				debug.Print.Printf(types, "\t MaxOccurs =  %s", seq.MaxOccurs)
				debug.Print.Printf(types, "\t MinOccurs =  %s\n", seq.MinOccurs)
			}
			debug.Print.Printf(types, "\t\t\t Extension Base \n")
			for _, extBase := range complexTp.ExtensionBase {
				debug.Print.Printf(types, "\t\t\t\t Base: %s\n", extBase.Base)
				for _, seq := range extBase.Sequence {
					debug.Print.Printf(types, "\t\t\t\t\t Name = %s", seq.Name)
					debug.Print.Printf(types, "\t Type =  %s", seq.Type)
					debug.Print.Printf(types, "\t Nillable =  %s", seq.Nillable)
					debug.Print.Printf(types, "\t MaxOccurs =  %s", seq.MaxOccurs)
					debug.Print.Printf(types, "\t MinOccurs =  %s\n", seq.MinOccurs)
				}
			}
		}

		schemaNumber++
		debug.Print.Printf(types, "*********************************************************\n")
	}

	debug.Print.Printf(types, "---------------------------------------------------------------------------------------------\n")

	debug.Print.Printf(message, "Message Element\n")
	for _, msg := range w.xmlData.Message {
		debug.Print.Printf(message, "\t Message Name: \t %s\n", msg.Name)
		debug.Print.Printf(message, "\t\t Element Part:\t %s\n", msg.Part.Element)
		debug.Print.Printf(message, "\t\t Element Name:\t %s\n", msg.Part.Name)
	}
	debug.Print.Printf(message, "---------------------------------------------------------------------------------------------\n")

	debug.Print.Printf(operations, "Operation Element\n")
	for _, oper := range w.xmlData.PortType {
		debug.Print.Printf(operations, "\t Name:\t %s\n", oper.Name)
		debug.Print.Printf(operations, "\t Input:\t %s\n", oper.Input)
		debug.Print.Printf(operations, "\t Output: %s\n", oper.Output)
		debug.Print.Printf(operations, "\n*******************************************************************************************************\n")

	}
	debug.Print.Printf(operations, "---------------------------------------------------------------------------------------------\n")

	debug.Print.Printf(bindings, "Binding Element \n")
	for _, bind := range w.xmlData.Binding {
		if bind.Protocol.Style != "" {
			debug.Print.Printf(bindings, "\t Name:\t\t %s\n", bind.Name)
			debug.Print.Printf(bindings, "\t Type:\t\t %s\n", bind.Type)
			debug.Print.Printf(bindings, "\t Soap Style:\t %s\n", bind.Protocol.Style)
			debug.Print.Printf(bindings, "\t Soap Transport: %s\n", bind.Protocol.Transport)
			debug.Print.Printf(bindings, "\t Concrete Operations:\n\n")
			for _, oper := range bind.ConcreteOperation {
				debug.Print.Printf(bindings, "\t\t Name:\t\t%s\n", oper.Name)
				debug.Print.Printf(bindings, "\t\t SoapAction:\t%s\n", oper.Operation.SoapAction)
				debug.Print.Printf(bindings, "\t\t SoapStyle:\t%s\n", oper.Operation.Style)
				debug.Print.Printf(bindings, "\t\t Input (use):\t %s\n", oper.InputSOAP.Use)
				debug.Print.Printf(bindings, "\t\t Output (use):\t %s\n\n", oper.OutputSOAP.Use)
			}
			debug.Print.Printf(bindings, "*******************************************************\n")
		} else {
			debug.Print.Printf(bindings, "\t Name:\t\t %s\n", bind.Name)
			debug.Print.Printf(bindings, "\t Type:\t\t %s\n", bind.Type)
			debug.Print.Printf(bindings, "\t Verb:\t\t %s\n", bind.Protocol.Verb)
			for _, oper := range bind.ConcreteOperation {
				debug.Print.Printf(bindings, "\t\t Name:\t\t%s\n", oper.Name)
				debug.Print.Printf(bindings, "\t\t Location:\t%s\n", oper.Operation.Location)
				debug.Print.Printf(bindings, "\t\t Input (use): %s\n", oper.InputHTTP.Type)
				debug.Print.Printf(bindings, "\t\t Output (use): %s\n\n", oper.OutputHTTP.Type)
			}
			debug.Print.Printf(bindings, "*******************************************************\n")
		}
	}
	debug.Print.Printf(bindings, "---------------------------------------------------------------------------------------------\n")
	debug.Print.Printf(service, "Service Element \n")
	for _, value := range w.xmlData.Service {
		debug.Print.Printf(service, "Service Name: %s\n", value.Name)
		for _, port := range value.Port {
			debug.Print.Printf(service, "\t Port Binding: %s\n", port.Binding)
			debug.Print.Printf(service, "\t Port Name: %s\n", port.Name)
			debug.Print.Printf(service, "\t Address Location: %s\n\n", port.Address)
		}
	}
	debug.Print.Printf(service, "---------------------------------------------------------------------------------------------\n")
}

*/

