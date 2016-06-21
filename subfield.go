package marc21

const (
	nameLength    = 1
	startTmpIndex = 2
)

type SubField struct {
	Name string
	Data []byte
}

func readSubfields(raw []byte) (sf []*SubField) {
	for i := 0; i < len(raw); {
		if raw[i] == delimeter {
			i++
			if i == len(raw) {
				return sf
			}
			name := string(raw[i])
			i++
			start := i
			for raw[i] != delimeter && raw[i] != fieldTerminator && i < len(raw) && raw[i] != recordTerminator {
				i++
				if i == len(raw) {
					sf = append(sf, &SubField{name, raw[start:i]})
					return sf
				}
			}
			sf = append(sf, &SubField{name, raw[start:i]})
		} else {
			i++
		}
	}
	return sf
}

func (sf *SubField) write() (bin []byte) {
	bin = append([]byte{delimeter, sf.Name[0]}, sf.Data...)
	return bin
}
