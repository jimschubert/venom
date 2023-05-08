package venom

import "testing"

func TestCleanPath(t *testing.T) {
	type args struct {
		input   string
		replace []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no replacement for valid string",
			args: args{input: "This.is.valid"},
			want: "This.is.valid",
		},
		{
			name: "replacement for unexpected characters",
			args: args{input: "first/second/third/last"},
			want: "first_second_third_last",
		},
		{
			name: "no replacement for valid string when passed custom replace override",
			args: args{input: "This.is.valid", replace: []string{"|"}},
			want: "This.is.valid",
		},
		{
			name: "replacement for unexpected characters when passed custom replace override",
			args: args{input: "first/second/third/last", replace: []string{"-"}},
			want: "first-second-third-last",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanPath(tt.args.input, tt.args.replace...); got != tt.want {
				t.Errorf("CleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
