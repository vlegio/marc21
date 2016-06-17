package marc21

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

const (
	delimeter        = 0x1f
	fieldTerminator  = 0x1e
	recordTerminator = 0x1d

	leaderSize    = 24
	maxRecordSize = 99999
)

var (
	errInvalidLength      = fmt.Errorf("Record length is invalid")
	errNoRecordTerminator = fmt.Errorf("Record must end in a RT")
	errInvalidLeader      = fmt.Errorf("Leader is Invalid")
)

type OStack struct {
	Name string
	Data []byte
	Next *OStack
	Prev *OStack
}

func (o *OStack) Add(name string, data []byte) {
	if o == nil {
		o = new(OStack)
		o.Name = name
		o.Data = data
		o.Next = nil
		o.Prev = nil
	} else {
		if o.Name < name {
			if o.Next != nil {
				o.Next.Add(name, data)
			} else {
				nO := new(OStack)
				nO.Name = name
				nO.Data = data
				o.Next = nO
				nO.Prev = o
			}
		} else {
			nO := new(OStack)
			nO.Name = name
			nO.Data = data
			if o.Prev != nil {
				cp := o.Prev
				cp.Next = nO
				nO.Prev = cp
				nO.Next = o
				o.Prev = nO
			} else {
				nO.Next = o
				o.Prev = nO
			}
		}
	}
}

func (o *OStack) First() {
	if o == nil {
		return
	}
	for o.Prev != nil {
		if o.Prev == nil {
			return
		}
		o = o.Prev
	}
}

func (o *OStack) Len() int {
	if o == nil {
		return 0
	}
	for o.Prev != nil {
		if o.Prev == nil {
			break
		}
		o = o.Prev
	}
	i := 0
	for i = 0; o.Next != nil; i++ {
		o = o.Next
	}
	return i
}

func (o *OStack) Get(name string) *OStack {
	fmt.Println(o)
	o.First()
	if o == nil {
		return o
	}
	for o.Next != nil {
		if o.Name == name {
			return o
		}
		o = o.Next
	}
	return nil
}

type location struct {
	offset int
	length int
}

type MarcReader struct {
	r      io.Reader
	offset uint64
}

func NewReader(r io.Reader) (m *MarcReader) {
	m = new(MarcReader)
	m.r = r
	m.offset = 0
	return m
}

func (m *MarcReader) ReadRecord() (r *MarcRecord, err error) {
	length, rawRecord, err := readRecord(m.r)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	r = NewMarcRecord(rawRecord, m.offset)
	m.offset += uint64(length)
	return r, nil
}

type MarcRecord struct {
	Offset               uint64
	Status               byte
	Type                 byte
	BibLevel             byte
	ControlType          byte
	CharacterEncoding    byte
	IndicatorCount       byte
	SubfieldCodeCount    byte
	EncodingLevel        byte
	CatalogingForm       byte
	MultipartLevel       byte
	LengthOFFieldPort    byte
	StartCharPos         byte
	LengthImplemenDefine byte
	Undefine             byte
	Directory            map[string][]location
	VariableField        map[string]*RawField
}

func NewEmptyMarcRecord() (m *MarcRecord) {
	m = new(MarcRecord)
	m.VariableField = make(map[string]*RawField)
	return m
}

func NewMarcRecord(rawRecord []byte, offset uint64) (m *MarcRecord) {
	m = new(MarcRecord)
	m.Offset = offset

	m.Status = rawRecord[5]
	m.Type = rawRecord[6]
	m.BibLevel = rawRecord[7]
	m.ControlType = rawRecord[8]
	m.CharacterEncoding = rawRecord[9]
	m.IndicatorCount = rawRecord[10]
	m.SubfieldCodeCount = rawRecord[11]
	m.EncodingLevel = rawRecord[17]
	m.CatalogingForm = rawRecord[18]
	m.MultipartLevel = rawRecord[19]
	m.LengthOFFieldPort = rawRecord[20]
	m.StartCharPos = rawRecord[21]
	m.LengthImplemenDefine = rawRecord[22]
	m.Undefine = rawRecord[23]
	m.Directory = decodeDirectory(rawRecord)
	m.VariableField = make(map[string]*RawField)
	for k, v := range m.Directory {
		raw := make([][]byte, len(v))
		subfields := new(OStack) // make(map[string][]byte)
		for i := range v {
			start := v[i].offset
			end := v[i].offset + v[i].length
			raw[i] = rawRecord[start:end]
			for j := range raw[i] {
				if raw[i][j] == delimeter {
					sf := string(raw[i][j+1])
					z := 2
					if raw[i][z] == delimeter {
					delim:
						z++
						if raw[i][z] == sf[0] {
							z++
							start := z
							for raw[i][z] != delimeter && raw[i][z] != fieldTerminator {
								z++
								if z == len(raw[i])-1 {
									break
								}
							}
							subfields.Add(sf, raw[i][start:z])
							continue
						}
						for {
							switch raw[i][z] {
							case delimeter:
								goto delim
							case fieldTerminator:
								subfields.Add(sf, []byte{})
							default:
								z++
							}
						}
					}

				}
			}
		}
		m.VariableField[k] = new(RawField)
		m.VariableField[k].Indicators = raw[0][0:1]
		m.VariableField[k].Tag = k
		m.VariableField[k].RawData = raw
		m.VariableField[k].Subfields = subfields
	}
	return m
}

