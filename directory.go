package marc21

const (
	startBaseAddress      = 12
	endBaseAddress        = 17
	directoryRecordLength = 12
	tagLength             = 3
	offsetStart           = 7
	offsetEnd             = 12
	lengthStart           = 3
	lengthEnd             = 7
)

type directory struct {
	tag    string
	offset int
	length int
}

func readDirectory(baseAddress int, record []byte) (dir []*directory) {
	for i := 0; record[i] != fieldTerminator; i += directoryRecordLength {
		tag := string(record[i : i+tagLength])
		offset := decodeDecimal(record[i+offsetStart : i+offsetEnd])
		length := decodeDecimal(record[i+lengthStart : i+lengthEnd])
		dir = append(dir, &directory{tag, baseAddress + offset, length})
	}
	return dir
}

func writeDirectory(curOffset int, vf *VariableField) (dirRaw []byte, vfRaw []byte) {
	vfRaw = vf.write()
	dirRaw = []byte(vf.Tag)
	dirRaw = append(dirRaw, encodeDecimal4(len(vfRaw))...)
	dirRaw = append(dirRaw, encodeDecimal(curOffset)...)
	return dirRaw, vfRaw
}
