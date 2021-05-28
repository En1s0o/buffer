package v1_test

import (
	"runtime"
	"testing"
	"time"
	_ "unsafe"

	v1 "github.com/En1s0o/buffer/pkg/bytes/v1"
)

func TestNewBuffer(t *testing.T) {
	{
		buf := v1.NewBuffer()
		_, _ = buf.WriteString("hello")
		str := buf.String()
		if str != "hello" {
			t.Errorf("NewBuffer: hello != %s", str)
		}
		_, _ = buf.WriteString("world")
	}

	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	{
		buf := v1.NewBuffer()
		if buf.Len() != 0 {
			t.Errorf("NewBuffer: buf.Len() == %d", buf.Len())
		}
	}
}
