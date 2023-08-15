package cache

import (
	"testing"
)

func TestMemUsage(t *testing.T) {
	memUsage := MemUsage()
	if memUsage <= 0 {
		t.Fatal("MemUsage() should not return <= 0")
	}
	if memUsage > 1 {
		t.Fatal("MemUsage() should not return > 1")
	}
}
