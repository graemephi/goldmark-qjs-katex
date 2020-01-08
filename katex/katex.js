import katex from "./katex/katex.mjs"

// KaTeX expects console.warn function to exist and uses it even if you tell it not to
// (through the strict cfg parameter) for unicode inputs.
function warn(msg) {
    console.log("KaTeX warning: " + msg);
}
function noop() {}

function render(tex, displayMode, warnings) {
    console.warn = warnings ? warn : noop;
    return katex.renderToString(tex, {
        throwOnError: false,
        displayMode: displayMode,
    });
}

globalThis.render = render;
