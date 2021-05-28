package v1

import (
	"bytes"
	"io"
)

// Buffer with automatic recycling mechanism.
type Buffer struct {
	inner *bytes.Buffer
}

func (b *Buffer) Bytes() []byte {
	return b.inner.Bytes()
}

func (b *Buffer) String() string {
	return b.inner.String()
}

func (b *Buffer) Len() int {
	return b.inner.Len()
}

func (b *Buffer) Cap() int {
	return b.inner.Cap()
}

func (b *Buffer) Truncate(n int) {
	b.inner.Truncate(n)
}

func (b *Buffer) Reset() {
	b.inner.Reset()
}

func (b *Buffer) Grow(n int) {
	b.inner.Grow(n)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	return b.inner.Write(p)
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	return b.inner.WriteString(s)
}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	return b.inner.ReadFrom(r)
}

// allow invoke multiple times.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	if nBytes := b.inner.Len(); nBytes > 0 {
		m, e := w.Write(b.Bytes())
		if m > nBytes {
			panic("Buffer.WriteTo: invalid Write count")
		}
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	return n, nil
}

func (b *Buffer) WriteByte(c byte) error {
	return b.inner.WriteByte(c)
}

func (b *Buffer) WriteRune(r rune) (n int, err error) {
	return b.inner.WriteRune(r)
}

// allow invoke multiple times.
func (b *Buffer) Read(p []byte) (n int, err error) {
	buf := b.Bytes()
	if len(buf) == 0 {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	n = copy(p, buf)
	return n, nil
}

func (b *Buffer) Next(n int) []byte {
	return b.inner.Next(n)
}
