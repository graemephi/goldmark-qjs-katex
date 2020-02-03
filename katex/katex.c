#include "katex.h"

#include <string.h>

#include "quickjs/quickjs-libc.h"
#include "quickjs/quickjs.h"

#include "katex.bytecode.h"

typedef struct State {
    JSRuntime *rt;
    JSContext *ctx;
    JSAtom render;
    JSValue global_obj;
    JSValue true_val;
    JSValue false_val;
} State;

typedef struct RenderArgs {
    JSValue tex;
    JSValue display_mode;
    JSValue warnings;
} RenderArgs;

// cgo only uses gcc and clang, so __thread portability is not an issue.
__thread State *tls_state = 0;

static State *init_qjs()
{
    if (tls_state) {
        return tls_state;
    }

    JSRuntime *rt = JS_NewRuntime();
    JSContext *ctx = JS_NewContextRaw(rt);
    JS_AddIntrinsicBaseObjects(ctx);
    JS_AddIntrinsicRegExp(ctx);
    JS_AddIntrinsicJSON(ctx);
    js_std_add_helpers(ctx, 0, 0);
    {
        extern JSModuleDef *js_init_module_std(JSContext *ctx, const char *name);
        js_init_module_std(ctx, "std");
    }
    js_std_eval_binary(ctx, qjsc_katex, qjsc_katex_size, 1);
    js_std_eval_binary(ctx, qjsc_api, qjsc_api_size, 0);

    tls_state = calloc(sizeof(State), 1);
    tls_state->rt = rt;
    tls_state->ctx = ctx;
    tls_state->render = JS_NewAtom(ctx, "render");
    tls_state->global_obj = JS_GetGlobalObject(ctx);
    tls_state->false_val = JS_NewBool(ctx, false);
    tls_state->true_val = JS_NewBool(ctx, true);

    return tls_state;
}

size_t render(void *dest, size_t dest_cap, void *src, size_t src_len, Mode mode)
{
    State *state = init_qjs();
    JSContext *ctx = state->ctx;

    size_t dest_len = 0;
    const char *buf = 0;

    RenderArgs args;
    args.tex = JS_NewStringLen(ctx, src, src_len);
    args.display_mode = (mode & Mode_Display) ? state->true_val : state->false_val;
    args.warnings = (mode & Mode_Warn) ? state->true_val : state->false_val;
    JSValue v = JS_Invoke(ctx, state->global_obj, state->render, 3, &args.tex);

    if (JS_IsString(v) == false) {
        dest_len = -1;
        goto done;
    }

    buf = JS_ToCStringLen(ctx, &dest_len, v);

    if (dest_len > dest_cap) {
        // QJS strings are not UTF-8, so although we can usually tell the buffer
        // is too small to use before converting, we don't know how long it will
        // be. So, we always convert to UTF-8 so we can report to the caller how
        // much they should allocate.
        goto done;
    }

    memcpy(dest, buf, dest_len);

done:
    JS_FreeValue(ctx, args.tex);
    JS_FreeValue(ctx, v);
    JS_FreeCString(ctx, buf);

    return dest_len;
}
