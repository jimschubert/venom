package venom

import (
	"fmt"
	"github.com/jimschubert/venom/internal"
	"os"
	"path/filepath"
	"text/template"
)

type templateWriter struct {
	name          string
	fileExtension string
	outDir        string
	doc           Documentation
	options       TemplateOptions
	funcs         functions
	includeIndex  bool
}

func (w *templateWriter) write() error {
	t, err := template.New(w.name).Funcs(newFuncMap(w.funcs)).ParseFS(w.options.Templates, "**/*.tmpl")
	if err != nil {
		return err
	}

	docRoot := filepath.Join(w.outDir, internal.CleanPath(w.doc.RootCommand.Name))
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	rootCommandPath := filepath.Join(docRoot, fmt.Sprintf("%s.%s", internal.CleanPath(w.doc.RootCommand.Name), w.fileExtension))
	rootCommand, err := os.Create(rootCommandPath)
	defer func(f *os.File) {
		_ = f.Close()
	}(rootCommand)

	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(rootCommand, fmt.Sprintf("%s_command.tmpl", w.name), struct {
		Command
		Doc Documentation
	}{
		Command: w.doc.RootCommand,
		Doc:     w.doc,
	}); err != nil {
		return err
	}

	w.options.Logger.Printf("[%s] Wrote file %s", w.name, rootCommandPath)

	var writeCommand func(c Command, t *template.Template) error
	writeCommand = func(c Command, t *template.Template) error {
		subCommandPath := filepath.Join(docRoot, fmt.Sprintf("%s.%s", internal.CleanPath(c.FullPath), w.fileExtension))
		subCommandFile, err := os.Create(subCommandPath)
		if err != nil {
			return err
		}
		defer func(subCommandFile *os.File) {
			_ = subCommandFile.Close()
		}(subCommandFile)

		if err := t.ExecuteTemplate(subCommandFile, fmt.Sprintf("%s_command.tmpl", w.name), struct {
			Command
			Doc Documentation
		}{
			Command: c,
			Doc:     w.doc,
		}); err != nil {
			return err
		}

		w.options.Logger.Printf("[%s] Wrote file %s", w.name, subCommandPath)

		for _, subcommand := range c.Subcommands {
			err = writeCommand(subcommand, t)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if w.doc.RootCommand.Subcommands != nil {
		for _, subcommand := range w.doc.RootCommand.Subcommands {
			err = writeCommand(subcommand, t)
			if err != nil {
				return err
			}
		}
	}

	if w.includeIndex {
		var indexName string
		if w.doc.RootCommand.Name == "index" {
			indexName = "README"
		} else {
			indexName = "index"
		}
		indexPath := filepath.Join(docRoot, fmt.Sprintf("%s.%s", indexName, w.fileExtension))
		index, err := os.Create(indexPath)
		defer func(f *os.File) {
			_ = f.Close()
		}(index)

		if err != nil {
			return err
		}

		err = t.ExecuteTemplate(index, fmt.Sprintf("%s_index.tmpl", w.name), w.doc)

		if err == nil {
			w.options.Logger.Printf("[%s] Wrote file %s", w.name, indexPath)
		}
		return err
	}

	return nil
}