package venom

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeJson(outDir string, doc Documentation, options TemplateOptions) error {
	data, err := options.JsonMarshaler(&doc)
	if err != nil {
		return err
	}

	cleanName := CleanPath(doc.RootCommand.Name)
	docRoot := filepath.Join(outDir, cleanName)
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	docJson := filepath.Join(outDir, cleanName, fmt.Sprintf("%s.json", cleanName))
	return os.WriteFile(docJson, data, 0700)
}

func init() {
	RegisterWriter(Json, writeJson)
}