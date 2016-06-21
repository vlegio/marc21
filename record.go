package marc21

import (
	"io"
)

type MarcRecord struct {
	Leader         *Leader
	directory      []*directory
	VariableFields []*VariableField
}

func ReadRecord(r io.Reader) (record *MarcRecord, err error) {
	lengthRaw := make([]byte, 5)

	_, err = r.Read(lengthRaw)
	if err != nil {
		return nil, err
	}

	length := decodeDecimal(lengthRaw)

	if length < leaderSize+2 || length > maxRecordSize {
		return nil, errInvalidLength
	}

	rawRecord := make([]byte, length)
	copy(rawRecord, lengthRaw)

	_, err = r.Read(rawRecord[5:])
	if err != nil {
		return nil, err
	}

	if rawRecord[len(rawRecord)-1] != recordTerminator {
		return nil, errNoRecordTerminator
	}

	record = new(MarcRecord)

	record.Leader = readLeader(rawRecord[0:leaderSize])

	record.directory = readDirectory(record.Leader.baseAddress, rawRecord[leaderSize:record.Leader.baseAddress])

	for _, dir := range record.directory {
		record.VariableFields = append(record.VariableFields, readVariableField(rawRecord[dir.offset:dir.offset+dir.length], dir.tag))
	}

	return record, nil
}

func (mr *MarcRecord) Write(w io.Writer) (err error) {
	var vfRaw []byte
	var dirRaw []byte
	for _, vf := range mr.VariableFields {
		curOffset := len(vfRaw)
		tmpDirRaw, tmpVfRaw := writeDirectory(curOffset, vf)
		dirRaw = append(dirRaw, tmpDirRaw...)
		vfRaw = append(vfRaw, tmpVfRaw...)
	}
	dirRaw = append(dirRaw, fieldTerminator)
	leaderRaw := mr.Leader.write(len(dirRaw))
	length := encodeDecimal(leaderSize + len(dirRaw) + len(vfRaw) + 1)
	_, err = w.Write(length)
	if err != nil {
		return err
	}
	_, err = w.Write(leaderRaw)
	if err != nil {
		return err
	}
	_, err = w.Write(dirRaw)
	if err != nil {
		return err
	}
	_, err = w.Write(vfRaw)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{recordTerminator})
	if err != nil {
		return err
	}
	return nil
}
