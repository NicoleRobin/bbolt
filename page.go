package bbolt

import (
	"fmt"
	"log"
	"os"
	"sort"
	"unsafe"
)

const pageHeaderSize = unsafe.Sizeof(page{})

const minKeysPerPage = 2

const branchPageElementSize = unsafe.Sizeof(branchPageElement{})
const leafPageElementSize = unsafe.Sizeof(leafPageElement{})

const (
	branchPageFlag   = 0x01
	leafPageFlag     = 0x02
	metaPageFlag     = 0x04
	freelistPageFlag = 0x10
)

const (
	bucketLeafFlag = 0x01
)

type pgid uint64

type page struct {
	id    pgid
	flags uint16
	count uint16
	// 表示当前页是否有溢出，溢出多少页
	overflow uint32
}

// typ returns a human readable page type string used for debugging.
func (p *page) typ() string {
	if (p.flags & branchPageFlag) != 0 {
		return "branch"
	} else if (p.flags & leafPageFlag) != 0 {
		return "leaf"
	} else if (p.flags & metaPageFlag) != 0 {
		return "meta"
	} else if (p.flags & freelistPageFlag) != 0 {
		return "freelist"
	}
	return fmt.Sprintf("unknown<%02x>", p.flags)
}

// meta returns a pointer to the metadata section of the page.
// meta() 返回指向page中的meta段的指针
func (p *page) meta() *meta {
	log.Printf("page:meta(), p:%+v", p)
	return (*meta)(unsafeAdd(unsafe.Pointer(p), unsafe.Sizeof(*p)))
}

// leafPageElement retrieves the leaf node by index
// leafPageElement() 根据index查找leafPageElement
func (p *page) leafPageElement(index uint16) *leafPageElement {
	log.Printf("page:leafPageElement(), p:%+v, p:%p, index:%d, leafPageElementSize:%d", p, p, index, leafPageElementSize)
	return (*leafPageElement)(unsafeIndex(unsafe.Pointer(p), unsafe.Sizeof(*p),
		leafPageElementSize, int(index)))
}

// leafPageElements retrieves a list of leaf nodes.
// leafPageElements() 查询page的所有leaf node
func (p *page) leafPageElements() []leafPageElement {
	log.Printf("page.leafPageElements(), p:%+v", p)
	if p.count == 0 {
		return nil
	}
	var elems []leafPageElement
	data := unsafeAdd(unsafe.Pointer(p), unsafe.Sizeof(*p))
	log.Printf("page.leafPageElements(), p:%p, unsafe.Pointer(p):%p, unsafe.Sizeof(*p):%d, data:%p", p, unsafe.Pointer(p), unsafe.Sizeof(*p), data)
	unsafeSlice(unsafe.Pointer(&elems), data, int(p.count))
	return elems
}

// branchPageElement retrieves the branch node by index
// branchPageElement() 根据index查找branchPageElement
func (p *page) branchPageElement(index uint16) *branchPageElement {
	return (*branchPageElement)(unsafeIndex(unsafe.Pointer(p), unsafe.Sizeof(*p),
		branchPageElementSize, int(index)))
}

// branchPageElements retrieves a list of branch nodes.
// branchPageElements() 查询page的所有branch node
func (p *page) branchPageElements() []branchPageElement {
	if p.count == 0 {
		return nil
	}
	var elems []branchPageElement
	data := unsafeAdd(unsafe.Pointer(p), unsafe.Sizeof(*p))
	unsafeSlice(unsafe.Pointer(&elems), data, int(p.count))
	return elems
}

// dump writes n bytes of the page to STDERR as hex output.
func (p *page) hexdump(n int) {
	buf := unsafeByteSlice(unsafe.Pointer(p), 0, 0, n)
	fmt.Fprintf(os.Stderr, "%x\n", buf)
}

// PageInfo represents human readable information about a page.
// PageInfo 表示人类可读的page相关信息
type PageInfo struct {
	ID            int
	Type          string
	Count         int
	OverflowCount int
}

// mergepgids copies the sorted union of a and b into dst.
// If dst is too small, it panics.
func mergepgids(dst, a, b pgids) {
	if len(dst) < len(a)+len(b) {
		panic(fmt.Errorf("mergepgids bad len %d < %d + %d", len(dst), len(a), len(b)))
	}
	// Copy in the opposite slice if one is nil.
	if len(a) == 0 {
		copy(dst, b)
		return
	}
	if len(b) == 0 {
		copy(dst, a)
		return
	}

	// Merged will hold all elements from both lists.
	merged := dst[:0]

	// Assign lead to the slice with a lower starting value, follow to the higher value.
	lead, follow := a, b
	if b[0] < a[0] {
		lead, follow = b, a
	}

	// Continue while there are elements in the lead.
	for len(lead) > 0 {
		// Merge largest prefix of lead that is ahead of follow[0].
		n := sort.Search(len(lead), func(i int) bool { return lead[i] > follow[0] })
		merged = append(merged, lead[:n]...)
		if n >= len(lead) {
			break
		}

		// Swap lead and follow.
		lead, follow = follow, lead[n:]
	}

	// Append what's left in follow.
	_ = append(merged, follow...)
}
