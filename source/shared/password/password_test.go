package password

import (
	"gowebapp/source/shared/testutil"
	"testing"
)

var res []byte
var resb bool

var randStr = testutil.RandAlphaStr

func BenchmarkHash8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, _ = Hash(randStr(8))
	}
}

func BenchmarkHash16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, _ = Hash(randStr(16))
	}
}

var (
	tkey    = "abcdefghijklmno"
	tval, _ = Hash(tkey)
)

func BenchmarkCompare16(b *testing.B) {
	var r bool
	for i := 0; i < b.N; i++ {
		r = Compare(randStr(16), tval)
	}
	resb = r
}
