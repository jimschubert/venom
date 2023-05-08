package venom

type TemplateOptions struct {
	Logger              Logger
	JsonMarshaler       MarshalFn
	YamlMarshaler       MarshalFn
	StripAnsiInMarkdown bool
}
