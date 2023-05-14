package venom

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type writerFn func() writer

var writers = make(map[Formats]writerFn)

type docCommandOptions struct {
	outDir     string
	formats    []string
	showHidden bool
}

// Initialize a new documentation command with cmd as the parent, providing options for customization
// The provided command will always provide documentation from the *root*, which allows the caller to list this automated
// command under other administrative/tooling/hidden commands as needed.
func Initialize(cmd *cobra.Command, options *Options) error {
	if options == nil {
		options = NewOptions()
	}

	if err := options.validate(); err != nil {
		return err
	}

	formats := make([]string, 0)
	for _, format := range options.formats.defined() {
		switch format {
		case Markdown:
			formats = append(formats, "markdown")
		case Yaml:
			formats = append(formats, "yaml")
		case Json:
			formats = append(formats, "json")
		case ReST:
			formats = append(formats, "rest")
		}
	}

	o := docCommandOptions{
		outDir:     options.outDir,
		showHidden: options.showHiddenCommands,
		formats:    formats,
	}

	docCommand := &cobra.Command{
		Use:    options.commandName,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := *options
			root := cmd.Root()
			root.InitDefaultHelpCmd()
			root.InitDefaultHelpFlag()

			opts.showHiddenCommands = o.showHidden
			opts.outDir = o.outDir

			definedFormats := getUserSelectedFormats(o, opts)

			if !definedFormats.IsValid() {
				return errors.New("invalid formats selected")
			}

			opts.formats = definedFormats

			documentation := NewDocumentation(root, &opts)
			if err := Write(documentation); err != nil {
				return err
			}

			return nil
		},
	}

	if !options.disableUserCommandOptions {
		docCommand.Flags().StringVar(&o.outDir, "out-dir", o.outDir, "The target output directory")
		docCommand.Flags().BoolVar(&o.showHidden, "show-hidden", o.showHidden, "Also show hidden commands")
		docCommand.Flags().StringSliceVar(&o.formats, "formats", o.formats,
			fmt.Sprintf("A comma-separated list of formats to output. Allowed: [%s]", strings.Join(formats, ",")))
	}

	cmd.AddCommand(docCommand)

	return nil
}

func getUserSelectedFormats(o docCommandOptions, opts Options) Formats {
	definedFormats := Formats(0)
	for _, format := range o.formats {
		switch strings.ToLower(format) {
		case "markdown", "md":
			if opts.formats.IsSet(Markdown) {
				definedFormats.Set(Markdown)
			} else {
				opts.templateOptions.Logger.Printf("Skipping markdown documentation because the application maintainers have not enabled this output format.")
			}
		case "yaml", "yml":
			if opts.formats.IsSet(Yaml) {
				definedFormats.Set(Yaml)
			} else {
				opts.templateOptions.Logger.Printf("Skipping yaml documentation because the application maintainers have not enabled this output format.")
			}
		case "json":
			if opts.formats.IsSet(Json) {
				definedFormats.Set(Json)
			} else {
				opts.templateOptions.Logger.Printf("Skipping json documentation because the application maintainers have not enabled this output format.")
			}
		case "rest", "rst":
			if opts.formats.IsSet(ReST) {
				definedFormats.Set(ReST)
			} else {
				opts.templateOptions.Logger.Printf("Skipping rest documentation because the application maintainers have not enabled this output format.")
			}
		default:
			opts.templateOptions.Logger.Printf("Skipping %s documentation because it is not currently supported.", format)
		}
	}
	return definedFormats
}

// registerWriter for a given format to allow writing via the writer function
func registerWriter(format Formats, writer writerFn) {
	writers[format] = writer
}

// Write to outDir the documentation for all given formats
func Write(documentation Documentation) error {
	var err error
	options := documentation.options
	formats := options.formats
	outDir := options.outDir

	if !formats.IsValid() {
		return errors.New("unexpected formats provided to Write")
	}

	// ensure proper initialization
	documentation.init()

	templateOptions := options.TemplateOptions()
	for _, format := range []Formats{Yaml, Json, Markdown, Man, ReST} {
		if formats.IsSet(format) {
			templateOptions.Logger.Printf("Generating documentation for %s", strings.ToLower(format.String()))
			if writeBuilder, ok := writers[format]; ok {
				w := writeBuilder()
				if wt, ok := w.(wantsTemplateOptions); ok {
					wt.SetTemplateOptions(templateOptions)
				}
				err = w.Write(outDir, documentation)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("missing output writer for format %v", format)
			}
		}
	}
	return nil
}
