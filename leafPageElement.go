package bbolt

import "unsafe"

// leafPageElement represents a node on a leaf page.
type leafPageElement struct {
	flags uint32
	pos   uint32
	ksize uint32
	vsize uint32
}

// key returns a byte slice of the node key.
func (n *leafPageElement) key() []byte {
	i := int(n.pos)
	j := i + int(n.ksize)
	return unsafeByteSlice(unsafe.Pointer(n), 0, i, j)
}

// value returns a byte slice of the node value.
func (n *leafPageElement) value() []byte {
	i := int(n.pos) + int(n.ksize)
	j := i + int(n.vsize)
	return unsafeByteSlice(unsafe.Pointer(n), 0, i, j)
}
