package parser

import (
	"io"
	"strings"
	"github.com/luisfurquim/goose"
	"github.com/PuerkitoBio/goquery"
)

//Makes the TagsT map
func Init() {
	Parser = goose.Alert(5)
	
    Tags.Attributes =   map[string]string{
    	"attributeformdefault"	:   "attributeFormDefault",
        "elementformdefault"	:	"elementFormDefault",
        "targetnamespace"		:	"targetNamespace",
        "xmlns"					:	"xmlns",
    }
	Tags.Import	=	map[string]string{
		"schemalocation"	:	"schemaLocation",
		"namespace"			:	"namespace",
	}
    Tags.Elements    =   map[string]string{
        "name"			:   "name",
        "type" 			:   "type",
        "nillable"		:   "nillable",
        "maxoccurs"		:   "maxOccurs",
        "minoccurs"		:   "minOccurs",
    }
    
    Tags.Annotation		=	"annotation>documentation"
    Tags.Documentation	=	"documentation"
   
}


//TODO1: fix the bug in paramater passing goquery.Find() method
//TODO2: Rethink the way the wsdl file using goquery.
func CheckWSDL(file io.Reader, ns string) {	
	Init()
	
	docQuery, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		Parser.Printf(1,"Error to read WSDL file. Error: %s\n", err)
	}
	Parser.Printf(9, "docQuery = %v\n", docQuery)
	
	if ns != "" {
		ns = ns + "\\\\:"
	}
	Parser.Printf(1,"ns = %s\n", ns+"schema")
	//docQuery.Find("schema").Each(func(i int, s *goquery.Selection) {
	//docQuery.Find("xsd\\:schema").Each(func(i int, s *goquery.Selection) {
	//docQuery.Find("xs\\:schema").Each(func(i int, s *goquery.Selection) {
	docQuery.Find(string(ns + "schema")).Each(func(i int, s *goquery.Selection) {
		Parser.Printf(9, "i = %d \n", i)
		for _, nodes := range s.Nodes {
            //Check each attribute of all schema tags
            for _, attr := range nodes.Attr {
				//To catch only string before ':' caracter
				keytest := strings.Split(attr.Key, ":")
				Parser.Printf(9,"valor de keytest = %s\n", keytest[0])
                if _, ok := Tags.Attributes[keytest[0]]; !ok {
					Parser.Printf(1, "Inconsistence: %s\n", attr.Key)
					untreatedTags = append(untreatedTags, untreatedTagT{goquery.NodeName(s), i, attr.Key })
				}
			}
            //Check each attribute of all import tags
            s.Find(ns + "import").Each(func(i int, selImp *goquery.Selection){
				Parser.Printf(1, "i = %d\n", i)

			})
         }
         
	})
	Parser.Printf(1,"Tags untreated: %s\n", untreatedTags)
}
