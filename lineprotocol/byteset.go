package lineprotocol

// newByteSet returns a set representation
// of the bytes in the given string.
func newByteSet(s string) *byteSet {
	var set byteSet
	for i := 0; i < len(s); i++ {
		set.set(s[i])
	}
	return &set
}

func newByteSetRange(i0, i1 uint8) *byteSet {
	var set byteSet
	for i := i0; i <= i1; i++ {
		set.set(i)
	}
	return &set
}

// byteSet is a compact bitset representing a set of byte values.
// It uses 32 bytes ([4]uint64) instead of 256 bytes ([256]bool).
type byteSet [4]uint64

// get reports whether x is in the set.
func (b *byteSet) get(x uint8) bool {
	return b[x/64]&(1<<(x%64)) != 0
}

// set ensures that x is in the set.
func (b *byteSet) set(x uint8) {
	b[x/64] |= 1 << (x % 64)
}

// union returns the union of b and b1.
func (b *byteSet) union(b1 *byteSet) *byteSet {
	return &byteSet{
		b[0] | b1[0],
		b[1] | b1[1],
		b[2] | b1[2],
		b[3] | b1[3],
	}
}

// intersect returns the intersection of b and b1.
func (b *byteSet) intersect(b1 *byteSet) *byteSet {
	return &byteSet{
		b[0] & b1[0],
		b[1] & b1[1],
		b[2] & b1[2],
		b[3] & b1[3],
	}
}

func (b *byteSet) without(b1 *byteSet) *byteSet {
	return b.intersect(b1.invert())
}

// invert returns everything not in b.
func (b *byteSet) invert() *byteSet {
	return &byteSet{
		^b[0],
		^b[1],
		^b[2],
		^b[3],
	}
}
