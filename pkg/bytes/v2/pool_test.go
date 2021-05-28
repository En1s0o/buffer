package v2_test

import (
	"log"
	"strconv"
	"testing"

	v2 "github.com/En1s0o/buffer/pkg/bytes/v2"
)

func TestNewBuffer(t *testing.T) {
	for i := 0; i < 10; i++ {
		buffer := v2.NewBuffer()
		if buffer.Len() != 0 {
			t.Errorf("NewBuffer: buffer.Len() == %d", buffer.Len())
		}
		_, _ = buffer.Write([]byte("hello " + strconv.Itoa(i)))
		log.Printf("%s", string(buffer.Bytes()))
		log.Printf("%s", v2.Stats())
		buffer.DeRef()
	}
	log.Printf("%s", v2.Stats())
}
