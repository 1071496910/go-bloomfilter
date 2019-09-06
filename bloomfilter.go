package bloomfilter

import (
	"github.com/1071496910/go-bitmap"
	"github.com/spaolacci/murmur3"
)

type BloomFilter interface {
	Add([]byte) error
	Test([]byte) (bool, error)
}

type bloomFilter struct {
	stroage bitmap.Bitmap
	m       int
	k       int
}

// baseHashes returns the four hash values of data that are used to create k
// hashes
func baseHashes(data []byte) [4]uint64 {
	a1 := []byte{1} // to grab another bit of data
	hasher := murmur3.New128()
	hasher.Write(data) // #nosec
	v1, v2 := hasher.Sum128()
	hasher.Write(a1) // #nosec
	v3, v4 := hasher.Sum128()
	return [4]uint64{
		v1, v2, v3, v4,
	}
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

func (b *bloomFilter) Add(data []byte) error {
	h := baseHashes(data)
	for i := uint(0); i < uint(b.k); i++ {
		b.stroage.Set(b.location(h, i))
	}
	return nil
}

func (b *bloomFilter) Test(data []byte) (bool, error) {
	h := baseHashes(data)
	for i := 0; i < b.k; i++ {
		if ok, err := b.stroage.Check(b.location(h, uint(i))); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (b *bloomFilter) location(h [4]uint64, i uint) uint {
	return uint(location(h, i) % uint64(b.m))
}
func NewBloomFilter(m int, k int) *bloomFilter {
	return &bloomFilter{
		stroage: bitmap.NewBitmap(m),
		m:       m,
		k:       k,
	}
}
