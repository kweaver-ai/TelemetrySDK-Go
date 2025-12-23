package eventsdk

import "testing"

func TestLevelSelf(t *testing.T) {
	tests := []struct {
		name string
		l    level
		want string
	}{
		{
			"",
			ERROR,
			"ERROR",
		},
		{
			"",
			WARN,
			"WARN",
		},
		{
			"",
			INFO,
			"INFO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Self(); got != tt.want {
				t.Errorf("Self() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelPrivate(t *testing.T) {
	tests := []struct {
		name string
		l    level
	}{
		{
			"",
			level(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.private()
		})
	}
}

func TestNewLevel(t *testing.T) {
	type args struct {
		l string
	}
	tests := []struct {
		name string
		args args
		want level
	}{
		{
			"",
			args{"WARN"},
			level("WARN"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLevel(tt.args.l); got != tt.want {
				t.Errorf("newLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelValid(t *testing.T) {
	tests := []struct {
		name string
		l    level
		want bool
	}{
		{
			"",
			level(""),
			false,
		},
		{
			"",
			level("ERROR"),
			true,
		},
		{
			"",
			level("WARN"),
			true,
		},
		{
			"",
			level("INFO"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
