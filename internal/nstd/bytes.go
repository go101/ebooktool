package nstd

import (
	"bytes"
)

func MergeByteSlices(bss ...[]byte) []byte {
	var n, allNils = 0, true
	for _, bs := range bss {
		n += len(bs)
		if bs != nil {
			allNils = false
		}
	}

	if n == 0 && allNils {
		return nil
	}

	var r = make([]byte, 0, n)
	for _, bs := range bss {
		r = append(r, bs...)
	}
	return r
}

type Bytes []byte

// The Bytes method convert a Bytes valule to []byte.
// It is not much necessary for most cases.
// Only useful in reflections.
func (bs Bytes) Bytes() []byte {
	return bs
}

func (bs Bytes) Index(sub []byte) int {
	return bytes.Index([]byte(bs), sub)
}

func (bs Bytes) Decap() Bytes {
	return bs[:len(bs):len(bs)]
}
