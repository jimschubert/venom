package venom

func writeMarkdown(outDir string, doc Documentation, options TemplateOptions) error {
	fns := functionsMarkdown{
		stripAnsi: options.StripAnsiInMarkdown,
	}

	writer := templateWriter{
		name:          "markdown",
		fileExtension: "md",
		outDir:        outDir,
		doc:           doc,
		options:       options,
		funcs:         fns,
		includeIndex:  true,
	}

	return writer.write()
}

func init() {
	registerWriter(Markdown, writeMarkdown)
}
