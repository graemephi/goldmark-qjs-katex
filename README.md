# goldmark-qjs-katex

This is an extension for [Goldmark](https://github.com/yuin/goldmark) that adds TeX rendering using [KaTeX](https://katex.org/). It embeds [QuickJS](https://bellard.org/quickjs/) and QuickJS-compiled KaTeX bytecode.

The parser follows pandoc's rules for TeX in markdown. Right now, `$` and `$$` are the only supported delimiters. Also, only KaTeX's default configuration is supported.

### Performance

```
goos: windows
goarch: amd64
pkg: github.com/graemephi/goldmark-qjs-katex
BenchmarkSequencesAndSeries/NoKaTeX-4              12210            981796 ns/op          381171 B/op       1497 allocs/op
BenchmarkSequencesAndSeries/NoCache-4                 10        1069105530 ns/op         4900488 B/op       1563 allocs/op
BenchmarkSequencesAndSeries/Cache-4                 5041           2329294 ns/op         4443003 B/op       1545 allocs/op
```

## Usage

```
import (
	"github.com/graemephi/goldmark-qjs-katex"
	
	"github.com/yuin/goldmark"
)
```
```
markdown := goldmark.New(
	goldmark.WithExtensions(&qjskatex.Extension{}),
)
```

Also, godoc.

As [explained in the KaTeX documentation](https://katex.org/docs/node#including-in-webpages), you will need to use the KaTeX stylesheet in the HTML page that is used to display the math (but not the JavaScript file). To do this, you can link to the [the minimized stylesheet hosted at jsDeliver](https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css). The KaTeX documentation provides [an example](https://katex.org/docs/browser#starter-template) (but note that you only need the stylesheet, not the script). The file `katex.min.css` is also provided in the `katex/katex` folder in this repository. The current version of `goldmark-qjs-katex` uses KaTeX version `v0.16.11`, so use this version to avoid issues (although using a version of the form `v0.16.*` should be safe as well).

## Building

If you just want to build, gcc must be installed, and all you need to do is

```
go build
```

However, if you modify `./katex/katex.js`, then you must recompile the JS source to bytecode with qjsc. qjsc can be compiled from source in `./katex/quickjs/`. In addition, if these modifications change the way TeX is rendered then you need to regenerate the test cases. This requires pandoc to be installed. Run

```
make
```
to do all that, or look in the Makefile to see how to do it.

## Dependencies

[Goldmark](https://github.com/yuin/goldmark), [KaTeX](https://katex.org/), [QuickJS](https://bellard.org/quickjs/).

## Licenses

### Goldmark
```
MIT License

Copyright (c) 2019 Yusuke Inuzuka
```
### KaTeX
```
The MIT License (MIT)

Copyright (c) 2013-2018 Khan Academy
```
### QuickJS
```
QuickJS is released under the MIT license.

Unless otherwise specified, the QuickJS sources are copyright Fabrice
Bellard and Charlie Gordon.
```
