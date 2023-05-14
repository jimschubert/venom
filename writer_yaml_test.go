package venom

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/jimschubert/venom/internal"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestYamlWrite(t *testing.T) {
	type args struct {
		doc     Documentation
		options TemplateOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "writes yaml doc",
			args: args{
				options: NewOptions().TemplateOptions(),
				doc: Documentation{
					AutoGenerationTag: "generated by: Simple doc yaml",
					RootCommand: Command{
						Name:  "simple",
						Short: "s",
						Long:  "simple doc yaml",
						LocalFlags: []Flag{
							{
								Name:        "testing",
								Shorthand:   "t",
								DefValue:    "true",
								NoOptDefVal: "true",
								Hidden:      false,
								RawUsage:    "--testing",
							},
						},
						Runnable: true,
						Subcommands: []Command{
							{
								Name:     "command",
								Short:    "c",
								Long:     "simple yaml subcommand",
								Hidden:   true,
								Runnable: false,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outDir := t.TempDir()
			w := writerYaml{
				options: tt.args.options,
			}
			if err := w.Write(outDir, tt.args.doc); (err != nil) != tt.wantErr {
				t.Fatalf("writerYaml() error = %v, wantErr %v", err, tt.wantErr)
			}

			clean := internal.CleanPath(tt.args.doc.RootCommand.Name)
			b, err := os.ReadFile(filepath.Join(outDir, clean, fmt.Sprintf("%s.yml", clean)))
			if err != nil {
				t.Fatalf("writerYaml() unable to read file at expected path")
			}

			writtenDoc := Documentation{}
			err = yaml.Unmarshal(b, &writtenDoc)
			if err != nil {
				t.Fatalf("writerYaml() unable to unmarshal the generated yml")
			}

			if diff := deep.Equal(tt.args.doc, writtenDoc); diff != nil {
				t.Fatalf("writerYaml():\n%v", strings.Join(diff, "\t\n"))
			}
		})
	}
}
