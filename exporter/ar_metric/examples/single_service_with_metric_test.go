package examples

import (
	"testing"
)

func TestFileExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"File",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FileExample()
		})
	}
}

func TestConsoleExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Console",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConsoleExample()
		})
	}
}

func TestStdoutExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Stdout",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StdoutExample()
		})
	}
}

func TestHTTPExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test发送到AnyRobot Metric数据接收器",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTTPExample()
		})
	}
}

func TestWithAllExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test调用全部接口的示例",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithAllExample()
		})
	}
}
