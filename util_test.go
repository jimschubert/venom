package venom

import "testing"

func Test_hangingIndent(t *testing.T) {
	type args struct {
		input     string
		hangWidth int
		maxWidth  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single line below max width",
			args: args{
				input:     "single line below max width",
				hangWidth: 15,
				maxWidth:  30,
			},
			want: "single line below max width",
		},
		{
			name: "multi-line short",
			args: args{
				input:     "The quick brown fox jumped over the lazy dog",
				hangWidth: 6,
				maxWidth:  18,
			},
			want: "The quick brown\n      fox jumped\n      over the\n      lazy dog",
		},
		{
			name: "multi-line long, command wrapping",
			args: args{
				input:     "-s, --supported         Set to try and demonstrate some other boolean flag, but this one with is written with a short option (default true)",
				hangWidth: 24,
				maxWidth:  120,
			},
			want: "-s, --supported         Set to try and demonstrate some other boolean flag, but this one with is written with a short\n                        option (default true)",
		},
		{
			name: "multi-line long",
			args: args{
				input:     "      --config string   config file which we would like to do something with as an example of how this will display when way too long (default is $HOME/.sample.yaml)",
				hangWidth: 24,
				maxWidth:  60,
			},
			want: "      --config string   config file which we would like to\n                        do something with as an example of\n                        how this will display when way too\n                        long (default is\n                        $HOME/.sample.yaml)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hangingIndent(tt.args.input, tt.args.hangWidth, tt.args.maxWidth); got != tt.want {
				t.Errorf("hangingIndent() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
