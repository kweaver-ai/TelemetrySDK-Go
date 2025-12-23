package examples

import "testing"

func TestExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test原始的业务系统",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTTPExample()
		})
	}
}
