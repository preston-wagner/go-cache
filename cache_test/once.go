package cache_test

import (
	"strconv"
	"testing"

	"github.com/preston-wagner/unicycle/sets"
)

func ItoaOnce(t *testing.T) func(int) string {
	keys := sets.Set[int]{}
	return func(key int) string {
		if keys.Has(key) {
			t.Error("wrapped function was called with the same key multiple times!")
		}
		keys.Add(key)
		return strconv.Itoa(key)
	}
}

func AtoiOnce(t *testing.T) func(string) (int, error) {
	keys := sets.Set[string]{}
	return func(key string) (int, error) {
		if keys.Has(key) {
			t.Error("wrapped function was called with the same key multiple times!")
		}
		keys.Add(key)
		return strconv.Atoi(key)
	}
}
