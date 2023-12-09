package day5

import (
	"reflect"
	"testing"
)

func Test_overlaps(t *testing.T) {
	sourceSections := []section{
		{
			sourceStart:      4,
			sourceEnd:        60,
			destinationStart: 4,
			destinationEnd:   60,
			delta:            0,
		},
		{
			sourceStart:      98,
			sourceEnd:        142,
			destinationStart: 98,
			destinationEnd:   142,
			delta:            0,
		},
	}

	sourceSectionsWithDeltas := []section{
		{
			sourceStart:      0,
			sourceEnd:        56,
			destinationStart: 4,
			destinationEnd:   60,
			delta:            4,
		},
		{
			sourceStart:      128,
			sourceEnd:        172,
			destinationStart: 98,
			destinationEnd:   142,
			delta:            -30,
		},
	}

	type args struct {
		listA []section
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
				listA: sourceSections,
				listB: []section{
					{
						sourceStart:      0,
						sourceEnd:        9,
						destinationStart: 4,
						destinationEnd:   13,
						delta:            4,
					},
					{
						sourceStart:      53,
						sourceEnd:        101,
						destinationStart: 70,
						destinationEnd:   118,
						delta:            17,
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
				},
				{
					sourceStart:      10,
					sourceEnd:        52,
					destinationStart: 10,
					destinationEnd:   52,
					delta:            0,
				},
				{
					sourceStart:      53,
					sourceEnd:        60,
					destinationStart: 70,
					destinationEnd:   77,
					delta:            17,
				},
				{
					sourceStart:      98,
					sourceEnd:        101,
					destinationStart: 115,
					destinationEnd:   118,
					delta:            17,
				},
				{
					sourceStart:      102,
					sourceEnd:        142,
					destinationStart: 102,
					destinationEnd:   142,
					delta:            0,
				},
			},
		},
		{
			name: "finds overlaps with sources that have deltas",
			args: args{
				listA: sourceSectionsWithDeltas,
				listB: []section{
					{
						sourceStart:      0,
						sourceEnd:        9,
						destinationStart: 4,
						destinationEnd:   13,
						delta:            4,
					},
					{
						sourceStart:      53,
						sourceEnd:        101,
						destinationStart: 70,
						destinationEnd:   118,
						delta:            17,
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
				},
				{
					sourceStart:      10,
					sourceEnd:        52,
					destinationStart: 10,
					destinationEnd:   52,
					delta:            0,
				},
				{
					sourceStart:      53,
					sourceEnd:        60,
					destinationStart: 70,
					destinationEnd:   77,
					delta:            17,
				},
				{
					sourceStart:      98,
					sourceEnd:        101,
					destinationStart: 115,
					destinationEnd:   118,
					delta:            17,
				},
				{
					sourceStart:      102,
					sourceEnd:        142,
					destinationStart: 102,
					destinationEnd:   142,
					delta:            0,
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

func Test_doesSectionCoverThisValue(t *testing.T) {
	type args struct {
		value    int
		sections []section
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 int
		want2 int
	}{
		{
			name: "does not cover it, value smaller than first section, should be true and first sourcestart",
			args: args{
				value: 0,
				sections: []section{
					{
						sourceStart:      5,
						sourceEnd:        10,
						destinationStart: 6,
						destinationEnd:   11,
						delta:            1,
					},
				},
			},
			want:  false,
			want1: 5,
			want2: 0,
		},
		{
			name: "covers it, value = first source start",
			args: args{
				value: 5,
				sections: []section{
					{
						sourceStart:      5,
						sourceEnd:        10,
						destinationStart: 6,
						destinationEnd:   11,
						delta:            1,
					},
				},
			},
			want:  true,
			want1: 10,
			want2: 1,
		},
		{
			name: "covers it, value = first source end",
			args: args{
				value: 10,
				sections: []section{
					{
						sourceStart:      5,
						sourceEnd:        10,
						destinationStart: 6,
						destinationEnd:   11,
						delta:            1,
					},
				},
			},
			want:  true,
			want1: 10,
			want2: 1,
		},
		{
			name: "does not cover it it, value > first source end",
			args: args{
				value: 11,
				sections: []section{
					{
						sourceStart:      5,
						sourceEnd:        10,
						destinationStart: 6,
						destinationEnd:   11,
						delta:            1,
					},
					{
						sourceStart:      18,
						sourceEnd:        20,
						destinationStart: 6,
						destinationEnd:   8,
						delta:            -12,
					},
				},
			},
			want:  false,
			want1: 18,
			want2: 0,
		},
		{
			name: "covers it, is in the second section",
			args: args{
				value: 19,
				sections: []section{
					{
						sourceStart:      5,
						sourceEnd:        10,
						destinationStart: 6,
						destinationEnd:   11,
						delta:            1,
					},
					{
						sourceStart:      18,
						sourceEnd:        20,
						destinationStart: 6,
						destinationEnd:   8,
						delta:            -12,
					},
				},
			},
			want:  true,
			want1: 20,
			want2: -12,
		},
		{
			name: "does not cover it, is past the sections",
			args: args{
				value: 24,
				sections: []section{
					{
						sourceStart:      5,
						sourceEnd:        10,
						destinationStart: 6,
						destinationEnd:   11,
						delta:            1,
					},
					{
						sourceStart:      18,
						sourceEnd:        20,
						destinationStart: 6,
						destinationEnd:   8,
						delta:            -12,
					},
				},
			},
			want:  false,
			want1: 0,
			want2: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := doesSectionCoverThisValue(tt.args.value, tt.args.sections)

			if got != tt.want {
				t.Errorf("doesSectionCoverThisValue() got = %v, want %v", got, tt.want)
			}

			if got1 != tt.want1 {
				t.Errorf("doesSectionCoverThisValue() got1 = %v, want %v", got1, tt.want1)
			}

			if got2 != tt.want2 {
				t.Errorf("doesSectionCoverThisValue() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
