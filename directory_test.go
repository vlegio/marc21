package marc21

import (
	"testing"
)

var (
	dirOne = []byte{0x30, 0x31, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}
)

func TestWriteDirectory(t *testing.T) {
	vf := new(VariableField)
	vf.Tag = "010"
	vf.HasIndicators = true
	vf.Indicators = []byte{0x31, 0x23}
	vf.Subfields = append(vf.Subfields, &SubField{"a", []byte("GOUNB")})
	dirRaw, vfRaw := writeDirectory(0, vf)
	if !isCompare(vfRaw, vfOne) {
		t.Errorf("Write variable field error. Wait % X, got % X", vfOne, vfRaw)
	}
	if !isCompare(dirOne, dirRaw) {
		t.Errorf("Write directory error. Wait % X, got % X", dirOne, dirRaw)
	}
}

func TestReadDirectory(t *testing.T) {
	dirs := readDirectory(0, append(dirOne, fieldTerminator))
	if len(dirs) != 1 {
		t.Error("Wait one directory entry, got", len(dirs))
	}
	if dirs[0].tag != "010" {
		t.Error("Wait tag \"010\", got", dirs[0].tag)
	}
	if dirs[0].offset != 0 {
		t.Error("Wait offset 0, got", dirs[0].offset)
	}
	if dirs[0].length != len(vfOne) {
		t.Error("Wait length", len(vfOne), "got", dirs[0].length)
	}
}
