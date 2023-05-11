package venom

import "io/fs"

type TemplateOptions struct {
	Logger              Logger
	JsonMarshaler       MarshalFn
	YamlMarshaler       MarshalFn
	StripAnsiInMarkdown bool
	Templates           fs.FS
}
