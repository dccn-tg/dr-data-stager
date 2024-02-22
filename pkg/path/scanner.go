package path

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/cyverse/go-irodsclient/fs"
	log "github.com/sirupsen/logrus"
)

const (
	blockSize = 4096
	separator = string(filepath.Separator)
)

// See zsyscall_linux_amd64.go/Getdents.
// len(buf)>0.
func getdents(fd int, buf []byte) (n int, err int) {
	var _p0 unsafe.Pointer
	_p0 = unsafe.Pointer(&buf[0])
	r0, _, errno := syscall.Syscall(syscall.SYS_GETDENTS64, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	err = int(errno)
	return
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

// NewScanner determines the path type and returns a corresponding
// implementation of the Scanner interface.
func NewScanner(path PathInfo) Scanner {
	switch path.Type {
	case TypeIrods:
		return IrodsCollectionScanner{base: path}
	default:
		return FileSystemScanner{base: path}
	}
}

// Scanner defines the interface for scanning files iteratively
// from a namespace `path`.
type Scanner interface {
	// ScanMakeDir gets a list of file-like objects iteratively under the given path, and
	// performs mkdir-like operations when the iteration visits a directory-like object.
	//
	// How the iteration is done depends on the implementation. How the mkdir-like operation is
	// performed is also based on the implementation of the `dirmaker`.  Set the `dirmaker` to
	// `nil` to skip the directory making operation.
	//
	// For example, it can be that the Scanner is implemented to loop over a local filesystem using
	// the `filepath.Walk`, while the `dirmaker` is implemented to create a remote iRODS collection.
	ScanMakeDir(ctx context.Context, buffer int, dirmaker *DirMaker) chan string

	// CountFilesInDir counts number of files in a given directory
	CountFilesInDir(ctx context.Context, dir string) int
}

// FileSystemScanner implements the `Scanner` interface for a POSIX-compliant filesystem.
type FileSystemScanner struct {
	dirmaker *DirMaker
	base     PathInfo
}

// ScanMakeDir gets a list of files iteratively under a file system `path`, and performs directory
// creation based on the implementation of the `dirmaker`.
//
// The output is a string channel with the buffer size provided by the `buffer` argument.
// Each element of the channel refers to a file path.  The channel is closed at the end of the scan.
func (s FileSystemScanner) ScanMakeDir(ctx context.Context, buffer int, dirmaker *DirMaker) chan string {

	files := make(chan string, buffer)

	s.dirmaker = dirmaker

	go func() {
		defer close(files)
		if s.base.Mode.IsDir() {
			s.fastWalk(ctx, s.base.Path, false, &files)
		} else {
			files <- s.base.Path
		}
	}()

	return files
}

func (s FileSystemScanner) CountFilesInDir(ctx context.Context, dir string) int {

	c := 0
	files := make(chan string, 10000)
	go func() {
		s.fastWalk(ctx, dir, false, &files)
		defer close(files)
	}()
	for range files {
		c++
	}
	return c
}

// fastWalk uses linux specific way (i.e. syscall.SYS_GETDENT64) to walk through
// files and directories under the given root recursively.  Files are pushed to
// a given channel of type string. The caller is responsible for
// initiating and closing the provided channel.
func (s FileSystemScanner) fastWalk(ctx context.Context, root string, followLink bool, files *chan string) {

	dir, err := os.Open(root)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		return
	}
	defer dir.Close()

	// Opendir.
	// See dir_unix.go/readdirnames.
	buf := make([]byte, blockSize)
	nbuf := len(buf)
	for {
		select {
		case <-ctx.Done():
			log.Debugf("fastWalk aborted: %s", root)
			return
		default:
			var errno int
			nbuf, errno = getdents(int(dir.Fd()), buf)
			if errno != 0 || nbuf <= 0 {
				return
			}

			// See syscall_linux.go/ParseDirent.
			subbuf := buf[0:nbuf]
			for len(subbuf) > 0 {
				dirent := (*syscall.Dirent)(unsafe.Pointer(&subbuf[0]))
				subbuf = subbuf[dirent.Reclen:]
				bytes := (*[10000]byte)(unsafe.Pointer(&dirent.Name[0]))

				// Using Reclen we compute the first multiple of 8 above the length of
				// Dirent.Name. This value can be used to compute the length of long
				// Dirent.Name faster by checking the last 8 bytes only.
				minlen := uintptr(dirent.Reclen) - unsafe.Offsetof(dirent.Name)
				if minlen > 8 {
					minlen -= 8
				} else {
					minlen = 0
				}

				var name = string(bytes[0 : minlen+uintptr(clen(bytes[minlen:]))])
				if name == "." || name == ".." { // Useless names
					continue
				}

				vpath := filepath.Join(root, name)

				switch dirent.Type {
				case syscall.DT_UNKNOWN:
					log.Warnf("unknonw file type: %s", vpath)
				case syscall.DT_REG:
					*files <- vpath
				case syscall.DT_DIR:
					// construct the directory to be created with dirmaker.
					if err := (*s.dirmaker).Mkdir(strings.TrimPrefix(vpath, s.base.Path)); err != nil {
						log.Errorf("Mkdir failure: %s", err.Error())
					}
					s.fastWalk(ctx, vpath, followLink, files)
				case syscall.DT_LNK:

					// TODO: walk through symlinks is not supported due to issue with
					//       eventual infinite walk loop of A -> B -> C -> A cannot be
					//       easily detected.
					// log.Warnf("skip symlink: %s\n", vpath)
					// continue

					if !followLink {
						log.Warnf("skip symlink: %s\n", vpath)
						continue
					}

					// follow the link; but only to its first level referent.
					referent, err := filepath.EvalSymlinks(vpath)
					if err != nil {
						log.Errorf("cannot resolve symlink: %s error: %s\n", vpath, err)
						continue
					}

					// avoid the situation that the symlink refers to its parent, which
					// can cause infinite filesystem walk loop.
					if referent == root {
						log.Warnf("skip path to avoid symlink loop: %s\n", vpath)
						continue
					}

					log.Warnf("symlink only followed to its first non-symlink referent: %s -> %s\n", vpath, referent)
					s.fastWalk(ctx, referent, false, files)

				default:
					log.Warnf("skip unhandled file: %s (type: %s)", vpath, string(dirent.Type))
					continue
				}
			}
		}
	}
}

// IrodsCollectionScanner implements the `Scanner` interface for iRODS.
type IrodsCollectionScanner struct {
	FileSystem *fs.FileSystem
	base       PathInfo
	dirmaker   *DirMaker
}

// ScanMakeDir gets a list of data objects iteratively under a iRODS collection `path`, and performs
// directory creation based on the implementation of `dirmaker`.
//
// The output is a string channel with the buffer size provided by the `buffer` argument.
// Each element of the channel refers to an iRODS data object.  The channel is closed at the end of the scan.
func (s IrodsCollectionScanner) ScanMakeDir(ctx context.Context, buffer int, dirmaker *DirMaker) chan string {

	files := make(chan string, buffer)

	s.dirmaker = dirmaker

	go func() {
		if s.base.Mode.IsDir() {
			s.collWalk(ctx, s.base.Path, &files)
		} else {
			files <- s.base.Path
		}
		defer close(files)
	}()

	return files
}

func (s IrodsCollectionScanner) CountFilesInDir(ctx context.Context, dir string) int {
	c := 0
	files := make(chan string, 10000)
	go func() {
		s.collWalk(ctx, dir, &files)
		defer close(files)
	}()
	for range files {
		c++
	}
	return c
}

// escapeSpecialCharsGenQuery addes "\" in front of the known special characters
// that cannot be passed to GenQuery directly.
func (s IrodsCollectionScanner) escapeSpecialCharsGenQuery(p string) string {

	// note that the special characters need to be handcrafted one by one.
	// so far, the one noticed not being accepted by iRODS GenQuery is "`".
	for _, c := range []string{"`"} {
		p = strings.ReplaceAll(p, c, fmt.Sprintf("\\%s", c))
	}

	return p
}

// collWalk uses the "iquest" command to query file objects and sub-collections within the collection referred
// by `path`.  It pushs file objects to the `files` channel and loop over the sub-collections iteratively.
//
// The caller is responsible for closing the `files` channel.
func (s IrodsCollectionScanner) collWalk(ctx context.Context, path string, files *chan string) {

	entries, err := s.FileSystem.List(path)
	if err != nil {
		log.Error(err)
		return
	}

	if len(entries) == 0 {
		return
	}

	// push collection entries into channel
	chanEntries := make(chan *fs.Entry, len(entries))
	go func() {
		defer close(chanEntries)
		for _, entry := range entries {
			chanEntries <- entry
		}
	}()

	for {
		select {
		case entry, more := <-chanEntries:
			if !more {
				// no more entries to handle
				return
			}
			if entry.Type == fs.FileEntry {
				*files <- entry.Path
			} else {
				if s.dirmaker != nil {
					// perform `MakeDir` with the `dirmaker`
					if err := (*s.dirmaker).Mkdir(strings.TrimPrefix(entry.Path, s.base.Path)); err != nil {
						log.Errorf("Mkdir failure: %s", err.Error())
					}
				}
				s.collWalk(ctx, entry.Path, files)
			}
		case <-ctx.Done():
			// aborted
			log.Debugf("collWalk aborted")
			return
		}
	}
}
