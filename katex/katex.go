// Package katex exposes a simplified API to KaTeX, run on QuickJS.
// Exported functions are thread-safe and can be called from any goroutine at any time.
package katex

/*
#include "katex.h"
const size_t Error_BadInput = -1;
*/
import "C"

import (
	"unsafe"
	"errors"
	"io"
)

// ErrBadInput indicates a malformed utf-8 string or and internal KaTeX error.
// These do not represent parse errors, which are rendered.
var ErrBadInput = errors.New("bad KaTeX input")

// ErrInconsistent indicates that equivalent calls into the qjs returned
// different results. This almost certainly means the qjs runtime internal state
// has been corrupted. This has never been observed, so methods to detect this
// and recover are not implemented.
var ErrInconsistent = errors.New("inconsistent results between calls into qjs")

func clen(buf []byte) C.size_t {
	return C.size_t(len(buf))
}

func ccap(buf []byte) C.size_t {
	return C.size_t(cap(buf))
}

func cref(buf []byte) unsafe.Pointer {
	if cap(buf) == 0 {
		return unsafe.Pointer(nil)
	}
	return unsafe.Pointer(&buf[0])
}

// Mode specifies how KaTeX is rendered.
type Mode int

const (
	// Display renders in display mode if set, inline otherwise.
	Display Mode = 1 << iota

	// Warn prints KaTeX warnings to std out if set, suppresses them otherwise.
	Warn
)

const (
	Inline 		Mode = 0
	NoWarn	    Mode = 0
	InlineWarn  Mode = Inline|Warn
	DisplayWarn Mode = Display|Warn
)

func (m Mode) String() string {
	switch m & Display {
	case Inline:
		return "inline"
	case Display:
		return "display"
	}
	return "none"
}

func render(dest []byte, src []byte, m C.Mode) ([]byte, error) {
	if len(src) == 0 {
		return dest[:0], nil
	}
	size := C.render(cref(dest), ccap(dest), cref(src), clen(src), m)
	// Cast to avoid spurious error on clang.
	if C.ptrdiff_t(size) == C.ptrdiff_t(C.Error_BadInput) {
		return dest[:0], ErrBadInput
	}
	if size > ccap(dest) {
		dest = make([]byte, size)
		new_size := C.render(cref(dest), ccap(dest), cref(src), clen(src), m)
		if size != new_size {
			return dest[:0], ErrInconsistent
		}
	}
	return dest[:size], nil
}

// Render renders a TeX string to HTML with KaTeX. The intended use of this
// function is for callers to reuse dest to minimise allocations and copying.
func Render(dest *[]byte, src []byte, m Mode) error {
	var err error
	*dest, err = render(*dest, src, C.Mode(m))
	return err
}

// RenderTo renders a TeX string to HTML with KaTeX.
func RenderTo(w io.Writer, src []byte, m Mode, optionalInitialBufSize ...uintptr) error {
	size := uintptr(4096)
	if len(optionalInitialBufSize) > 0 {
		size = optionalInitialBufSize[0]
	}
	dest, err := render(make([]byte, size), src, C.Mode(m))
	if err == nil {
		w.Write(dest)
	}
	return err
}
