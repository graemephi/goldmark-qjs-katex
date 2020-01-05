import katex from "./katex/katex.mjs";

let inline = {
    throwOnError: false,
    displayMode: false,
    output: 'html'
};

let display = {
    throwOnError: false,
    displayMode: true,
    output: 'html'
};

function renderInlineMaths(tex) {
    return katex.renderToString(tex, inline);
}

function renderDisplayMaths(tex) {
    return katex.renderToString(tex, display);
}

globalThis.renderInlineMaths = renderInlineMaths;
globalThis.renderDisplayMaths = renderDisplayMaths;
