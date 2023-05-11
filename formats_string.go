package venom

import (
	"bytes"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Markdown-1]
	_ = x[Man-2]
	_ = x[Yaml-4]
	_ = x[ReST-8]
	_ = x[Json-16]
}

const (
	_Formats_name_0 = "MarkdownMan"
	_Formats_name_1 = "Yaml"
	_Formats_name_2 = "ReST"
	_Formats_name_3 = "Json"
)

var (
	_Formats_index_0 = [...]uint8{0, 8, 11}
)

func (i Formats) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _Formats_name_0[_Formats_index_0[i]:_Formats_index_0[i+1]]
	case i == 4:
		return _Formats_name_1
	case i == 8:
		return _Formats_name_2
	case i == 16:
		return _Formats_name_3
	default:
		buf := bytes.Buffer{}
		d := i.defined()
		last := len(d) - 1
		for i, format := range d {
			buf.WriteString(format.String())
			if i < last {
				buf.WriteString("|")
			}
		}
		return buf.String()
	}
}
