package cache

import (
	"log"

	"github.com/pbnjay/memory"
)

// returns a value between 0 and 1 representing the current percentage of memory used
func MemUsage() float64 {
	totalMemory := memory.TotalMemory()
	log.Println(totalMemory)
	freeMemory := memory.FreeMemory()
	log.Println(freeMemory)
	return 1 - (float64(freeMemory) / float64(totalMemory))
}
