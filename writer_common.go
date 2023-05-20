package venom

import (
	"fmt"
	"github.com/jimschubert/venom/internal"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type writerForTemplates struct {
	name          string
	fileExtension string
	outDir        string
	doc           Documentation
	options       TemplateOptions
	funcs         functions
	includeIndex  bool
}

func (w *writerForTemplates) write() error {
	templateName := strings.ToLower(w.name)
	t, err := template.New(w.name).Funcs(newFuncMap(w.funcs)).ParseFS(w.options.Templates, "**/*.tmpl")
	if err != nil {
		return err
	}

	indexTemplateName := fmt.Sprintf("%s_index.tmpl", templateName)
	commandTemplateName := fmt.Sprintf("%s_command.tmpl", templateName)

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

	if t.Lookup(commandTemplateName) != nil {

		if err := t.ExecuteTemplate(rootCommand, commandTemplateName, struct {
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

			if err := t.ExecuteTemplate(subCommandFile, commandTemplateName, struct {
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
	} else {
		w.options.Logger.Printf("[%s] Skipping commands: no template found for %q", w.name, commandTemplateName)
	}

	if w.includeIndex {
		// if the writer supports an index, but user has customized without the targeted index, we just skip and log
		if t.Lookup(indexTemplateName) != nil {
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

			err = t.ExecuteTemplate(index, indexTemplateName, w.doc)

			if err == nil {
				w.options.Logger.Printf("[%s] Wrote file %s", w.name, indexPath)
			}
			return err
		} else {
			w.options.Logger.Printf("[%s] Skipping index: no template found for %q", w.name, indexTemplateName)
		}
	}

	return nil
}

type writerForMarshals struct {
	name          string
	fileExtension string
	outDir        string
	doc           Documentation
	marshaller    MarshalFn
	logger        Logger
}

func (w *writerForMarshals) write() error {
	data, err := w.marshaller(&w.doc)
	if err != nil {
		return err
	}

	cleanName := internal.CleanPath(w.doc.RootCommand.Name)
	docRoot := filepath.Join(w.outDir, cleanName)
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	docJson := filepath.Join(w.outDir, cleanName, fmt.Sprintf("%s.%s", cleanName, w.fileExtension))
	err = os.WriteFile(docJson, data, 0700)
	if err == nil {
		w.logger.Printf("[%s] Wrote file %s", w.name, docJson)
	}
	return err
}

func trimIndent(input string, max int) string {
	tmp := strings.FieldsFunc(input, func(r rune) bool {
		return '\n' == r
	})
	lines := make([]string, 0)
	for _, s := range tmp {
		spaceCount := max
		current := strings.TrimLeftFunc(s, func(r rune) bool {
			shouldDrop := (0 < spaceCount) && unicode.IsSpace(r)
			spaceCount--
			return shouldDrop
		})
		if len(current) > 0 {
			lines = append(lines, current)
		}
	}

	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}
