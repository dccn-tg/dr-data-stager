package path

import (
	"context"
	"os"
)

// PathInfo defines a data structure of the path information.
type PathInfo struct {
	// Path is the path in question.
	Path string
	// PathType is the namespace type of the path.
	Type PathType
	// Mode is the `os.FileMode` of the path.
	Mode os.FileMode
	// Size
	Size int64
	// Checksum
	Checksum string
}

func (p PathInfo) CountFiles(ctx context.Context) int {
	if p.Mode.IsRegular() {
		return 1
	}
	scanner := NewScanner(p)
	return scanner.CountFilesInDir(ctx, p.Path)
}

func (p PathInfo) SameAs(ctx context.Context, o PathInfo) bool {

	return p.Size == o.Size

	// if p.Size != o.Size {
	// 	return false
	// }

	// if o.Checksum == "" || p.Checksum == "" {
	// 	return false
	// }

	// return o.Checksum == p.Checksum
}
