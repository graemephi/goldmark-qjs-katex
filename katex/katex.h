#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

typedef enum Mode
{
    Mode_Inline,
    Mode_Display,
	Mode_InlineWarn,
	Mode_DisplayWarn,

    Mode_Warn = Mode_InlineWarn
} Mode;

// Returns length of resulting string on sucess, -1 on failure. If the length is
// too large to fit in the destination buffer, no bytes will be written, but the
// length that would have been written otherwise is returned.
size_t render(void *dest, size_t dest_cap, void *src, size_t src_len, Mode mode);
