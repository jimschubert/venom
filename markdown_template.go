package venom

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

func writeMarkdown(outDir string, doc Documentation, options TemplateOptions) error {
	fns := MarkdownFns
	t, err := template.New("markdown").Funcs(template.FuncMap{
		"header":        fns.FormatHeader,
		"text":          fns.FormatText,
		"flag":          fns.FormatFlag,
		"see_also_path": fns.SeeAlsoPath,
		"example":       fns.FormatExample,
		"autogen":       fns.FormatAutoGenTag,
		"is_local":      fns.IsLocalFlag,
	}).ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		return err
	}

	docRoot := filepath.Join(outDir, CleanPath(doc.RootCommand.Name))
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	rootCommand, err := os.Create(filepath.Join(docRoot, fmt.Sprintf("%s.md", CleanPath(doc.RootCommand.Name))))
	defer func(f *os.File) {
		_ = f.Close()
	}(rootCommand)

	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(rootCommand, "markdown_command.tmpl", struct {
		Command
		Doc Documentation
	}{
		Command: doc.RootCommand,
		Doc:     doc,
	}); err != nil {
		return err
	}

	var writeCommand func(c Command, t *template.Template) error
	writeCommand = func(c Command, t *template.Template) error {
		subCommandFile, err := os.Create(filepath.Join(docRoot, fmt.Sprintf("%s.md", CleanPath(c.FullPath))))
		if err != nil {
			return err
		}
		defer func(subCommandFile *os.File) {
			_ = subCommandFile.Close()
		}(subCommandFile)

		if err := t.ExecuteTemplate(subCommandFile, "markdown_command.tmpl", struct {
			Command
			Doc Documentation
		}{
			Command: c,
			Doc:     doc,
		}); err != nil {
			return err
		}

		for _, subcommand := range c.Subcommands {
			err = writeCommand(subcommand, t)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if doc.RootCommand.Subcommands != nil {
		for _, subcommand := range doc.RootCommand.Subcommands {
			err = writeCommand(subcommand, t)
			if err != nil {
				return err
			}
		}
	}

	var indexName string
	if doc.RootCommand.Name == "index" {
		indexName = "README"
	} else {
		indexName = "index"
	}
	index, err := os.Create(filepath.Join(docRoot, fmt.Sprintf("%s.md", indexName)))
	defer func(f *os.File) {
		_ = f.Close()
	}(index)

	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(index, "markdown_index.tmpl", doc)

	return err
}

func init() {
	RegisterWriter(Markdown, writeMarkdown)
}