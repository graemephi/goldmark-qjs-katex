package katex_test

import (
	"bytes"
	"testing"

	"github.com/graemephi/goldmark-qjs-katex/katex"
)

func TestUnicode(t *testing.T) {
	in := []byte("🐢")
	want := []byte("<span class=\"katex-display\"><span class=\"katex\"><span class=\"katex-mathml\"><math xmlns=\"http://www.w3.org/1998/Math/MathML\"><semantics><mrow><mtext>🐢</mtext></mrow><annotation encoding=\"application/x-tex\">🐢</annotation></semantics></math></span><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0em;vertical-align:0em;\"></span><span class=\"mord\">🐢</span></span></span></span></span>")
	got := make([]byte, len(want))
	err := katex.Render(&got, in, katex.Display)
	if err != nil {
		t.Errorf("Failed to convert %s: %s", in, err)
		return
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got, want:\n. %s\n. %s", got, want)
	}
}

func TestNonUTF8(t *testing.T) {
	in := []byte("\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff")
	dest := []byte{}
	err := katex.Render(&dest, in, katex.Display)
	if err != nil {
		t.Errorf("Failed to convert [70-ff]: %s", err)
		return
	}
	open := []byte("<span class=\"katex-display\">")
	close := []byte("</span>")
	ok := !bytes.Equal(dest[:len(open)], open)
	ok = ok || !bytes.Equal(dest[len(dest)-len(close):], close)
	ok = ok || (bytes.Count(dest, []byte("<span")) != bytes.Count(dest, []byte("</span")))
	if ok {
		t.Errorf("KaTeX did not render complete HTML: %s", dest)
	}
}

func TestRenderTo(t *testing.T) {
	in := []byte("\\sin(x)=x")
	want := []byte("<span class=\"katex\"><span class=\"katex-mathml\"><math xmlns=\"http://www.w3.org/1998/Math/MathML\"><semantics><mrow><mi>s</mi><mi>i</mi><mi>n</mi><mo stretchy=\"false\">(</mo><mi>x</mi><mo stretchy=\"false\">)</mo><mo>=</mo><mi>x</mi></mrow><annotation encoding=\"application/x-tex\">sin(x)=x</annotation></semantics></math></span><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:1em;vertical-align:-0.25em;\"></span><span class=\"mord mathdefault\">s</span><span class=\"mord mathdefault\">i</span><span class=\"mord mathdefault\">n</span><span class=\"mopen\">(</span><span class=\"mord mathdefault\">x</span><span class=\"mclose\">)</span><span class=\"mspace\" style=\"margin-right:0.2777777777777778em;\"></span><span class=\"mrel\">=</span><span class=\"mspace\" style=\"margin-right:0.2777777777777778em;\"></span></span><span class=\"base\"><span class=\"strut\" style=\"height:0.43056em;vertical-align:0em;\"></span><span class=\"mord mathdefault\">x</span></span></span></span>")
	a := make([]byte, len(want))
	b := bytes.Buffer{}
	err := katex.Render(&a, in, katex.Inline)
	if err != nil {
		t.Errorf("Failed to convert %s: %s", in, err)
		return
	}
	err = katex.RenderTo(&b, in, katex.Inline)
	if err != nil {
		t.Errorf("Failed to convert %s: %s", in, err)
		return
	}
	if !bytes.Equal(a, b.Bytes()) {
		t.Errorf("Render and RenderTo are inconsistent:\n. %s\n. %s", a, b.Bytes())
	}
}

func TestTooLarge(t *testing.T) {
	dest := []byte{}
	err := katex.Render(&dest, make([]byte, 1<<30), katex.Inline)
	if err != katex.ErrTooLarge {
		t.Errorf("accepted too large input")
	}
}
