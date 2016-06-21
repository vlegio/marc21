package marc21

import (
	"testing"
)

var (
	vfOne = []byte{0x31, 0x23, 0x1F, 0x61, 0x47, 0x4F, 0x55, 0x4E, 0x42, 0x1E}
	vfTwo = []byte{0x47, 0x4F, 0x55, 0x4E, 0x42}
)

func TestWriteVariableField(t *testing.T) {
	vf := new(VariableField)
	vf.HasIndicators = true
	vf.Indicators = []byte{0x31, 0x23}
	vf.Subfields = append(vf.Subfields, &SubField{"a", []byte("GOUNB")})
	res := vf.write()
	if !isCompare(vfOne, res) {
		t.Errorf("Write variable field error. Wait % X, got % X", vfOne, res)
	}
}

func TestWriteVariableFieldTwo(t *testing.T) {
	vf := new(VariableField)
	vf.HasIndicators = false
	vf.RawData = []byte("GOUNB")
	res := vf.write()
	if !isCompare(vfTwo, res) {
		t.Errorf("Write variable field error. Wait % X, got % X", vfTwo, res)
	}
}

func TestReadVariableFieldOne(t *testing.T) {
	vf := readVariableField(vfOne, "010")
	if !vf.HasIndicators {
		t.Error("Has indicators must be true, got false")
	}
	if !isCompare([]byte{0x31, 0x23}, vf.Indicators) {
		t.Errorf("Wait indicators 31 23, got %X", vf.Indicators)
	}
	if len(vf.Subfields) != 1 {
		t.Error("Subfield count must be 1, got:", len(vf.Subfields))
	}
	if vf.Subfields[0].Name != "a" {
		t.Error("Subfield[0].Name must be \"a\", got", vf.Subfields[0].Name)
	}
	if !isCompare([]byte("GOUNB"), vf.Subfields[0].Data) {
		t.Errorf("SubField[0].Data must be % X, got % X", []byte("GOUNB"), vf.Subfields[0].Data)
	}
}
