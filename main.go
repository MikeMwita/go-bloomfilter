package main

import (
	"fmt"
	"github.com/MikeMwita/go-bloomfilter.git/bloomfilter"
)

func main() {
	bf := bloomfilter.NewBloomFilter(100, 0.01)

	bf.Add([]byte("apple"))
	bf.Add([]byte("banana"))
	bf.Add([]byte("grape"))

	fmt.Println("\n=== Bloom Filter Results ===")
	fmt.Println("Exists (apple):", bf.Exists([]byte("apple")))   //  true
	fmt.Println("Exists (banana):", bf.Exists([]byte("banana"))) // true
	fmt.Println("Exists (grape):", bf.Exists([]byte("grape")))   // true
	fmt.Println("Exists (mango):", bf.Exists([]byte("mango")))   // false (not added)

	// Save and Load Bloom Filter
	bf.Save("bloom.gob")

	// Load a new instance
	bf2 := bloomfilter.NewBloomFilter(100, 0.01)
	bf2.Load("bloom.gob")

	fmt.Println("\n=== After Loading from File ===")
	fmt.Println("Exists (apple):", bf2.Exists([]byte("apple")))  // true
}
