package venom

import (
	"github.com/jimschubert/venom/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
	"time"
)

// MarshalFn is a common interface allowing the caller to provide yaml/json marshaler functions
type MarshalFn func(in interface{}) (out []byte, err error)

// Logger allows the user to provide any logger fulfilling this interface
type Logger interface {
	// Printf is a common signature used by log.Logger, logrus.Logger, and others
	Printf(format string, v ...any)
}

// Documentation represents the "top-level" of documentation to be passed to a template
type Documentation struct {
	GenerationDate    string   `yaml:"generationDate,omitempty" json:"generationDate,omitempty"`
	AutoGenerationTag string   `yaml:"autoGenerationTag,omitempty" json:"autoGenerationTag,omitempty"`
	RootCommand       Command  `yaml:"rootCommand,omitempty" json:"rootCommand,omitempty"`
	options           *Options `yaml:"-"`
}

func (d *Documentation) init() {
	if d.GenerationDate == "" {
		d.GenerationDate = time.Now().Format("2-Jan-2006")
	}
}

// Write these docs
func (d *Documentation) Write() error {
	if d != nil {
		return Write(*d)
	}
	return nil
}

// ParentCommand provides the name of a command's parent
type ParentCommand struct {
	Name     string `yaml:"name,omitempty" json:"name,omitempty"`
	Short    string `yaml:"short,omitempty" json:"short,omitempty"`
	FullPath string `yaml:"fullPath,omitempty" json:"fullPath,omitempty"`
}

// Command is a different representation of cobra.Command
type Command struct {
	Name            string            `yaml:"name,omitempty" json:"name,omitempty"`
	Usage           string            `yaml:"usage,omitempty" json:"usage,omitempty"`
	Aliases         []string          `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	SuggestFor      []string          `yaml:"suggestFor,omitempty" json:"suggestFor,omitempty"`
	Short           string            `yaml:"short,omitempty" json:"short,omitempty"`
	Long            string            `yaml:"long,omitempty" json:"long,omitempty"`
	GroupID         string            `yaml:"groupID,omitempty" json:"groupID,omitempty"`
	ValidArgs       []string          `yaml:"validArgs,omitempty" json:"validArgs,omitempty"`
	ArgAliases      []string          `yaml:"argAliases,omitempty" json:"argAliases,omitempty"`
	Deprecated      string            `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Annotations     map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Version         string            `yaml:"version,omitempty" json:"version,omitempty"`
	Hidden          bool              `yaml:"hidden" json:"hidden"`
	Runnable        bool              `yaml:"runnable" json:"runnable"`
	RawFlagUsages   string            `yaml:"rawFlagUsages,omitempty" json:"rawFlagUsages,omitempty"`
	Parent          *ParentCommand    `yaml:"parent,omitempty" json:"parent,omitempty"`
	Subcommands     []Command         `yaml:"subcommands,omitempty" json:"subcommands,omitempty"`
	LocalFlags      []Flag            `yaml:"localFlags,omitempty" json:"localFlags,omitempty"`
	InheritedFlags  []Flag            `yaml:"inheritedFlags,omitempty" json:"inheritedFlags,omitempty"`
	PersistentFlags []Flag            `yaml:"persistentFlags,omitempty" json:"persistentFlags,omitempty"`
	Examples        []string          `yaml:"examples,omitempty" json:"examples,omitempty"`
	FullPath        string            `yaml:"fullPath,omitempty" json:"fullPath,omitempty"`
}

// Flag is a representation of pflag.Flag
type Flag struct {
	Name                string `json:"name,omitempty" yaml:"name,omitempty"`
	Shorthand           string `json:"shorthand,omitempty" yaml:"shorthand,omitempty"`
	Usage               string `json:"usage,omitempty" yaml:"usage,omitempty"`
	DefValue            string `json:"defValue,omitempty" yaml:"defValue,omitempty"`
	NoOptDefVal         string `json:"noOptDefVal,omitempty" yaml:"noOptDefVal,omitempty"`
	Deprecated          string `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Hidden              bool   `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	ShorthandDeprecated string `json:"shorthandDeprecated,omitempty" yaml:"shorthandDeprecated,omitempty"`
	Inherited           bool   `json:"inherited,omitempty" yaml:"inherited,omitempty"`
	Persistent          bool   `json:"persistent,omitempty" yaml:"persistent,omitempty"`
	RawUsage            string `json:"rawUsage,omitempty" yaml:"rawUsage,omitempty"`
}

