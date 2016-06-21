package marc21

const (
	leaderSize = 24
)

type Leader struct {
	length               int
	Status               byte
	Type                 byte
	BibLevel             byte
	ControlType          byte
	CharacterEncoding    byte
	IndicatorCount       byte
	SubfieldCodeCount    byte
	baseAddress          int
	EncodingLevel        byte
	CatalogingForm       byte
	MultipartLevel       byte
	LengthOFFieldPort    byte
	StartCharPos         byte
	LengthImplemenDefine byte
	Undefine             byte
}

func readLeader(raw []byte) (l *Leader) {
	l = new(Leader)
	l.length = decodeDecimal(raw[0:5])
	l.Status = raw[5]
	l.Type = raw[6]
	l.BibLevel = raw[7]
	l.ControlType = raw[8]
	l.CharacterEncoding = raw[9]
	l.IndicatorCount = raw[10]
	l.SubfieldCodeCount = raw[11]
	l.baseAddress = decodeDecimal(raw[startBaseAddress:endBaseAddress])
	l.EncodingLevel = raw[17]
	l.CatalogingForm = raw[18]
	l.MultipartLevel = raw[19]
	l.LengthOFFieldPort = raw[20]
	l.StartCharPos = raw[21]
	l.LengthImplemenDefine = raw[22]
	l.Undefine = raw[23]
	return l
}

func (l *Leader) write(dirLen int) (raw []byte) {
	raw = []byte{l.Status, l.Type, l.BibLevel, l.ControlType, l.CharacterEncoding, l.IndicatorCount, l.SubfieldCodeCount}
	raw = append(raw, encodeDecimal(dirLen+leaderSize)...)
	raw = append(raw, []byte{l.EncodingLevel, l.CatalogingForm, l.MultipartLevel, l.LengthOFFieldPort, l.StartCharPos, l.LengthImplemenDefine, l.Undefine}...)
	return raw
}
