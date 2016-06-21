package marc21

import (
	"strconv"
)

type MarcError string

func (me MarcError) Error() string {
	return string(me)
}

const (
	delimeter        = byte(0x1f)
	fieldTerminator  = byte(0x1e)
	recordTerminator = byte(0x1d)
	maxRecordSize    = 99999
)

var (
	errInvalidLength      = MarcError("Record length is invalid")
	errNoRecordTerminator = MarcError("Record must end in a RT")
	errInvalidLeader      = MarcError("Leader is Invalid")
)

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
