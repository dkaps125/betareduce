package betareduce

import (
	"fmt"
	"os"
)

func p_out(s string, args ...interface{}) {
	if !debug {
		return
	}
	p_err(s, args...)
}

func p_err(s string, args ...interface{}) {
	var pre string
	fmt.Printf(pre+s, args...)
}

func p_dieif(b bool, s string, args ...interface{}) {
	if b {
		p_err(s, args...)
		os.Exit(1)
	}
}

func p_die(s string, args ...interface{}) {
	p_err(s, args...)
	os.Exit(1)
}
