package diff

import (
	"reflect"
	"testing"

	"github.com/hexops/gotextdiff"
)

func TestSubtract(t *testing.T) {
	type args struct {
		a gotextdiff.Unified
		b gotextdiff.Unified
	}
	tests := []struct {
		name string
		args args
		want gotextdiff.Unified
	}{
		{
			name: "Equal diffs",
			args: args{
				a: gotextdiff.Unified{
					Hunks: []*gotextdiff.Hunk{
						{
							FromLine: 1,
							ToLine:   3,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal1\n"},
								{Kind: gotextdiff.Insert, Content: "insert1\n"},
								{Kind: gotextdiff.Equal, Content: "equal2\n"},
								{Kind: gotextdiff.Delete, Content: "delete1\n"},
							},
						},
					},
				},
				b: gotextdiff.Unified{
					Hunks: []*gotextdiff.Hunk{
						{
							FromLine: 1,
							ToLine:   3,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal1\n"},
								{Kind: gotextdiff.Insert, Content: "insert1\n"},
								{Kind: gotextdiff.Equal, Content: "equal2\n"},
								{Kind: gotextdiff.Delete, Content: "delete1\n"},
							},
						},
					},
				},
			},
			want: gotextdiff.Unified{
				Hunks: []*gotextdiff.Hunk{},
			},
		},
		{
			name: "Add hunk at the beginning",
			args: args{
				a: gotextdiff.Unified{
					Hunks: []*gotextdiff.Hunk{
						{
							FromLine: 1,
							ToLine:   3,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal1\n"},
								{Kind: gotextdiff.Insert, Content: "insert1\n"},
								{Kind: gotextdiff.Equal, Content: "equal2\n"},
								{Kind: gotextdiff.Delete, Content: "delete1\n"},
							},
						},
						{
							FromLine: 4,
							ToLine:   6,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal3\n"},
								{Kind: gotextdiff.Insert, Content: "insert2\n"},
								{Kind: gotextdiff.Equal, Content: "equal4\n"},
								{Kind: gotextdiff.Delete, Content: "delete2\n"},
							},
						},
					},
				},
				b: gotextdiff.Unified{
					Hunks: []*gotextdiff.Hunk{
						{
							FromLine: 4,
							ToLine:   6,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal3\n"},
								{Kind: gotextdiff.Insert, Content: "insert2\n"},
								{Kind: gotextdiff.Equal, Content: "equal4\n"},
								{Kind: gotextdiff.Delete, Content: "delete2\n"},
							},
						},
					},
				},
			},
			want: gotextdiff.Unified{
				Hunks: []*gotextdiff.Hunk{
					{
						FromLine: 1,
						ToLine:   3,
						Lines: []gotextdiff.Line{
							{Kind: gotextdiff.Equal, Content: "equal1\n"},
							{Kind: gotextdiff.Insert, Content: "insert1\n"},
							{Kind: gotextdiff.Equal, Content: "equal2\n"},
							{Kind: gotextdiff.Delete, Content: "delete1\n"},
						},
					},
				},
			},
		},
		{
			name: "Add hunk at the end",
			args: args{
				a: gotextdiff.Unified{
					Hunks: []*gotextdiff.Hunk{
						{
							FromLine: 1,
							ToLine:   3,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal1"},
								{Kind: gotextdiff.Insert, Content: "insert1"},
								{Kind: gotextdiff.Equal, Content: "equal2"},
								{Kind: gotextdiff.Delete, Content: "delete1"},
							},
						},
						{
							FromLine: 4,
							ToLine:   6,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal3"},
								{Kind: gotextdiff.Insert, Content: "insert2"},
								{Kind: gotextdiff.Equal, Content: "equal4"},
								{Kind: gotextdiff.Delete, Content: "delete2"},
							},
						},
					},
				},
				b: gotextdiff.Unified{
					Hunks: []*gotextdiff.Hunk{
						{
							FromLine: 1,
							ToLine:   3,
							Lines: []gotextdiff.Line{
								{Kind: gotextdiff.Equal, Content: "equal1"},
								{Kind: gotextdiff.Insert, Content: "insert1"},
								{Kind: gotextdiff.Equal, Content: "equal2"},
								{Kind: gotextdiff.Delete, Content: "delete1"},
							},
						},
					},
				},
			},
			want: gotextdiff.Unified{
				Hunks: []*gotextdiff.Hunk{
					{
						FromLine: 4,
						ToLine:   6,
						Lines: []gotextdiff.Line{
							{Kind: gotextdiff.Equal, Content: "equal3"},
							{Kind: gotextdiff.Insert, Content: "insert2"},
							{Kind: gotextdiff.Equal, Content: "equal4"},
							{Kind: gotextdiff.Delete, Content: "delete2"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := subtract(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}
