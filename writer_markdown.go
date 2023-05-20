package venom

import (
	"fmt"
	"github.com/jimschubert/stripansi"
	"github.com/jimschubert/venom/internal"
	"strings"
)

type functionsMarkdown struct {
	stripAnsi      bool
	maxOptionWidth int
}

func (m functionsMarkdown) FormatOptions(input string) string {
	leftAligned := trimIndent(input, 2)
	width := 0

loop:
	for i := 0; i < len(leftAligned); i++ {
		first := rune(leftAligned[i])
		// lookahead
		for j := 1; j < 3 && (i+j < len(leftAligned)); j++ {
			idx := i + j
			if ' ' == first && ' ' == rune(leftAligned[idx]) {
				last := rune(leftAligned[idx+1])
				if ' ' != last && '-' != last {
					width = idx + 1
					break loop
				}
			}
		}
	}

	columnWrapped := hangingIndent(leftAligned, width, m.maxOptionWidth)
	return columnWrapped
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

type writerMarkdown struct {
	options TemplateOptions
}

func (w *writerMarkdown) Write(outDir string, doc Documentation) error {
	fns := functionsMarkdown{
		stripAnsi:      w.options.StripAnsiInMarkdown,
		maxOptionWidth: w.options.MaxOptionWidthInMarkdown,
	}

	helper := writerForTemplates{
		name:          Markdown.String(),
		fileExtension: "md",
		outDir:        outDir,
		doc:           doc,
		options:       w.options,
		funcs:         fns,
		includeIndex:  true,
	}

	return helper.write()
}

func (w *writerMarkdown) SetTemplateOptions(options TemplateOptions) {
	w.options = options
}

func init() {
	registerWriter(Markdown, func() writer {
		return &writerMarkdown{}
	})
}

var (
	_ functions = (*functionsMarkdown)(nil)
)
