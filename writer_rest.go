package venom

import (
	"fmt"
	"github.com/jimschubert/stripansi"
	"strings"
)

type functionsRest struct {
}

func (f functionsRest) FormatHeader(input string) string {
	return strings.ReplaceAll(input, " ", "_")
}

func (f functionsRest) FormatText(input string) string {
	return stripansi.String(input)
}

func (f functionsRest) FormatOptions(input string) string {
	return f.indented(input)
}

func (f functionsRest) FormatFlag(input Flag) string {
	return f.indented(input.Usage)
}

func (f functionsRest) SeeAlsoPath(input string) string {
	ref := strings.ReplaceAll(input, " ", "_")
	return fmt.Sprintf("`%s <%s.rst>`_", input, ref)
}

func (f functionsRest) FormatExample(input string) string {
	return f.indented(trimIndent(input, -1))
}

func (f functionsRest) FormatAutoGenTag(input string) string {
	return input
}

func (f functionsRest) IsLocalFlag(input Flag) bool {
	return !input.Persistent && !input.Inherited
}

func (f functionsRest) indented(input string) string {
	var result []rune
	addIndent := true
	for _, char := range input {
		switch char {
		// trims space at the start of each line!
		case '\t', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			if addIndent {
				continue
			}
		case '\n':
			if !addIndent && len(result) > 0 {
				addIndent = true
			}
		}

		if addIndent {
			result = append(result, ' ', ' ')
			addIndent = false
		}

		result = append(result, char)
	}
	return string(result)
}

type writerRest struct {
	options TemplateOptions
}

func (w *writerRest) Write(outDir string, doc Documentation) error {
	fns := functionsRest{}

	helper := writerForTemplates{
		name:          ReST.String(),
		fileExtension: "rst",
		outDir:        outDir,
		doc:           doc,
		options:       w.options,
		funcs:         fns,
		includeIndex:  false,
	}

	return helper.write()
}

func (w *writerRest) SetTemplateOptions(options TemplateOptions) {
	w.options = options
}

func init() {
	registerWriter(ReST, func() writer {
		return &writerRest{}
	})
}
