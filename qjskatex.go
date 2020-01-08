// package qjsKaTeX is an extension for goldmark (http://github.com/yuin/goldmark) to perform server-side KaTeX rendering.
package qjskatex

import (
	"github.com/graemephi/goldmark-qjs-katex/katex"

	gm "github.com/yuin/goldmark"
	gma "github.com/yuin/goldmark/ast"
	gmp "github.com/yuin/goldmark/parser"
	gmr "github.com/yuin/goldmark/renderer"
	gmt "github.com/yuin/goldmark/text"
	gmu "github.com/yuin/goldmark/util"
)

type node struct {
	gma.BaseInline

	mode katex.Mode
	pos  gmt.Segment

	// Renderer Optimisation: use a single buffer to render TeX into for the
	// entire run. But Goldmark only lets us set per-instance state, not
	// per-run, which is problematic if multiple gorouties are rendering using
	// the same goldmark instance. We have a parser Context, but no way to
	// access it from the gmr. So we put a pointer to that single buffer on
	// every node. Kinda gross, but it's up to 1.6x faster on TeX heavy pages.
	// Before (initialBufSize=4096): BenchmarkSequencesAndSeries-4          20         436353295 ns/op         6995413 B/op       2169 allocs/op
	// Before (initialBufSize=8192): BenchmarkSequencesAndSeries-4          20         356185385 ns/op         8565546 B/op       2083 allocs/op
	// After: 						 BenchmarkSequencesAndSeries-4          20         278764605 ns/op         3978978 B/op       1532 allocs/op
	buf *[]byte
}

var texNode = gma.NewNodeKind("TeX")

func (n *node) Kind() gma.NodeKind {
	return texNode
}

func (n *node) Dump(source []byte, level int) {
	gma.DumpHelper(n, source, level, map[string]string{
		"pos":  `"` + string(n.pos.Value(source)) + `"`,
		"mode": n.mode.String(),
	}, nil)
}

type parser struct{}

var ctxKey = gmp.NewContextKey()

func (p *parser) Trigger() []byte {
	return []byte{'$'}
}

func blank(s []byte) bool {
	result := true
	for c := 0; c < len(s); c++ {
		if !gmu.IsSpace(s[c]) {
			result = false
			break
		}
	}
	return result
}

func (p *parser) Parse(parent gma.Node, block gmt.Reader, pc gmp.Context) gma.Node {
	// Pandoc only parses TeX as inline. Both $ and $$ always behave like `, and
	// never like ```.  We give TeX the same parsing rules as `, except that inline
	// TeX must not have leading or trailing spaces, so we do not strip them.

	buf := block.Source()
	ln, pos := block.Position()
	lStart := pos.Start
	lEnd := pos.Stop
	line := buf[lStart:lEnd]

	if len(line) < 2 {
		return nil
	}

	start := 0
	end := 0
	advance := 0
	var mode katex.Mode

	if line[1] == '$' {
		// $$
		mode = katex.Display
		start = lStart + 2
		c_offset := 2

		for end == 0 {
			for c := c_offset; c < len(line); c++ {
				if line[c] == '$' {
					c++
					if c == len(line) {
						break
					}
					if line[c] == '$' {
						end = lStart + c - 1
						advance = 2
						break
					}
				}
			}
			if lEnd == len(buf) {
				// End of buffer, no closing $$
				break
			}
			if end == 0 {
				line = buf[lEnd:]
				c := 0
				for c < len(line) && line[c] != '\n' {
					c++
				}
				if blank(line[:c]) {
					// End of paragraph, no closing $$
					break
				}
				lStart = lEnd
				lEnd = lStart + c
				line = buf[lStart:lEnd]
				ln++
				c_offset = 0
			}
		}
	} else if !gmu.IsSpace(line[1]) {
		// $
		mode = katex.Inline
		start = lStart + 1

		for end == 0 {
			for c := 1; c < len(line); c++ {
				if line[c] == '\\' {
					c++
					continue
				}
				if line[c] == '$' {
					if !gmu.IsSpace(line[c-1]) || line[c-2] == '\\' {
						end = lStart + c
						advance = 1
						break
					}
				}
			}
			if lEnd == len(buf) {
				// End of buffer, no closing $
				break
			}
			if end == 0 {
				line = buf[lEnd:]
				c := 0
				for c < len(line) && line[c] != '\n' {
					c++
				}
				if blank(line[:c]) {
					// End of paragraph, no closing $
					break
				}
				lStart = lEnd
				lEnd = lStart + c
				line = buf[lStart:lEnd]
				ln++
			}
		}
	}

	if end-start <= 0 {
		return nil
	}

	// Consider parsing `[$ab$](c.td)` (1), `[a$b](c.td/$)` (2), `[$[]$](c.td)`
	// (3). We want to (1) to parse as as TeX-formatted link, and (2) to parse
	// as a link, because $ are valid in URLs (this is not the case for code
	// span backticks). But (3) can and should parse like (1). To detect and
	// disallow (2), but allow (1) and (3), we use the following rule:
	//
	// A TeX block inside a link is parsed as TeX iff a ] does not appear before
	// the first [ inside the TeX delimiters.
	//
	// Pandoc seems to enforce balanced brackets inside of link-enclosed TeX
	// blocks instead. This doesn't seem 'more correct' (or less) to me because
	// ordinary TeX can clearly contain unbalanced brackets. The goal here is to
	// prevent people from having to manually escape $ if they link something
	// like [sub$domain](foo.bar/sub$domain), for more complicated stuff they
	// can escape as usual.
	if pc.IsInLinkLabel() {
		src := buf[start:end]
		ok := true
		for c := 0; c < len(src); c++ {
			if src[c] == '\\' {
				c++
				continue
			}
			if src[c] == '[' {
				break
			}
			if src[c] == ']' {
				ok = false
				break
			}
		}
		if ok == false {
			return nil
		}
	}

	block.SetPosition(ln, gmt.NewSegment(end+advance, lEnd))

	var renderBuf *[]byte
	if v := pc.Get(ctxKey); v != nil {
		renderBuf = (v).(*[]byte)
	} else {
		renderBuf = new([]byte)
		*renderBuf = make([]byte, 4096)
		pc.Set(ctxKey, renderBuf)
	}

	return &node{
		mode: mode,
		pos:  gmt.NewSegment(start, end),
		buf:  renderBuf,
	}
}

type renderer struct {
	warn katex.Mode
}

func (r *renderer) render(w gmu.BufWriter, source []byte, gmnode gma.Node, entering bool) (gma.WalkStatus, error) {
	n := gmnode.(*node)
	tex := source[n.pos.Start:n.pos.Stop]
	err := katex.Render(n.buf, tex, n.mode|r.warn)
	if err == nil {
		w.Write(*n.buf)
	}
	return gma.WalkStop, err
}

func (r *renderer) RegisterFuncs(reg gmr.NodeRendererFuncRegisterer) {
	reg.Register(texNode, r.render)
}

type Extension struct {
	EnableWarnings bool
}

func (e *Extension) Extend(m gm.Markdown) {
	warn := katex.NoWarn
	if e.EnableWarnings {
		warn = katex.Warn
	}
	var p gmp.InlineParser = &parser{}
	var r gmr.NodeRenderer = &renderer{warn}
	m.Parser().AddOptions(gmp.WithInlineParsers(gmu.PrioritizedValue{ Value: p, Priority: 150}))
	m.Renderer().AddOptions(gmr.WithNodeRenderers(gmu.PrioritizedValue{ Value: r, Priority: 150}))
}
