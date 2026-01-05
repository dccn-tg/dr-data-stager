package path

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	ifs "github.com/cyverse/go-irodsclient/fs"
	"github.com/dccn-tg/dr-data-stager/pkg/dr"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

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
			s.goWalk(ctx, s.base.Path, false, &files)
			//s.fastWalk(ctx, s.base.Path, false, &files)
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
		s.goWalk(ctx, s.base.Path, false, &files)
		//s.fastWalk(ctx, dir, false, &files)
		defer close(files)
	}()
	for range files {
		c++
	}
	return c
}

func (s FileSystemScanner) goWalk(ctx context.Context, root string, followLink bool, files *chan string) {

	filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {

		if e != nil {
			log.Warnf("skip file: %s due to %s\n", p, e)
		}

		switch {
		case d.Type().IsDir():
			if s.dirmaker != nil {
				if err := (*s.dirmaker).Mkdir(ctx, strings.TrimPrefix(p, s.base.Path)); err != nil {
					log.Errorf("Mkdir failure: %s\n", err.Error())
				}
			}
		case d.Type().IsRegular():
			*files <- p
		case d.Type() == fs.ModeSymlink:
			log.Warnf("skip symlink: %s\n", p)
		default:
			log.Warnf("skip unsupported file type: %s\n", p)
		}

		return nil
	})
}

// IrodsCollectionScanner implements the `Scanner` interface for iRODS.
type IrodsCollectionScanner struct {
	base     PathInfo
	dirmaker *DirMaker
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
			// ensure the top-level directory at destination exist
			if err := (*s.dirmaker).Mkdir(ctx, ""); err != nil {
				log.Errorf("Mkdir failure: %s", err.Error())
			}
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

// collWalk uses the "iquest" command to query file objects and sub-collections within the collection referred
// by `path`.  It pushs file objects to the `files` channel and loop over the sub-collections iteratively.
//
// The caller is responsible for closing the `files` channel.
func (s IrodsCollectionScanner) collWalk(ctx context.Context, path string, files *chan string) {

	entries, err := ctx.Value(dr.KeyFilesystem).(*ifs.FileSystem).List(path)
	if err != nil {
		log.Errorf("%s\n", err)
		return
	}

	if len(entries) == 0 {
		return
	}

	// push collection entries into channel
	chanEntries := make(chan *ifs.Entry, len(entries))
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
			if entry.Type == ifs.FileEntry {
				*files <- entry.Path
			} else {
				if s.dirmaker != nil {
					// perform `MakeDir` with the `dirmaker`
					if err := (*s.dirmaker).Mkdir(ctx, strings.TrimPrefix(entry.Path, s.base.Path)); err != nil {
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
