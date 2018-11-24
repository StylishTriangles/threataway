package view

import (
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New("template")
	}
}

func TestAddFlash(t *testing.T) {
	v := New("test")
	v.AddFlash("test", FlashNotice)
	v.AddFlash("test2", FlashSuccess)
	expected := []Flash{
		{Message: "test", Class: FlashNotice},
		{Message: "test2", Class: FlashSuccess},
	}
	returned := v.peekFlashes()
	if len(returned) != len(expected) {
		t.Error(
			"Expected length", len(expected),
			"got", len(returned))
	}
	for i := range expected {
		if expected[i].Class != returned[i].Class || expected[i].Message != returned[i].Message {
			t.Error(
				"Expected", expected[i],
				"got", returned[i])
		}
	}
}
