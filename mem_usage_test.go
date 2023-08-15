package cache

import (
	"log"
	"testing"
)

func TestMemUsage(t *testing.T) {
	memUsage := MemUsage()
	log.Println(memUsage)
	if memUsage == 0 {
		t.Fatal("MemUsage() should not return 0")
	}
}
