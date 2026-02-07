package main

import "testing"

func TestAdaptiveThreshold(t *testing.T) {
	tests := []struct {
		name   string
		bufLen int
		expect int
	}{
		{
			name:   "empty buffer",
			bufLen: 0,
			expect: 32,
		},
		{
			name:   "small buffer",
			bufLen: 100,
			expect: 72,
		},
		{
			name:   "medium buffer",
			bufLen: 10000,
			expect: 432,
		},
		{
			name:   "large buffer",
			bufLen: 100000,
			expect: 1296,
		},
		{
			name:   "very large buffer exceeds max",
			bufLen: 1000000,
			expect: 4032,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adaptiveThreshold(tt.bufLen)

			t.Log(got)

			if got != tt.expect {
				t.Errorf("adaptiveThreshold(%d) = %d, want %d", tt.bufLen, got, tt.expect)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name string
		x    int
		min  int
		max  int
		want int
	}{
		{"below min", 5, 10, 20, 10},
		{"above max", 25, 10, 20, 20},
		{"within range", 15, 10, 20, 15},
		{"equals min", 10, 10, 20, 10},
		{"equals max", 20, 10, 20, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clamp(tt.x, tt.min, tt.max); got != tt.want {
				t.Errorf("clamp(%d, %d, %d) = %d, want %d", tt.x, tt.min, tt.max, got, tt.want)
			}
		})
	}
}
