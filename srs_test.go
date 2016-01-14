package mian

import (
	"os"
	"testing"
)

func TestMian(m *testing.M) {
	os.Exit(m.Run())
}
