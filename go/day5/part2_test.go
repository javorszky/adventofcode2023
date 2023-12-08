package day5

import (
	"reflect"
	"testing"
)

func Test_overlaps(t *testing.T) {
	type args struct {
		listA []int
		listB []section
	}
	tests := []struct {
		name string
		args args
		want []section
	}{
		{
			name: "finds overlap",
			args: args{
				listA: []int{4, 56, 98, 44},
				listB: []section{
					{
						sourceStart:      0,
						sourceEnd:        9,
						destinationStart: 4,
						destinationEnd:   13,
						delta:            4,
						transformFunc: func(i int) int {
							return i + 4
						},
					},
					{
						sourceStart:      53,
						sourceEnd:        101,
						destinationStart: 70,
						destinationEnd:   118,
						delta:            17,
						transformFunc: func(i int) int {
							return i + 17
						},
					},
				},
			},
			want: []section{
				{
					sourceStart:      4,
					sourceEnd:        9,
					destinationStart: 8,
					destinationEnd:   13,
					delta:            4,
					transformFunc: func(i int) int {
						return i + 4
					},
				},
				{
					sourceStart:      5,
					sourceEnd:        52,
					destinationStart: 5,
					destinationEnd:   52,
					delta:            0,
					transformFunc: func(i int) int {
						return i
					},
				},
				{
					sourceStart:      53,
					sourceEnd:        60,
					destinationStart: 70,
					destinationEnd:   87,
					delta:            17,
					transformFunc: func(i int) int {
						return i + 17
					},
				},
				{
					sourceStart:      98,
					sourceEnd:        101,
					destinationStart: 115,
					destinationEnd:   118,
					delta:            17,
					transformFunc: func(i int) int {
						return i + 17
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := overlaps(tt.args.listA, tt.args.listB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("overlaps() = %v, want %v", got, tt.want)
			}
		})
	}
}
