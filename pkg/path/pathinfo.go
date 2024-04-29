package path

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dccn-tg/dr-data-stager/pkg/dr"
	"github.com/cyverse/go-irodsclient/fs"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
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
	// checksum
	checksum string
}

func (p PathInfo) CountFiles(ctx context.Context) int {
	if p.Mode.IsRegular() {
		return 1
	}
	scanner := NewScanner(p)
	return scanner.CountFilesInDir(ctx, p.Path)
}

func (p PathInfo) GetChecksum() string {
	if p.checksum != "" {
		return p.checksum
	}

	if p.Type == TypeFileSystem {
		// calculate checksum for local file
		f, err := os.Open(p.Path)
		if err != nil {
			log.Errorf("%s\n", err)
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Errorf("%s\n", err)
		}

		p.checksum = fmt.Sprintf("%x", h.Sum(nil))
	}

	return p.checksum
}

func (p PathInfo) SameAs(ctx context.Context, o PathInfo) bool {

	if p.Size != o.Size {
		return false
	}

	sum1 := o.GetChecksum()
	sum2 := p.GetChecksum()

	if sum1 == "" || sum2 == "" {
		return false
	}

	return sum1 == sum2
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

			// iRODS file entry contains checksum if it is available
			info.checksum = entry.CheckSum
			info.Size = entry.Size
			return info, nil
		}

		if entry.Type == fs.DirectoryEntry {
			info.Mode = os.ModeDir
			return info, nil
		}

	}

	// local file
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
