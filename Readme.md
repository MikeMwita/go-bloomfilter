# Bloom Filter in Go

A lightweight, toy  implementation of a **Bloom Filter** in Go. This version is designed for learning and demonstration purposes rather than production use. It leverages **bitwise operations** for efficient storage and **MurmurHash-inspired hashing** for element indexing.

Unlike traditional naive implementations, this version:
- Uses **a compact bit array stored in `[]uint64`** instead of a boolean array.
- Supports **variable hash function counts**, optimizing accuracy and performance.
- Computes an **optimal bit array size** and **hash function count** based on expected input size and desired false positive rate.

## Features
 - **Probabilistic membership testing** (false positives possible, no false negatives).

-  **Efficient bitwise storage** (`[]uint64` instead of simple booleans).

-  **Fast hashing with `murmurHash3`-like function** (FNV-1a with a seed).

 - **Mathematically tuned for accuracy and performance** (false positive rate control).

## How It Works
1. Converts input data into **multiple hash values**.
2. Maps hash outputs to bit positions in the **bit array**.
3. **Setting bits** when adding an element.
4. **Checking bits** when testing membership.
5. If all required bits are set, the element **probably exists**.
6. If any bit is unset, the element **definitely does not exist**.

