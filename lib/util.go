package betareduce

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

var Debug bool

func P_out(s string, args ...interface{}) {
	if !Debug {
		return
	}
	P_err(s, args...)
}

func P_err(s string, args ...interface{}) {
	var pre string
	fmt.Printf(pre+s, args...)
}

func P_dieif(b bool, s string, args ...interface{}) {
	if b {
		P_err(s, args...)
		os.Exit(1)
	}
}

func P_die(s string, args ...interface{}) {
	P_err(s, args...)
	os.Exit(1)
}

func GetValue(contents []byte) string {
	var ret string

	switch contents[0] {
	case 's':
		ret = string(contents[1:])
	case 'i':
		ret = strconv.Itoa(int(binary.BigEndian.Uint64(contents[1:])))
	default:
		ret = string(contents[1:])
	}

	return ret
}
