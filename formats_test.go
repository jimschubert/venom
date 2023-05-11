package venom

import (
	"reflect"
	"testing"
)

func ptrFormat(f Formats) *Formats {
	return &f
}

func TestFormats_IsSet(t *testing.T) {
	type args struct {
		format Formats
	}
	tests := []struct {
		name string
		f    Formats
		args args
		want bool
	}{
		{
			name: "Single set test (positive)",
			f:    Markdown,
			args: args{format: Markdown},
			want: true,
		},
		{
			name: "Single set test (negative)",
			f:    Markdown,
			args: args{format: Yaml},
			want: false,
		},
		{
			name: "Multiple set test (positive)",
			f:    Markdown | Yaml | Json,
			args: args{format: Markdown | Json},
			want: true,
		},
		{
			name: "Multiple set test (negative)",
			f:    Markdown | Json,
			args: args{format: Yaml | Man},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsSet(tt.args.format); got != tt.want {
				t.Errorf("IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormats_IsValid(t *testing.T) {
	tests := []struct {
		name string
		f    Formats
		want bool
	}{
		{
			name: "IsValid single (positive)",
			f:    Man,
			want: true,
		},
		{
			name: "IsValid multiple (positive)",
			f:    Man | Yaml,
			want: true,
		},
		{
			name: "IsValid (negative)",
			f:    Formats(0),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormats_Set(t *testing.T) {
	type args struct {
		format Formats
	}
	tests := []struct {
		name string
		f    Formats
		args args
		want *Formats
	}{
		{
			name: "Set from empty",
			f:    Formats(0),
			args: args{format: Man | Yaml},
			want: ptrFormat(Man | Yaml),
		},
		{
			name: "Set from existing",
			f:    Markdown,
			args: args{format: Man | Yaml},
			want: ptrFormat(Markdown | Man | Yaml),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Set(tt.args.format); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormats_Unset(t *testing.T) {
	type args struct {
		format Formats
	}
	tests := []struct {
		name string
		f    Formats
		args args
		want *Formats
	}{
		{
			name: "Unset from empty",
			f:    Formats(0),
			args: args{format: Man | Yaml},
			want: ptrFormat(Formats(0)),
		},
		{
			name: "Unset from existing",
			f:    Markdown | Man | Yaml,
			args: args{format: Markdown | Yaml},
			want: ptrFormat(Man),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Unset(tt.args.format); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormats_defined(t *testing.T) {
	tests := []struct {
		name string
		f    Formats
		want []Formats
	}{
		{
			name: "one",
			f:    Yaml,
			want: []Formats{Yaml},
		},
		{
			name: "two",
			f:    Yaml | Markdown | Json,
			want: []Formats{Yaml, Json, Markdown},
		},
		{
			name: "multiple",
			f:    Markdown | Yaml | Man | Json | ReST,
			want: []Formats{Yaml, Json, Markdown, Man, ReST},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.defined(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defined() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormats_String(t *testing.T) {
	tests := []struct {
		name string
		i    Formats
		want string
	}{
		{
			name: "Yaml",
			i:    Yaml,
			want: "Yaml",
		},
		{
			name: "Json",
			i:    Json,
			want: "Json",
		},
		{
			name: "Markdown",
			i:    Markdown,
			want: "Markdown",
		},
		{
			name: "Man",
			i:    Man,
			want: "Man",
		},
		{
			name: "ReST",
			i:    ReST,
			want: "ReST",
		},
		{
			name: "multiple",
			i:    Yaml | Markdown | Json,
			want: "Yaml|Json|Markdown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
