package venom

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
)

// Options provides a builder-pattern of user-facing optional functionality when constructing via venom.Initialize
type Options struct {
	commandName        string
	formats            Formats
	outDir             string
	ShowHiddenCommands bool
	templateOptions    *TemplateOptions
}

// WithCommandName allows the caller to provide the target command name, which is used to construct a cobra.Command for documentation.
func (o *Options) WithCommandName(name string) *Options {
	o.commandName = name
	return o
}

// WithLogger allows the caller to define a target log implementation for any warnings or errors from the initialized command.
func (o *Options) WithLogger(logger Logger) *Options {
	o.templateOptions.Logger = logger
	return o
}

// WithFormats allows the caller to define formats which differ from the Options defaults.
func (o *Options) WithFormats(formats Formats) *Options {
	o.formats = formats
	return o
}

// WithJsonMarshal allows the caller to define a custom JSON marshaling function, default is json.Marshal from the Go standard library.
func (o *Options) WithJsonMarshal(fn MarshalFn) *Options {
	o.templateOptions.JsonMarshaler = fn
	return o
}

// WithYamlMarshal allows the caller to define a custom YAML marshaling function, default is yaml.Marshal (v3) as depended by cobra.
func (o *Options) WithYamlMarshal(fn MarshalFn) *Options {
	o.templateOptions.YamlMarshaler = fn
	return o
}

// WithOutDirectory allows the caller to define an output directory which differs from the default dist/docs.
func (o *Options) WithOutDirectory(out string) *Options {
	o.outDir = out
	return o
}

// WithShowHiddenCommands allows the caller to signify whether details about hidden commands should be present in the final output
func (o *Options) WithShowHiddenCommands() *Options {
	o.ShowHiddenCommands = true
	return o
}

// TemplateOptions provides the value of current TemplateOptions
func (o *Options) TemplateOptions() TemplateOptions {
	return *(*o).templateOptions
}

func (o *Options) validate() error {
	if o.commandName == "" {
		return errors.New("command name can't be empty")
	}

	if !o.formats.IsValid() {
		return errors.New("invalid documentation format(s) provided")
	}

	if o.templateOptions.JsonMarshaler == nil {
		return errors.New("invalid json marshal provided")
	}

	if o.templateOptions.YamlMarshaler == nil {
		return errors.New("invalid yaml marsha provided")
	}

	return nil
}

// NewOptions provides a new set of options with default command name ("docs") and formats (Markdown).
// The value returned here follows the builder pattern for easily discovering and applying available options.
func NewOptions() *Options {
	return &Options{
		commandName: "docs",
		formats:     Markdown,
		outDir:      "dist/docs",
		templateOptions: &TemplateOptions{
			JsonMarshaler: json.Marshal,
			YamlMarshaler: yaml.Marshal,
		},
	}
}
