package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hashF := func(k []byte) uint32 {
		i, _ := strconv.Atoi(string(k))
		return uint32(i)
	}
	hash := NewMap(3, hashF)
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("asking for %s, should have yielded %s", k, v)
		}
	}
}
