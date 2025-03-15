package bloomfilter

import (
	"encoding/binary"
	"encoding/gob"
	"hash/fnv"
	"math"
	"os"
	"sync"
)

// BloomFilter represents a probabilistic data structure for fast set membership testing.
// It allows false positives but never false negatives.
type BloomFilter struct {
	bitArray  []uint64
	size      uint
	hashCount uint
	mu        sync.Mutex
}

// NewBloomFilter initializes a Bloom filter with an optimal size and number of hash functions.
// The size (m) and number of hash functions (k) are derived from:
//   - m = -(n * ln(p)) / (ln(2)^2)   [where n = expected elements, p = false positive rate]
//   - k = (m/n) * ln(2)
func NewBloomFilter(numElements uint, falsePositiveRate float64) *BloomFilter {
	m := uint(math.Ceil(-float64(numElements) * math.Log(falsePositiveRate) / math.Pow(math.Ln2, 2)))
	k := uint(math.Ceil((float64(m) / float64(numElements)) * math.Ln2))
	bitArraySize := (m + 63) / 64 // Rounding up to fit uint64 slots

	return &BloomFilter{
		bitArray:  make([]uint64, bitArraySize),
		size:      m,
		hashCount: k,
	}
}

// murmurHash3 generates multiple deterministic hashes based on a seed value.
func murmurHash3(data []byte, seed uint32) uint {
	hash := fnv.New32a()
	hash.Write(data)
	sum := hash.Sum32() ^ seed

	sum ^= sum >> 16
	sum *= 0x85ebca6b
	sum ^= sum >> 13
	sum *= 0xc2b2ae35
	sum ^= sum >> 16

	return uint(sum) % math.MaxUint32
}


// Add inserts an element into the Bloom filter by setting bits at multiple positions.
func (bf *BloomFilter) Add(data []byte) {
	bf.mu.Lock() // Ensure thread safety
	defer bf.mu.Unlock()

	for i := uint(0); i < bf.hashCount; i++ {
		hashVal := murmurHash3(data, uint32(i)) % bf.size
		bf.setBit(hashVal)
	}
}

// Exists checks whether an element is probably in the Bloom filter.
func (bf *BloomFilter) Exists(data []byte) bool {
	for i := uint(0); i < bf.hashCount; i++ {
		hashVal := murmurHash3(data, uint32(i)) % bf.size
		if !bf.getBit(hashVal) {
			return false // If any bit is unset, the element was never added
		}
	}
	return true // All required bits are set → element might exist
}

// setBit sets a specific bit in the bit array.
func (bf *BloomFilter) setBit(index uint) {
	if index/64 >= uint(len(bf.bitArray)) {
		return // Ignore out-of-bounds bits
	}

	bf.bitArray[index/64] |= 1 << (index % 64)
}

// getBit checks if a specific bit is set.
func (bf *BloomFilter) getBit(index uint) bool {
	return (bf.bitArray[index/64] & (1 << (index % 64))) != 0
}

// Save writes the Bloom filter to a file for persistence.
func (bf *BloomFilter) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(bf.bitArray)
}

// Load reads a Bloom filter from a file.
func (bf *BloomFilter) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&bf.bitArray)
}
