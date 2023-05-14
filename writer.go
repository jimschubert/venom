package venom

type writer interface {
	Write(outDir string, doc Documentation) error
}

type wantsTemplateOptions interface {
	SetTemplateOptions(options TemplateOptions)
}
