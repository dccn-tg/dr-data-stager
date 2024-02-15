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
}

func (p PathInfo) CountFiles(ctx context.Context) int {
	if p.Mode.IsRegular() {
		return 1
	}
	scanner := NewScanner(p)
	return scanner.CountFilesInDir(ctx, p.Path)
}
