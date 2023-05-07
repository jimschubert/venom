package venom

import (
	"github.com/go-test/deep"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func withChildren(cmd *cobra.Command, children ...*cobra.Command) *cobra.Command {
	for _, child := range children {
		cmd.AddCommand(child)
	}
	return cmd
}

func withDefaults(command *Command) Command {
	if command.Subcommands == nil {
		command.Subcommands = []Command{}
	}
	if command.LocalFlags == nil {
		command.LocalFlags = []Flag{}
	}
	if command.InheritedFlags == nil {
		command.InheritedFlags = []Flag{}
	}
	if command.PersistentFlags == nil {
		command.PersistentFlags = []Flag{}
	}
	if command.Annotations == nil {
		command.Annotations = map[string]string{}
	}
	return *command
}

func TestNewCommandFromCobra(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "Simple command, no children",
			args: args{cmd: &cobra.Command{
				Use:   "pinky",
				Short: "p",
			}},
			want: withDefaults(&Command{
				Name:  "pinky",
				Usage: "pinky",
				Short: "p",
			}),
		},
		{
			name: "Simple command, with children",
			args: args{cmd: withChildren(&cobra.Command{
				Use:   "pinky",
				Short: "p",
			}, withChildren(&cobra.Command{
				Use: "and",
			}, &cobra.Command{
				Use: "brain",
			},
			))},
			want: withDefaults(&Command{
				Name:  "pinky",
				Usage: "pinky",
				Short: "p",
				Subcommands: []Command{
					withDefaults(&Command{
						Name:   "and",
						Usage:  "pinky and",
						Parent: &ParentCommand{Name: "pinky"},
						Subcommands: []Command{
							withDefaults(&Command{
								Name:   "brain",
								Usage:  "pinky and brain",
								Parent: &ParentCommand{Name: "and"},
							}),
						},
					}),
				},
			}),
		},
		{
			name: "Full command",
			args: args{cmd: &cobra.Command{
				Use:        "sample",
				Short:      "s",
				Long:       "sample test",
				Deprecated: "(deprecated)",
				Version:    "1.0.0",
				Hidden:     true,
				Aliases:    []string{"example", "testing"},
				SuggestFor: []string{"a", "b"},
				Run: func(cmd *cobra.Command, args []string) {

				},
				Annotations: map[string]string{
					"First": "Second",
				},
			}},
			want: withDefaults(&Command{
				Name:       "sample",
				Short:      "s",
				Usage:      "sample",
				Long:       "sample test",
				Deprecated: "(deprecated)",
				Version:    "1.0.0",
				Hidden:     true,
				Aliases:    []string{"example", "testing"},
				SuggestFor: []string{"a", "b"},
				Runnable:   true,
				Annotations: map[string]string{
					"First": "Second",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommandFromCobra(tt.args.cmd)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewCommandFromCobra():\n%v", strings.Join(diff, "\t\n"))
			}
		})
	}
}
