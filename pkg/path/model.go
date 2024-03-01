package path

// PathType represents the namespace type a path is referring to.
type PathType int

const (
	// TypeFileSystem is the namespace type for local filesystem.
	TypeFileSystem PathType = iota
	// TypeIrods is the the namespace type for iRODS.
	TypeIrods
)
