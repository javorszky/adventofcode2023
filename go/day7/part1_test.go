package day7

import (
	"reflect"
	"testing"
)

func Test_sortHands(t *testing.T) {
	type args struct {
		hands []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "sorts list of hands that all differ on the last char",
			args: args{
				hands: []string{
					"AAAAQ",
					"AAAA3",
					"AAAA9",
					"AAAAT",
					"AAAAA",
					"AAAA2",
				},
			},
			want: []string{
				"AAAA2",
				"AAAA3",
				"AAAA9",
				"AAAAT",
				"AAAAQ",
				"AAAAA",
			},
		},
		{
			name: "sorts some weird things 1",
			args: args{
				hands: []string{
					"22272",
					"QQQQ2",
					"22223",
					"42444",
					"63333",
				},
			},
			want: []string{
				"22223",
				"22272",
				"42444",
				"63333",
				"QQQQ2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortHands(tt.args.hands); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortHands() = %v, want %v", got, tt.want)
			}
		})
	}
}
