package venom

import (
	"fmt"
	"github.com/jimschubert/stripansi"
	"github.com/jimschubert/venom/internal"
	"strings"
	"unicode"
)

type functionsMarkdown struct {
	stripAnsi bool
}

func (m functionsMarkdown) FormatOptions(input string) string {
	tmp := strings.FieldsFunc(input, func(r rune) bool {
		return '\n' == r
	})
	lines := make([]string, 0)
	for _, s := range tmp {
		lines = append(lines, strings.TrimLeftFunc(s, unicode.IsSpace))
	}
	return strings.Join(lines, "\n")
}

func (m functionsMarkdown) FormatHeader(input string) string {
	return input
}

func (m functionsMarkdown) FormatText(input string) string {
	// TODO: consider cleaning HTML characters here, for now assume people do the right thing by indenting any code blocks in Long descriptions
	if !m.stripAnsi {
		return input
	} else {
		return stripansi.String(input)
	}
}

func (m functionsMarkdown) FormatFlag(input Flag) string {
	return strings.TrimSuffix(input.RawUsage, "\n")
}

func (m functionsMarkdown) SeeAlsoPath(input string) string {
	return internal.CleanPath(input)
}

func (m functionsMarkdown) FormatExample(input string) string {
	// first, handles code blocks with triple tick
	replaced := strings.TrimPrefix(strings.TrimSuffix(input, "\n```"), "```\n")
	// next, handles single-line code with triple tick
	return fmt.Sprintf("```\n%s\n```", strings.TrimPrefix(strings.TrimSuffix(replaced, "```"), "```"))
}

func (m functionsMarkdown) FormatAutoGenTag(input string) string {
	return input
}

func (m functionsMarkdown) IsLocalFlag(input Flag) bool {
	return !input.Persistent && !input.Inherited
}

var (
	_ functions = (*functionsMarkdown)(nil)
)
