package marc21

import (
	"testing"
)

var (
	lead = []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x30, 0x30, 0x30, 0x32, 0x34, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
)

func TestLeaderWrite(t *testing.T) {
	l := new(Leader)
	l.Status = 0x1
	l.Type = 0x1
	l.BibLevel = 0x1
	l.ControlType = 0x1
	l.CharacterEncoding = 0x1
	l.IndicatorCount = 0x1
	l.SubfieldCodeCount = 0x1
	l.EncodingLevel = 0x1
	l.CatalogingForm = 0x1
	l.MultipartLevel = 0x1
	l.LengthOFFieldPort = 0x1
	l.StartCharPos = 0x1
	l.LengthImplemenDefine = 0x1
	l.Undefine = 0x1
	res := l.write(0)
	if !isCompare(res, lead) {
		t.Errorf("Fail write leader. Must be % X, got % X", lead, res)
	}
}

func TestLeaderRead(t *testing.T) {
	leader := append([]byte{0x30, 0x30, 0x30, 0x32, 0x34}, lead...)
	l := readLeader(leader)
	if l.length != leaderSize {
		t.Error("l.length uncorrect, must be 24, got", l.length)
	}
	if l.Status != 0x1 {
		t.Error("l.Status uncorrect, must be 0x1, got", l.Status)
	}
	if l.Type != 0x1 {
		t.Error("l.Type uncorrect, must be 0x1, got", l.Type)
	}
	if l.BibLevel != 0x1 {
		t.Error("l.BibLevel uncorrect, must be 0x1, got", l.BibLevel)
	}
	if l.ControlType != 0x1 {
		t.Error("l.ControlType uncorrect, must be 0x1, got", l.ControlType)
	}
	if l.CharacterEncoding != 0x1 {
		t.Error("l.CharacterEncoding uncorrect, must be 0x1, got", l.CharacterEncoding)
	}
	if l.IndicatorCount != 0x1 {
		t.Error("l.IndicatorCount uncorrect, must be 0x1, got", l.IndicatorCount)
	}
	if l.SubfieldCodeCount != 0x1 {
		t.Error("l.SubfieldCodeCount uncorrect, must be 0x1, got", l.SubfieldCodeCount)
	}
	if l.baseAddress != 24 {
		t.Error("l.baseAddress uncorrect, must be 24, got", l.baseAddress)
	}
	if l.EncodingLevel != 0x1 {
		t.Error("l.EncodingLevel uncorrect, must be 0x1, got", l.EncodingLevel)
	}
	if l.CatalogingForm != 0x1 {
		t.Error("l.CatalogingForm uncorrect, must be 0x1, got", l.CatalogingForm)
	}
	if l.MultipartLevel != 0x1 {
		t.Error("l.MultipartLevel uncorrect, must be 0x1, got", l.MultipartLevel)
	}
	if l.LengthOFFieldPort != 0x1 {
		t.Error("l.LengthOFFieldPort uncorrect, must be 0x1, got", l.LengthOFFieldPort)
	}
	if l.StartCharPos != 0x1 {
		t.Error("l.StartCharPos uncorrect, must be 0x1, got", l.StartCharPos)
	}
	if l.LengthImplemenDefine != 0x1 {
		t.Error("l.LengthImplemenDefine uncorrect, must be 0x1, got", l.LengthImplemenDefine)
	}
	if l.Undefine != 0x1 {
		t.Error("l.Undefine uncorrect, must be 0x1, got", l.Undefine)
	}
}
