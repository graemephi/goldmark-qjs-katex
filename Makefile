katex/katex.bytecode.h: katex/katex.js
	qjsc -N qjsc_api -m -c -o $@ $<

gen_test.go: katex/filter.js katex/katex.bytecode.h
	go generate

test: gen_test.go katex/katex.bytecode.h
	go test

bench: gen_test.go katex/katex.bytecode.h
	go test -bench . -benchtime 10s
