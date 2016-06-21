package marc21

type VariableField struct {
	Tag           string
	HasIndicators bool
	Indicators    []byte
	RawData       []byte
	Subfields     []*SubField
}

func readVariableField(rawField []byte, tag string) (vf *VariableField) {
	vf = new(VariableField)
	vf.Tag = tag
	vf.RawData = rawField
	vf.Subfields = readSubfields(vf.RawData)
	if len(vf.Subfields) > 0 {
		for i := 0; vf.RawData[i] != delimeter; i++ {
			vf.Indicators = append(vf.Indicators, vf.RawData[i])
		}
		vf.HasIndicators = len(vf.Indicators) > 0
	}
	return vf
}

func (vf *VariableField) write() (bin []byte) {
	if len(vf.Subfields) > 0 {
		if vf.HasIndicators {
			bin = vf.Indicators
		}
		for _, sf := range vf.Subfields {
			bin = append(bin, sf.write()...)
		}
		bin = append(bin, fieldTerminator)
	} else {
		bin = vf.RawData
	}
	return bin
}
