package paperfishGo

func listXsdTypes(ct interface{}, xsdSymTab XsdSymTabT) {
	var e ElementT
	var n string
	var ok bool

	switch typ := ct.(type) {
	case (*ComplexTypeT):
		for _, e = range typ.Sequence {
			if e.Type != "" {
				n = bName(e.Type)
				if _, ok = xsdSymTab[n]; !ok {
					xsdSymTab[n] = &XsdSymT{}
				}
			} else if len(e.ComplexTypes) == 1 {
				listXsdTypes(&e.ComplexTypes[0], xsdSymTab)
			}

		}
	case (*SimpleTypeT):
		if typ.List.ItemType != "" {
			n = bName(typ.List.ItemType)
		} else {
			n = bName(typ.RestrictionBase.Base)
		}
		if _, ok = xsdSymTab[n]; !ok {
			xsdSymTab[n] = &XsdSymT{}
		}
	}
}
