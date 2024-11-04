package converter

import "testing"

func Test_scoreToResult(t *testing.T) {
	type args struct {
		score map[string]float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				score: map[string]float64{
					"g":  5,
					"g5": 5,
					"3":  5,
				},
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScoreToResult(tt.args.score, -1); got != tt.want {
				t.Errorf("scoreToResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
