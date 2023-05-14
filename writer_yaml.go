package venom

type writerYaml struct {
	options TemplateOptions
}

func (w *writerYaml) SetTemplateOptions(options TemplateOptions) {
	w.options = options
}

func (w *writerYaml) Write(outDir string, doc Documentation) error {
	helper := writerForMarshals{
		name:          Yaml.String(),
		fileExtension: "yml",
		outDir:        outDir,
		doc:           doc,
		marshaller:    w.options.YamlMarshaler,
		logger:        w.options.Logger,
	}
	return helper.write()
}

func init() {
	registerWriter(Yaml, func() writer {
		return &writerYaml{}
	})
}
