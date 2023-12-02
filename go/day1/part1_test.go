package day1

import "testing"

func Test_replaceSpelled(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "replaces writtens",
			args: args{line: "one1two2three3four4five5six6seven7eight8nine9"},
			want: "112233445566778899",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceSpelled(tt.args.line); got != tt.want {
				t.Errorf("replaceSpelled() = %v, want %v", got, tt.want)
			}
		})
	}
}
