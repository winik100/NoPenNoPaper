package main

import (
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestHalfAndFifth(t *testing.T) {

	tests := []struct {
		name  string
		f     func(int) int
		value int
		want  int
	}{
		{
			name:  "Half - Even",
			f:     half,
			value: 10,
			want:  5,
		},
		{
			name:  "Half - Odd",
			f:     half,
			value: 13,
			want:  6,
		},
		{
			name:  "Half - Zero",
			f:     half,
			value: 0,
			want:  0,
		},
		{
			name:  "Fifth - Divisible by Five",
			f:     fifth,
			value: 35,
			want:  7,
		},
		{
			name:  "Fifth - Not Divisible by Five",
			f:     fifth,
			value: 47,
			want:  9,
		},
		{
			name:  "Fifth - Zero",
			f:     fifth,
			value: 0,
			want:  0,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			res := testCase.f(testCase.value)

			testHelpers.Equal(t, res, testCase.want)
		})
	}
}
