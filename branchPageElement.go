package bbolt

import "unsafe"

// branchPageElement represents a node on a branch page.
// branchPageElement 表示一个branch page上的node
type branchPageElement struct {
	pos   uint32
	ksize uint32
	pgid  pgid
}

// key returns a byte slice of the node key.
// key 返回node的key
func (n *branchPageElement) key() []byte {
	return unsafeByteSlice(unsafe.Pointer(n), 0, int(n.pos), int(n.pos)+int(n.ksize))
}
