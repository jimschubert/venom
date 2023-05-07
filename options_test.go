package venom

import (
	"log"
	"testing"
)

func TestOptions_validate(t *testing.T) {
	type fields struct {
		commandName string
		formats     Formats
		logger      *log.Logger
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
			o := &Options{
				commandName: tt.fields.commandName,
				formats:     tt.fields.formats,
				logger:      tt.fields.logger,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
