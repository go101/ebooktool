package nstd

import "fmt"

func Panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func Assert(b bool, msg interface{}) {
	if !b {
		panic(msg)
	}
}

func Assertf(b bool, format string, a ...interface{}) {
	if !b {
		Panicf(format, a...)
	}
}
