package log

import (
	"testing"
)

func TestNewSyncLogger(t *testing.T) {
	s := NewSyncLogger()
	s.Close()
}
