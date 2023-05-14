package venom

type writerJson struct {
	options TemplateOptions
}

func (w *writerJson) SetTemplateOptions(options TemplateOptions) {
	w.options = options
}

func (w *writerJson) Write(outDir string, doc Documentation) error {
	helper := writerForMarshals{
		name:          Json.String(),
		fileExtension: "json",
		outDir:        outDir,
		doc:           doc,
		marshaller:    w.options.JsonMarshaler,
		logger:        w.options.Logger,
	}
	return helper.write()
}

func init() {
	registerWriter(Json, func() writer {
		return &writerJson{}
	})
}
