package marc21

import (
	"testing"
)

var (
	sfOne   = []byte{0x1F, 0x61, 0x47, 0x4F, 0x55, 0x4E, 0x42}
	sfTwo   = []byte{0x1F, 0x61, 0x47, 0x4F, 0x55, 0x4E, 0x42, 0x1F}
	sfThree = []byte{0x00, 0x1F, 0x61, 0x47, 0x4F, 0x55, 0x4E, 0x42, 0x1F}
)

func isCompare(one, two []byte) bool {
	if len(one) != len(two) {
		return false
	}
	for index := 0; index < len(one); index++ {
		if one[index] != two[index] {
			return false
		}
	}
	return true
}

func TestWriteSubfield(t *testing.T) {
	sf := new(SubField)
	sf.Name = "a"
	sf.Data = []byte("GOUNB")
	bin := sf.write()
	if !isCompare(bin, sfOne) {
		t.Errorf("Write subfield error. Wait % X, got % X", sfOne, bin)
	}
}

func TestReadSubfields(t *testing.T) {
	sf := readSubfields(sfOne)
	if len(sf) != 1 {
		t.Error("Must be one record, record count:", len(sf))
	}
	if sf[0].Name != "a" {
		t.Error("Subfield name must be \"a\", got:", sf[0].Name)
	}
	if string(sf[0].Data) != "GOUNB" {
		t.Error("Subfield data must be \"GOUNB\", got:", string(sf[0].Data))
	}
}

func TestReadSubfieldsTwo(t *testing.T) {
	sf := readSubfields(sfTwo)
	if len(sf) != 1 {
		t.Error("Must be one record, record count:", len(sf))
	}
	if sf[0].Name != "a" {
		t.Error("Subfield name must be \"a\", got:", sf[0].Name)
	}
	if string(sf[0].Data) != "GOUNB" {
		t.Error("Subfield data must be \"GOUNB\", got:", string(sf[0].Data))
	}
}

func TestReadSubfieldsThree(t *testing.T) {
	sf := readSubfields(sfThree)
	if len(sf) != 1 {
		t.Error("Must be one record, record count:", len(sf))
	}
	if sf[0].Name != "a" {
		t.Error("Subfield name must be \"a\", got:", sf[0].Name)
	}
	if string(sf[0].Data) != "GOUNB" {
		t.Error("Subfield data must be \"GOUNB\", got:", string(sf[0].Data))
	}
}
