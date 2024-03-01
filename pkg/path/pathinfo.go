package path

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/Donders-Institute/dr-data-stager/pkg/dr"
	"github.com/cyverse/go-irodsclient/fs"
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

// GetPathInfo resolves the PathInfo of the given path.
func GetPathInfo(ctx context.Context, path string) (PathInfo, error) {

	var info PathInfo

	re := regexp.MustCompile(`^(i|irods):`)

	if re.MatchString(path) {

		ipath := strings.TrimSuffix(re.ReplaceAllString(path, ""), "/")

		info.Path = ipath
		info.Type = TypeIrods

		entry, err := ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).Stat(ipath)
		if err != nil {
			return info, err
		}

		if entry.Type == fs.FileEntry {
			info.Mode = 0
			info.Checksum = entry.CheckSum
			info.Size = entry.Size
			return info, nil
		}

		if entry.Type == fs.DirectoryEntry {
			info.Mode = os.ModeDir
			return info, nil
		}

	}

	info.Path = path
	info.Type = TypeFileSystem

	fi, err := os.Stat(path)

	if err != nil {
		return info, err
	}

	info.Mode = fi.Mode()
	info.Size = fi.Size()

	return info, nil
}
