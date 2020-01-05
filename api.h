#include <stdbool.h>
#include <stdint.h>

typedef enum DisplayMode
{
    DisplayMode_Inline,
    DisplayMode_Display
} DisplayMode;

// Returns length of resulting string on sucess, -1 on failure. If the length is
// too large to fit in the destination buffer, no bytes will be written, but the
// length that would have been written otherwise is returned.
size_t render_maths(void *dest, size_t dest_cap, void *src, size_t src_len, DisplayMode mode);
