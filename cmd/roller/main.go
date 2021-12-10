package main

import (
	"encoding/binary"
	"fmt"
)

type XRand struct {
	x uint64
}

func (r *XRand) next() uint64 {
	x := r.x
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	r.x = x
	return x * 0x2545F4914F6CDD1D
}

func (r *XRand) bytes() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, r.next())
	for i, s := range b {
		if 0 == s {
			b[i] = 1
		}
	}
	return b
}

func (r *XRand) seed(a, b uint64) {
	a, b = r.splits(a), r.splits(b)
	r.x = r.splits(a ^ b)
}

func (r *XRand) splits(seed uint64) uint64 {
	result := seed + 0x9E3779B97f4A7C15
	result = (result ^ (result >> 30)) * 0xBF58476D1CE4E5B9
	result = (result ^ (result >> 27)) * 0x94D049BB133111EB
	result = (result ^ (result >> 31))
	return result
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	r := XRand{}
	r.x = r.splits(987654)

	fmt.Printf("a: %d \n", r.x)
	fmt.Printf("d1: %d \n", r.next())
	fmt.Printf("d2: %d \n", r.next())
	fmt.Printf("d3: %d \n", r.next())
	fmt.Printf("d4: %d \n", r.next())

	fmt.Printf("x: %x \n", r.x)
}
