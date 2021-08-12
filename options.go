package bbolt

import (
	"os"
	"time"
)

// DefaultOptions represent the options used if nil options are passed into Open().
// No timeout is used which will cause Bolt to wait indefinitely for a lock.
var DefaultOptions = &Options{
	Timeout:      0,
	NoGrowSync:   false,
	FreelistType: FreelistArrayType,
}

// Options represents the options that can be set when opening a database.
type Options struct {
	// Timeout is the amount of time to wait to obtain a file lock.
	// When set to zero it will wait indefinitely. This option is only
	// available on Darwin and Linux.
	Timeout time.Duration

	// Sets the DB.NoGrowSync flag before memory mapping the file.
	NoGrowSync bool

	// Do not sync freelist to disk. This improves the database write performance
	// under normal operation, but requires a full database re-sync during recovery.
	NoFreelistSync bool

	// FreelistType sets the backend freelist type. There are two options. Array which is simple but endures
	// dramatic performance degradation if database is large and framentation in freelist is common.
	// The alternative one is using hashmap, it is faster in almost all circumstances
	// but it doesn't guarantee that it offers the smallest page id available. In normal case it is safe.
	// The default type is array
	FreelistType FreelistType

	// Open database in read-only mode. Uses flock(..., LOCK_SH |LOCK_NB) to
	// grab a shared lock (UNIX).
	ReadOnly bool

	// Sets the DB.MmapFlags flag before memory mapping the file.
	MmapFlags int

	// InitialMmapSize is the initial mmap size of the database
	// in bytes. Read transactions won't block write transaction
	// if the InitialMmapSize is large enough to hold database mmap
	// size. (See DB.Begin for more information)
	//
	// If <=0, the initial map size is 0.
	// If initialMmapSize is smaller than the previous database size,
	// it takes no effect.
	InitialMmapSize int

	// PageSize overrides the default OS page size.
	// 默认4096
	PageSize int

	// NoSync sets the initial value of DB.NoSync. Normally this can just be
	// set directly on the DB itself when returned from Open(), but this option
	// is useful in APIs which expose Options but not the underlying DB.
	NoSync bool

	// OpenFile is used to open files. It defaults to os.OpenFile. This option
	// is useful for writing hermetic tests.
	OpenFile func(string, int, os.FileMode) (*os.File, error)

	// Mlock locks database file in memory when set to true.
	// It prevents potential page faults, however
	// used memory can't be reclaimed. (UNIX only)
	Mlock bool
}
