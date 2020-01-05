package qjsk

/*
#include "api.h"
*/
import "C"
import "unsafe"

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

func renderMaths(dest []byte, src []byte, displayMode C.DisplayMode) []byte {
	if len(src) == 0 {
		return dest[:0]
	}
	size := C.render_maths(cref(dest), ccap(dest), cref(src), clen(src), displayMode)
	if size < 0 {
		const msg = " !!!! Bad KaTeX render !!!! "
		copy(dest[:], msg)
		return dest[:len(msg)]
	}
	if size > ccap(dest) {
		dest = make([]byte, size)
		new_size := C.render_maths(cref(dest), ccap(dest), cref(src), clen(src), displayMode)
		if size != new_size {
			panic("inconsistent results between calls into qjs")
		}
	}
	return dest[:size]
}

func RenderInlineMaths(dest *[]byte, src []byte) {
	*dest = renderMaths(*dest, src, C.DisplayMode_Inline)
}

func RenderDisplayMaths(dest *[]byte, src []byte) {
	*dest = renderMaths(*dest, src, C.DisplayMode_Display)
}
