package venom

import "text/template"

// functions defines the common set of functions for template providers
type functions interface {
	FormatHeader(input string) string
	FormatText(input string) string
	FormatOptions(input string) string
	FormatFlag(input Flag) string
	SeeAlsoPath(input string) string
	FormatExample(input string) string
	FormatAutoGenTag(input string) string
	IsLocalFlag(input Flag) bool
}

func newFuncMap(fns functions) template.FuncMap {
	return template.FuncMap{
		"header":        fns.FormatHeader,
		"text":          fns.FormatText,
		"options":       fns.FormatOptions,
		"flag":          fns.FormatFlag,
		"see_also_path": fns.SeeAlsoPath,
		"example":       fns.FormatExample,
		"autogen":       fns.FormatAutoGenTag,
		"is_local":      fns.IsLocalFlag,
	}
}
