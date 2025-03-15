package bloomfilter

import (
	"fmt"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)
	bf.Add([]byte("test"))

	if !bf.Exists([]byte("test")) {
		t.Errorf("Expected 'test' to exist in Bloom Filter")
	}

	if bf.Exists([]byte("random")) {
		t.Errorf("Expected 'random' to NOT exist in Bloom Filter")
	}
}

func TestFalsePositiveRate(t *testing.T) {
	bf := NewBloomFilter(2000, 0.005)

	// Insert 1000 elements
	for i := 0; i < 1000; i++ {
		bf.Add([]byte(fmt.Sprintf("item-%d", i)))
	}

	// Check 1000 elements not inserted
	falsePositives := 0
	for i := 1000; i < 2000; i++ {
		if bf.Exists([]byte(fmt.Sprintf("item-%d", i))) {
			falsePositives++
		}
	}

	allowedFalsePositives := 15
	if falsePositives > allowedFalsePositives {
		t.Errorf("False positive rate too high: %d / 1000", falsePositives)
	}
}


func BenchmarkBloomFilter_Add(b *testing.B) {
	bf := NewBloomFilter(1000, 0.01)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(string(rune(i))))
	}
}

func BenchmarkBloomFilter_Exists(b *testing.B) {
	bf := NewBloomFilter(1000, 0.01)
	bf.Add([]byte("test"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bf.Exists([]byte("test"))
	}
}