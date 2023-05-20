package venom

import (
	"embed"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
)

//go:embed templates/*.tmpl
var templates embed.FS

// TemplateOptions are those options provided to the templating system
type TemplateOptions struct {
	Logger                   Logger
	JsonMarshaler            MarshalFn
	YamlMarshaler            MarshalFn
	StripAnsiInMarkdown      bool
	MaxOptionWidthInMarkdown int
	Templates                fs.FS
}

// Options provides a builder-pattern of user-facing optional functionality when constructing via venom.Initialize
type Options struct {
	commandName               string
	formats                   Formats
	outDir                    string
	showHiddenCommands        bool
	disableUserCommandOptions bool
	templateOptions           *TemplateOptions
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

// WithOutDirectory allows the caller to define an output directory which differs from the default ./docs.
func (o *Options) WithOutDirectory(out string) *Options {
	o.outDir = out
	return o
}

// WithShowHiddenCommands allows the caller to signify whether details about hidden commands should be present in the final output
func (o *Options) WithShowHiddenCommands() *Options {
	o.showHiddenCommands = true
	return o
}

// WithStripAnsiInMarkdown allows the caller to require ANSI characters to be stripped when processing markdown files.
func (o *Options) WithStripAnsiInMarkdown() *Options {
	o.templateOptions.StripAnsiInMarkdown = true
	return o
}

// DisableUserCommandOptions allows the caller to define fixed options such as output directory, supported doc formats, etc.
// Default behavior is a documentation command which provides defaults defined by the caller, but exposing a subset of options to the user.
func (o *Options) DisableUserCommandOptions() *Options {
	o.disableUserCommandOptions = true
	return o
}

// WithCustomTemplates allows the user to provide an implementation which provide custom templates for any template-driven format.
func (o *Options) WithCustomTemplates(fs fs.FS) *Options {
	o.templateOptions.Templates = fs
	return o
}

func (o *Options) WithMaxOptionWidthInMarkdown(width int) *Options {
	o.templateOptions.MaxOptionWidthInMarkdown = width
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
		return errors.New("invalid yaml marshal provided")
	}

	if o.templateOptions.Logger == nil {
		return errors.New("invalid logger provided")
	}

	if o.templateOptions.Templates == nil {
		return errors.New("invalid templates provided")
	}

	if o.templateOptions.MaxOptionWidthInMarkdown < 24 {
		return errors.New("invalid max options width in markdown provided; minimum is 24")
	}
	return nil
}

// NewOptions provides a new set of options with default command name ("docs") and formats (Markdown).
// The value returned here follows the builder pattern for easily discovering and applying available options.
func NewOptions() *Options {
	return &Options{
		commandName: "docs",
		formats:     Markdown,
		outDir:      "docs",
		templateOptions: &TemplateOptions{
			JsonMarshaler:            json.Marshal,
			YamlMarshaler:            yaml.Marshal,
			Logger:                   log.Default(),
			Templates:                templates,
			MaxOptionWidthInMarkdown: 120,
		},
	}
}
