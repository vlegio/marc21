package marc21

import (
	"testing"
)

func Test_decodeDecimal(t *testing.T) {
	first := decodeDecimal([]byte("12345"))
	if first != 12345 {
		t.Error("Error on decodeDecimal, wait 12345, got:", first)
	}
	second := decodeDecimal([]byte("12"))
	if second != 12 {
		t.Error("Error on decodeDecimal, wait 12, got:", second)
	}
	third := decodeDecimal([]byte("1:"))
	if third != 20 {
		t.Error("Error on decodeDecimal, wait 20, got:", third)
	}
}

func Test_encodeDecimal(t *testing.T) {
	first := string(encodeDecimal(5))
	if first != "00005" {
		t.Error("Error on encodeDecimal, wait \"00005\", got:", first)
	}
	second := string(encodeDecimal(12345))
	if second != "12345" {
		t.Error("Error on encodeDecimal,wait \"12345\", got:", second)
	}
}

func Test_encodeDecimal4(t *testing.T) {
	first := string(encodeDecimal4(5))
	if first != "0005" {
		t.Error("Error on encodeDecimal4, wait \"0005\", got:", first)
	}
	second := string(encodeDecimal4(1234))
	if second != "1234" {
		t.Error("Error on encodeDecimal4,wait \"1234\", got:", second)
	}
}