type RawField struct {
	Tag        string
	Indicators []byte
	RawData    [][]byte
	Subfields  *OStack
}

func decodeDirectory(record []byte) (dir map[string][]location) {
	baseAddress := decodeDecimal(record[12:17])
	dir = make(map[string][]location)
	for i := 24; record[i] != fieldTerminator; i += 12 {
		tag := string(record[i : i+3])
		dir[tag] = append(dir[tag],
			location{baseAddress + decodeDecimal(record[i+7:i+12]), decodeDecimal(record[i+3 : i+7])})
	}
	return dir
}

func decodeDecimal(n []byte) (result int) {
	result, err := strconv.Atoi(string(n))
	if err != nil {
		result = 0
		for i := range n {
			result = (10 * result) + int(n[i]-'0')
		}
	}
	return result
}

func encodeDecimal(n int) (result []byte) {
	str := strconv.Itoa(n)
	for len(str) < 5 {
		str = "0" + str
	}
	return []byte(str)
}

func encodeDecimal4(n int) (result []byte) {
	str := strconv.Itoa(n)
	for len(str) < 4 {
		str = "0" + str
	}
	return []byte(str)
}

func readRecord(r io.Reader) (length int, record []byte, err error) {
	tmp := make([]byte, 5)

	_, err = r.Read(tmp)
	if err != nil {
		return 0, nil, err
	}

	length = decodeDecimal(tmp)
	if length < leaderSize+2 || length > maxRecordSize {
		return 0, nil, errInvalidLength
	}

	record = make([]byte, length)
	copy(record, tmp)
	_, err = r.Read(record[5:])
	if err != nil {
		return 0, nil, err
	}

	if record[len(record)-1] != recordTerminator {
		return 0, nil, errNoRecordTerminator
	}

	return length, record, nil
}

func (m *MarcRecord) Write(w io.Writer) (err error) {
	record := new(bytes.Buffer)
	binary.Write(record, binary.LittleEndian, m.Status)
	binary.Write(record, binary.LittleEndian, m.Type)
	binary.Write(record, binary.LittleEndian, m.BibLevel)
	binary.Write(record, binary.LittleEndian, m.ControlType)
	binary.Write(record, binary.LittleEndian, m.CharacterEncoding)
	binary.Write(record, binary.LittleEndian, m.IndicatorCount)
	binary.Write(record, binary.LittleEndian, m.SubfieldCodeCount)

	rawFileds := new(bytes.Buffer)
	directory := new(bytes.Buffer)
	for tag, field := range m.VariableField {
		start := rawFileds.Len()
		if field.Subfields.Len() > 0 {
			field.Subfields.First()
			binary.Write(rawFileds, binary.LittleEndian, field.Indicators)
			for subfield := field.Subfields; subfield != nil; subfield = subfield.Next {
				if subfield.Name == "" {
					continue
				}
				binary.Write(rawFileds, binary.LittleEndian, byte(0x1f))
				binary.Write(rawFileds, binary.LittleEndian, []byte(subfield.Name))
				binary.Write(rawFileds, binary.LittleEndian, subfield.Data)
			}
			binary.Write(rawFileds, binary.LittleEndian, byte(0x1e))
		} else {
			var rawData []byte
			for _, data := range field.RawData {
				rawData = append(rawData, data...)
			}
			if len(rawData) < 1 {
				continue
			}
			binary.Write(rawFileds, binary.LittleEndian, rawData)
		}
		stop := rawFileds.Len()
		binary.Write(directory, binary.LittleEndian, []byte(tag))
		binary.Write(directory, binary.LittleEndian, encodeDecimal4(stop-start))
		binary.Write(directory, binary.LittleEndian, encodeDecimal(start))
	}

	binary.Write(directory, binary.LittleEndian, byte(0x1e))

	binary.Write(record, binary.LittleEndian, encodeDecimal(24+directory.Len())) //24 - Leader size
	binary.Write(record, binary.LittleEndian, m.EncodingLevel)
	binary.Write(record, binary.LittleEndian, m.CatalogingForm)
	binary.Write(record, binary.LittleEndian, m.MultipartLevel)
	binary.Write(record, binary.LittleEndian, m.LengthOFFieldPort)
	binary.Write(record, binary.LittleEndian, m.StartCharPos)
	binary.Write(record, binary.LittleEndian, m.LengthImplemenDefine)
	binary.Write(record, binary.LittleEndian, m.Undefine)
	binary.Write(record, binary.LittleEndian, directory.Bytes())
	binary.Write(record, binary.LittleEndian, rawFileds.Bytes())
	result := append(encodeDecimal(record.Len()+6), record.Bytes()...)
	result = append(result, recordTerminator)
	w.Write(result)
	return nil
}
