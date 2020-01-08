#!/usr/local/bin/qjs -m
import "./katex.js"
import * as std from "std";

function assert(cond) {
    if (!cond) throw new std.Error();
}

function detex_block(block) {
    if (Array.isArray(block)) {
        block.forEach(detex_block);
    } else if (block.t === 'Math') {
        assert(block.c.length === 2);
        assert((block.c[0].t === 'InlineMath') || (block.c[0].t === 'DisplayMath'));
        assert(typeof(block.c[1]) === 'string');
        block.t = 'RawInline';
        block.c = ['html', globalThis.render(block.c[1], block.c[0].t === 'DisplayMath', false)];
    } else if (Array.isArray(block.c)) {
        block.c.forEach(detex_block);
    }
}

function detex(doc) {
    doc.blocks.forEach(detex_block);
}

let json = JSON.parse(std.in.readAsString());

detex(json);

let result = JSON.stringify(json);

std.out.puts(result);
