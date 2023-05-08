package venom

import (
	"fmt"
	"strings"
)

type markdownFunctions struct {
}

func (m markdownFunctions) FormatHeader(input string) string {
	return input
}

func (m markdownFunctions) FormatText(input string) string {
	replacer := strings.NewReplacer(
		" < ", " &lt; ",
		" > ", " &gt; ",
	)

	return replacer.Replace(input)
}

func (m markdownFunctions) FormatFlag(input Flag) string {
	return strings.TrimSuffix(input.RawUsage, "\n")
}

func (m markdownFunctions) SeeAlsoPath(input string) string {
	return CleanPath(input)
}

func (m markdownFunctions) FormatExample(input string) string {
	// first, handles code blocks with triple tick
	replaced := strings.TrimPrefix(strings.TrimSuffix(input, "\n```"), "```\n")
	// next, handles single-line code with triple tick
	return fmt.Sprintf("```\n%s\n```", strings.TrimPrefix(strings.TrimSuffix(replaced, "```"), "```"))
}

func (m markdownFunctions) FormatAutoGenTag(input string) string {
	return input
}

func (m markdownFunctions) IsLocalFlag(input Flag) bool {
	return !input.Persistent && !input.Inherited
}

var (
	MarkdownFns           = markdownFunctions{}
	_           Functions = (*markdownFunctions)(nil)
)
