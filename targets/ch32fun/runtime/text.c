#include "text.h"

#include <string.h>

/*
 * Bounds-checked character access for immutable strings. Returns -1 when the
 * index is negative or past the terminating NUL, so callers can test the
 * result without risking an out-of-bounds read.
 */
int32_t jikascript_char_at(const char *text, int32_t index) {
    if (text == 0 || index < 0 || (size_t)index >= strlen(text)) {
        return -1;
    }
    return (int32_t)(unsigned char)text[index];
}
