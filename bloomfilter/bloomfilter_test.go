package bloomfilter

import (
	"fmt"
	"sync"
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


func TestResize(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)

	for i := 0; i < 200; i++ {
		bf.Add([]byte(fmt.Sprintf("element-%d", i)))
	}

	// Verify Bloom filter has resized
	if bf.size <= 100 {
		t.Errorf("Expected Bloom filter to resize, but size remained %d", bf.size)
	}
}


func TestConcurrentInserts(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01)
	var wg sync.WaitGroup

	// Simulate 100 concurrent inserts
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bf.Add([]byte(fmt.Sprintf("item-%d", i)))
		}(i)
	}
	wg.Wait()

	if !bf.Exists([]byte("item-50")) {
		t.Errorf("Expected 'item-50' to exist in Bloom Filter")
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

// Benchmark Insert Speed
func BenchmarkBloomFilter_Insert(b *testing.B) {
	bf := NewBloomFilter(10000, 0.01)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("item-%d", i)))
	}
}

// Benchmark Lookup Speed (Checking Existing Items)
func BenchmarkBloomFilter_Lookup(b *testing.B) {
	bf := NewBloomFilter(10000, 0.01)

	// Preload Bloom filter
	for i := 0; i < 10000; i++ {
		bf.Add([]byte(fmt.Sprintf("item-%d", i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bf.Exists([]byte(fmt.Sprintf("item-%d", i%10000)))
	}
}

