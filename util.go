package venom

import (
	"bytes"
	"strings"
	"unicode"
)

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func hangingIndent(input string, hangWidth int, maxWidth int) string {
	if len(input) == maxWidth {
		return input
	}

	buf := bytes.Buffer{}
	word := bytes.Buffer{}
	available := maxWidth - 1

	boundary := 0
	space := strings.Repeat(" ", hangWidth)
	for _, i := range strings.TrimRight(input, " ") {
		if unicode.IsSpace(i) {
			// word boundary, copy to buf
			if word.Len() > 0 {
				buf.Write(word.Bytes())
				word.Reset()
				buf.WriteRune(i)
				boundary = 1 // 'i' is one whitespace on buf.
			} else {
				word.WriteRune(i)
			}

			if i == '\n' {
				// reset available
				if boundary > 0 {
					buf.Truncate(buf.Len() - boundary)
				}
				available = maxWidth - 1 - len(space)
				boundary = 0

				buf.WriteString(space)
			} else {
				available--
			}

			continue
		}

		if available <= 1 {
			// trim any trailing whitespace from the buffer before writing the newline
			if boundary > 0 {
				buf.Truncate(buf.Len() - boundary)
			}
			available = maxWidth - 1 - len(space)

			buf.WriteRune('\n')
			buf.WriteString(space)
		}
		word.WriteRune(i)
		available--
	}
	if word.Len() > 0 {
		buf.Write(word.Bytes())
	}
	return buf.String()
}
