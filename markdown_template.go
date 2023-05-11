package venom

import (
	_ "embed"
	"fmt"
	"github.com/jimschubert/venom/internal"
	"os"
	"path/filepath"
	"text/template"
)

func writeMarkdown(outDir string, doc Documentation, options TemplateOptions) error {
	fns := markdownFunctions{
		stripAnsi: options.StripAnsiInMarkdown,
	}

	t, err := template.New("markdown").Funcs(newFuncMap(fns)).ParseFS(options.Templates, "**/*.tmpl")
	if err != nil {
		return err
	}

	docRoot := filepath.Join(outDir, internal.CleanPath(doc.RootCommand.Name))
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	rootCommandPath := filepath.Join(docRoot, fmt.Sprintf("%s.md", internal.CleanPath(doc.RootCommand.Name)))
	rootCommand, err := os.Create(rootCommandPath)
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

	options.Logger.Printf("[markdown] Wrote file %s", rootCommandPath)

	var writeCommand func(c Command, t *template.Template) error
	writeCommand = func(c Command, t *template.Template) error {
		subCommandPath := filepath.Join(docRoot, fmt.Sprintf("%s.md", internal.CleanPath(c.FullPath)))
		subCommandFile, err := os.Create(subCommandPath)
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

		options.Logger.Printf("[markdown] Wrote file %s", subCommandPath)

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
	indexPath := filepath.Join(docRoot, fmt.Sprintf("%s.md", indexName))
	index, err := os.Create(indexPath)
	defer func(f *os.File) {
		_ = f.Close()
	}(index)

	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(index, "markdown_index.tmpl", doc)

	if err == nil {
		options.Logger.Printf("[markdown] Wrote file %s", indexPath)
	}

	return err
}

func init() {
	registerWriter(Markdown, writeMarkdown)
}
