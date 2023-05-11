package internal

import "regexp"

var nonPathRegex = regexp.MustCompile(`[^a-zA-Z0-9.\-]+`)

// CleanPath will replace all invalid characters in input with underscore by default. A varargs is abused here to allow
// for passing a single optional replace character (acting as a function overload).
func CleanPath(input string, replace ...string) string {
	var replaceWith string
	if len(replace) > 1 {
		panic("multiple replace characters are not supported")
	}

	if len(replace) == 1 {
		replaceWith = replace[0]
	} else {
		replaceWith = "_"
	}

	return nonPathRegex.ReplaceAllString(input, replaceWith)
}
