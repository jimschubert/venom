package venom

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
	"testing"
)

func TestOptions_validate(t *testing.T) {
	defaultTemplateOptions := TemplateOptions{
		JsonMarshaler: json.Marshal,
		YamlMarshaler: yaml.Marshal,
		Logger:        log.Default(),
	}
	type fields struct {
		commandName   string
		formats       Formats
		logger        *log.Logger
		jsonMarshaler MarshalFn
		yamlMarshaler MarshalFn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validate fails for empty commandName",
			fields: fields{
				commandName: "",
				formats:     Markdown,
			},
			wantErr: true,
		},
		{
			name: "validate fails for empty formats",
			fields: fields{
				commandName: "asdf",
			},
			wantErr: true,
		},
		{
			name: "validate fails for invalid formats",
			fields: fields{
				commandName: "asdf",
				formats:     Formats(1<<7) | Formats(1<<6),
			},
			wantErr: true,
		},
		{
			name: "validate succeeds for valid inputs",
			fields: fields{
				commandName: "asdf",
				formats:     Markdown | Yaml,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := TemplateOptions{
				Logger:        defaultTemplateOptions.Logger,
				YamlMarshaler: defaultTemplateOptions.YamlMarshaler,
				JsonMarshaler: defaultTemplateOptions.JsonMarshaler,
			}

			if tt.fields.logger != nil {
				opts.Logger = tt.fields.logger
			}

			if tt.fields.jsonMarshaler != nil {
				opts.JsonMarshaler = tt.fields.jsonMarshaler
			}

			if tt.fields.yamlMarshaler != nil {
				opts.YamlMarshaler = tt.fields.yamlMarshaler
			}

			o := &Options{
				commandName:     tt.fields.commandName,
				formats:         tt.fields.formats,
				templateOptions: &opts,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
