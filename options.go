package venom

import (
	"errors"
)

// Options provides a builder-pattern of user-facing optional functionality when constructing via venom.Initialize
type Options struct {
	commandName string
	formats     Formats
	logger      Logger
}

// WithCommandName allows the caller to provide the target command name, which is used to construct a cobra.Command for documentation.
func (o *Options) WithCommandName(name string) *Options {
	o.commandName = name
	return o
}

// WithLogger allows the caller to define a target log implementation for any warnings or errors from the initialized command.
func (o *Options) WithLogger(logger Logger) *Options {
	o.logger = logger
	return o
}

// WithFormats allows the caller to define formats which differ from the Options defaults.
func (o *Options) WithFormats(formats Formats) *Options {
	o.formats = formats
	return o
}

func (o *Options) validate() error {
	if o.commandName == "" {
		return errors.New("command name can't be empty")
	}

	if !o.formats.IsValid() {
		return errors.New("invalid documentation format(s) provided")
	}

	return nil
}

// NewOptions provides a new set of options with default command name ("docs") and formats (Markdown).
// The value returned here follows the builder pattern for easily discovering and applying available options.
func NewOptions() *Options {
	return &Options{
		commandName: "docs",
		formats:     Markdown,
	}
}
