package functions

import (
	"testing"
)

func TestCheckLuhn(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"This one is correct Luhn",
			args{
				"4561261212345467",
			},
			true,
		},
		{
			"This one is NOT correct Luhn",
			args{
				"4561261212345464",
			},
			false,
		},
		{
			"Empty is NOT correct Luhn",
			args{
				"",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := CheckLuhn(tt.args.payload); got != tt.want {
					t.Errorf("CheckLuhn() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