func processFlags(flagSet *pflag.FlagSet, fn func(f *Flag)) []Flag {
	results := make([]Flag, 0)
	if flagSet != nil && flagSet.HasFlags() {
		flagSet.VisitAll(func(cobraFlag *pflag.Flag) {
			current := Flag{
				Name:                cobraFlag.Name,
				Shorthand:           cobraFlag.Shorthand,
				Usage:               cobraFlag.Usage,
				DefValue:            cobraFlag.DefValue,
				NoOptDefVal:         cobraFlag.NoOptDefVal,
				Deprecated:          cobraFlag.Deprecated,
				Hidden:              cobraFlag.Hidden,
				ShorthandDeprecated: cobraFlag.ShorthandDeprecated,
				RawUsage:            internal.FlagUsage(cobraFlag),
			}

			if fn != nil {
				fn(&current)
			}

			results = append(results, current)
		})
	}
	return results
}

// NewCommandFromCobra creates a venom.Command from a cobra.Command
func NewCommandFromCobra(cmd *cobra.Command, options *Options) Command {
	var parentCommand *ParentCommand
	cobraParent := cmd.Parent()
	if cobraParent != nil {
		parentCommand = &ParentCommand{
			Name:  cobraParent.Name(),
			Short: cobraParent.Short,
		}
	}

	annotations := make(map[string]string)
	if cmd.Annotations != nil && len(cmd.Annotations) > 0 {
		for key, value := range cmd.Annotations {
			annotations[key] = value
		}
	}

	paths := make([]string, 0)

	command := Command{
		Name:          cmd.Name(),
		Usage:         cmd.UseLine(),
		Parent:        parentCommand,
		Aliases:       cmd.Aliases,
		SuggestFor:    cmd.SuggestFor,
		Short:         cmd.Short,
		Long:          cmd.Long,
		GroupID:       cmd.GroupID,
		ValidArgs:     cmd.ValidArgs,
		ArgAliases:    cmd.ArgAliases,
		Deprecated:    cmd.Deprecated,
		Annotations:   annotations,
		Version:       cmd.Version,
		Hidden:        cmd.Hidden,
		Runnable:      cmd.Runnable(),
		RawFlagUsages: strings.TrimSuffix(cmd.Flags().FlagUsages(), "\n"),
	}

	paths = append(paths, cmd.Name())
	if cobraParent != nil {
		for cobraParent != nil {
			paths = append(paths, cobraParent.Name())
			cobraParent = cobraParent.Parent()
		}
	}

	// reverse paths in-place
	for i, j := 0, len(paths)-1; i < j; i, j = i+1, j-1 {
		paths[i], paths[j] = paths[j], paths[i]
	}

	if parentCommand != nil {
		parentCommand.FullPath = strings.Join(paths[:len(paths)-1], " ")
	}
	command.FullPath = strings.Join(paths, " ")

	if cmd.HasExample() {
		examples := []string{cmd.Example}
		command.Examples = examples
	}

	subcommands := make([]Command, 0)
	for _, c := range cmd.Commands() {
		if c.Hidden && !options.showHiddenCommands {
			continue
		}
		sub := NewCommandFromCobra(c, options)
		subcommands = append(subcommands, sub)
	}

	cmd.Groups()

	command.Subcommands = subcommands

	command.LocalFlags = processFlags(cmd.LocalFlags(), nil)
	command.InheritedFlags = processFlags(cmd.InheritedFlags(), func(f *Flag) {
		f.Inherited = true
	})
	command.PersistentFlags = processFlags(cmd.PersistentFlags(), func(f *Flag) {
		f.Persistent = true
	})

	return command
}

func NewDocumentation(cmd *cobra.Command, options *Options) Documentation {
	venomCommand := NewCommandFromCobra(cmd, options)
	doc := Documentation{
		RootCommand: venomCommand,
		options:     options,
	}
	if !cmd.DisableAutoGenTag {
		doc.AutoGenerationTag = "Auto-generated by jimschubert/venom"
	}
	doc.init()
	return doc
}
