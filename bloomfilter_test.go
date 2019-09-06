package bloomfilter

import (
	"strconv"
	"testing"
)

func BenchmarkName(b *testing.B) {
	bf := NewBloomFilter(100000, 7)
	for i := 0; i < b.N; i++ {
		data := []byte("hello" + strconv.Itoa(i))
		bf.Add(data)
		if ok, err := bf.Test(data); err != nil {
			b.Error(err)
		} else if !ok {
			b.Error("bug")
		}
	}
}
