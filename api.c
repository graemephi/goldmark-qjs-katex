#include "api.h"

#include <string.h>

#include "quickjs/quickjs-libc.h"
#include "quickjs/quickjs.h"

#include "api.bytecode.h"

_Thread_local bool initialised = false;
_Thread_local JSRuntime *rt;
_Thread_local JSContext *ctx;
_Thread_local JSAtom renderInlineMaths;
_Thread_local JSAtom renderDisplayMaths;

static void init_qjs()
{
    if (initialised) {
        return;
    }

    rt = JS_NewRuntime();
    ctx = JS_NewContextRaw(rt);
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

    renderInlineMaths = JS_NewAtom(ctx, "renderInlineMaths");
    renderDisplayMaths = JS_NewAtom(ctx, "renderDisplayMaths");

    initialised = true;
}

size_t render_maths(void *dest, size_t dest_cap, void *src, size_t src_len, DisplayMode mode)
{
    init_qjs();

    size_t dest_len = 0;

    JSAtom fn = (mode == DisplayMode_Inline) ? renderInlineMaths : renderDisplayMaths;
    JSValue global_obj = JS_GetGlobalObject(ctx);
    JSValue jsSrc = JS_NewStringLen(ctx, src, src_len);
    JSValue v = JS_Invoke(ctx, global_obj, fn, 1, &jsSrc);

    if (JS_IsString(v) == false) {
        // Should never happen.
        dest_len = -1;
        goto done;
    }

    const char *buf = JS_ToCStringLen(ctx, &dest_len, v);

    if (dest_len > dest_cap) {
        // QJS strings are not UTF-8, so although we can usually tell the buffer
        // is too small to use before converting, we don't know how long it will
        // be. So, we always convert to UTF-8 so we can report to the caller how
        // much they should allocate.
        goto done;
    }

    memcpy(dest, buf, dest_len);

done:
    JS_FreeValue(ctx, global_obj);
    JS_FreeValue(ctx, jsSrc);
    JS_FreeValue(ctx, v);

    return dest_len;
}
