package betareduce

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func GetAddrComponents(address string) (string, int) {
	P_out("Getting comps for %s\n", address)
	toks := strings.Split(address, ":")
	addr := toks[0]
	port, _ := strconv.Atoi(toks[1])

	return addr, port
}

var semBootSend = make(chan (int), 1)
var semBootRecv = make(chan (int), 1)

func inBootSend() {
	semBootSend <- 1
}

func outBootSend() {
	<-semBootSend
}

func inBootRecv() {
	semBootRecv <- 1
}

func outBootRecv() {
	<-semBootRecv
}

var semClientSend = make(chan (int), 1)
var semClientRecv = make(chan (int), 1)

func inClientSend() {
	semClientSend <- 1
}

func outClientSend() {
	<-semClientSend
}

func inClientRecv() {
	semClientRecv <- 1
}

func outClientRecv() {
	<-semClientRecv
}

var semRepSend = make(chan (int), 1)
var semRepRecv = make(chan (int), 1)

func inRepSend() {
	semRepSend <- 1
}

func outRepSend() {
	<-semRepSend
}

func inRepRecv() {
	semRepRecv <- 1
}

func outRepRecv() {
	<-semRepRecv
}
