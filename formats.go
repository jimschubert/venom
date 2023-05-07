package venom

// Formats defines the flag of supported documentation formats
type Formats byte

const (
	// Markdown will result in Markdown/CommonMark style output
	Markdown Formats = 1 << iota
	// Man will result in manpage format
	Man
	// Yaml will result in the YAML 1.1+ format
	Yaml
	// ReST will result in Restructured Text format
	ReST
	// Json will result in JavaScript Object Notation (JSON) format
	Json
)

// IsSet determines if the desired flag(s) are set
func (f *Formats) IsSet(format Formats) bool {
	return (*f)&format != 0
}

// Set will define the desired flag(s)
func (f *Formats) Set(format Formats) *Formats {
	*f = (*f) | format
	return f
}

// Unset will remove the desired flag(s)
func (f *Formats) Unset(format Formats) *Formats {
	*f = (*f) &^ format
	return f
}

// IsValid determines if this set of Formats flags are valid; anything set but not defined in the Formats flag set will return false.
func (f *Formats) IsValid() bool {
	return f.IsSet(Markdown) || f.IsSet(Man) || f.IsSet(Yaml) || f.IsSet(ReST) || f.IsSet(Json)
}
