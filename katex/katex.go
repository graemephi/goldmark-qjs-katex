// Package katex exposes a simplified API to KaTeX, run on QuickJS.
// Exported functions are thread-safe and can be called from any goroutine at any time.
//
// Use it like this if you're doing a lot of TeX rendering:
//    // var dest, src []byte
//    katex.Render(&dest, src, katex.Inline)
//
// Or like this, if you aren't:
//    // var dest io.Writer
//    // var src []byte
//    katex.RenderTo(dest, src, katex.Inline)
package katex

/*
#include "katex.h"
*/
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

// ErrBadInput indicates an internal KaTeX error, possibly due to differences
// between QuickJS and browser runtimes. These do not represent parse errors,
// which are rendered.
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

// Mode specifies how KaTeX is rendered with flags.
type Mode int

const (
	// Display renders in display mode if set, inline otherwise.
	Display Mode = 1 << iota

	// Warn prints KaTeX warnings to std out if set, suppresses them otherwise.
	Warn
)

// Named values for Mode flag combinations.
const (
	Inline      Mode = 0
	InlineWarn  Mode = Inline | Warn
	DisplayWarn Mode = Display | Warn
)

func (m Mode) String() string {
	switch m {
	case Inline:
		return "inline"
	case Display:
		return "display"
	case InlineWarn:
		return "inline|warn"
	case DisplayWarn:
		return "display|warn"
	}
	return "none"
}

// Warnings returns a Mode with the warning flag set or unset.
func Warnings(on bool) Mode {
	if on {
		return Warn
	}
	return Mode(0)
}

func render(dest []byte, src []byte, m C.Mode) ([]byte, error) {
	if len(src) == 0 {
		return dest[:0], nil
	}
	size := C.render(cref(dest), ccap(dest), cref(src), clen(src), m)
	if int(size) == -1 {
		return dest[:0], ErrBadInput
	}
	if size > ccap(dest) {
		dest = make([]byte, size)
		newSize := C.render(cref(dest), ccap(dest), cref(src), clen(src), m)
		if size != newSize {
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
func RenderTo(w io.Writer, src []byte, m Mode) error {
	size := len(src) * 150
	dest, err := render(make([]byte, size), src, C.Mode(m))
	if err == nil {
		w.Write(dest)
	}
	return err
}
