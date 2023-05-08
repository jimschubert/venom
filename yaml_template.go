package venom

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeYaml(outDir string, doc Documentation, options TemplateOptions) error {
	data, err := options.YamlMarshaler(&doc)
	if err != nil {
		return err
	}

	cleanName := CleanPath(doc.RootCommand.Name)
	docRoot := filepath.Join(outDir, cleanName)
	if err := os.MkdirAll(docRoot, 0700); err != nil {
		return err
	}

	docYaml := filepath.Join(outDir, cleanName, fmt.Sprintf("%s.yml", cleanName))
	return os.WriteFile(docYaml, data, 0700)
}

func init() {
	RegisterWriter(Yaml, writeYaml)
}
