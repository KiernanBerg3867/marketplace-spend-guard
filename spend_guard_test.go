package main

import "testing"

func TestAuthorize(t *testing.T) {
	tests := []struct {
		name    string
		order   Order
		used    float64
		allowed bool
	}{
		{"room remains", Order{Budget: 10}, 2, true},
		{"cap reached", Order{Budget: 10}, 10, false},
		{"invalid budget", Order{Budget: 0}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := authorize(tt.order, tt.used).Allowed; got != tt.allowed {
				t.Fatalf("allowed=%v, want %v", got, tt.allowed)
			}
		})
	}
}
