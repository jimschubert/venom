package venom

import (
	"fmt"
	"github.com/jimschubert/venom/internal"
	"os"
	"path/filepath"
)

func writerJson(outDir string, doc Documentation, options TemplateOptions) error {
	data, err := options.JsonMarshaler(&doc)
	if err != nil {
		return err
	}

	cleanName := internal.CleanPath(doc.RootCommand.Name)
	docRoot := filepath.Join(outDir, cleanName)
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	docJson := filepath.Join(outDir, cleanName, fmt.Sprintf("%s.json", cleanName))
	err = os.WriteFile(docJson, data, 0700)
	if err == nil {
		options.Logger.Printf("[json] Wrote file %s", docJson)
	}
	return err
}

func init() {
	registerWriter(Json, writerJson)
}
