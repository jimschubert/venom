package venom

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Documentation represents the "top-level" of documentation to be passed to a template
type Documentation struct {
	GenerationDate    string  `yaml:"generationDate,omitempty" json:"generationDate,omitempty"`
	AutoGenerationTag string  `yaml:"autoGenerationTag,omitempty" json:"autoGenerationTag,omitempty"`
	RootCommand       Command `yaml:"rootCommand,omitempty" json:"rootCommand,omitempty"`
}

// ParentCommand provides the name of a command's parent
type ParentCommand struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
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
}

// Flag is a representation of pflag.Flag
type Flag struct {
	Name                string
	Shorthand           string
	Usage               string
	DefValue            string
	NoOptDefVal         string
	Deprecated          string
	Hidden              bool
	ShorthandDeprecated string
	Inherited           bool
	Persistent          bool
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
				Inherited:           true,
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
func NewCommandFromCobra(cmd *cobra.Command) Command {
	var parentCommand *ParentCommand
	cobraParent := cmd.Parent()
	if cobraParent != nil {
		parentCommand = &ParentCommand{Name: cobraParent.Name()}
	}

	annotations := make(map[string]string)
	if cmd.Annotations != nil && len(cmd.Annotations) > 0 {
		for key, value := range cmd.Annotations {
			annotations[key] = value
		}
	}

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
		RawFlagUsages: cmd.Flags().FlagUsages(),
	}

	if cmd.HasExample() {
		examples := []string{cmd.Example}
		command.Examples = examples
	}

	cmdLen := len(cmd.Commands())
	subcommands := make([]Command, cmdLen)
	for i, c := range cmd.Commands() {
		sub := NewCommandFromCobra(c)
		subcommands[i] = sub
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
