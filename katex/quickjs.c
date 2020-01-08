// Compile quickjs as a single compilation unit.
// This simplifies the build process from go.

#define _GNU_SOURCE
#define CONFIG_VERSION "2019-12-21-qjsk"

#include "quickjs/quickjs-libc.c"
#include "quickjs/quickjs.c"

#undef malloc
#undef realloc
#undef free
#define is_digit is_digit2
#define compute_stack_size compute_stack_size2

#include "quickjs/libregexp.c"
#include "quickjs/libunicode.c"
#include "quickjs/cutils.c"
