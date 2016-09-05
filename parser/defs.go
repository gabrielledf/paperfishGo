package parser

import (
	"github.com/luisfurquim/goose"
)



var Parser goose.Alert

//TODO: review this data structures
//Whenever you add a new struct or new field to handle a xml schema tag, you must add the tag name in TagsT structure
type SimplexTag struct {
	name			string
	restriction		map[string]string
}

type ComplexTag struct {
	name			string
	sequence		string
	complexContent	map[string]string
}

type TagsSchemaT struct {
	//Refers to attributeFormDefault, elementFormDefault, targetNamespace and xmlns attr from schema tag
	Attributes		map[string]string
	Import			map[string]string
	Annotation		string
	Documentation	string
	Elements		map[string]string
	SimplexType		SimplexTag
	ComplexType		ComplexTag		
}

var Tags TagsSchemaT

//Structure to store untreated schema tags
type untreatedTagT struct {
	Father		string
	Sequence	int
	Tag			string
}

var untreatedTags []untreatedTagT
