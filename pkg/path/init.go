package path

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	"github.com/Donders-Institute/dr-data-stager/pkg/dr"
	"github.com/cyverse/go-irodsclient/fs"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

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

// ScanAndSync walks through the files retrieved from the `bufio.Scanner`,
// sync each file from the `srcColl` collection to the `dstColl` collection.
//
// The sync operation is performed in a concurrent way.  The degree of concurrency is
// defined by number of sync workers, `nworkers`.
//
// Files being successfully synced will be returned as a map with key as the filename
// and value as the checksum of the file.
func ScanAndSync(ctx context.Context, config config.Configuration, src, dst PathInfo, nworkers int) (processed chan SyncError) {

	processed = make(chan SyncError)

	// initiate a source scanner and performs the scan.
	scanner := NewScanner(src)
	dirmaker := NewDirMaker(dst, config)

	files := scanner.ScanMakeDir(ctx, nworkers*8, &dirmaker)

	// create worker group
	var wg sync.WaitGroup
	wg.Add(nworkers)

	// spin off workers
	for i := 1; i <= nworkers; i++ {
		go syncWorker(ctx, &wg, src, dst, files, processed)
	}

	go func() {
		// wait for all workers to finish.
		wg.Wait()
		// close processed channels.
		close(processed)
	}()

	return
}

func syncWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	src, dst PathInfo,
	files chan string,
	processed chan SyncError) {

	// determin the basedir of the source
	srcbase := src.Path
	if src.Mode.IsRegular() {
		srcbase = filepath.Dir(src.Path)
	}

	for {
		select {
		case fsrc, more := <-files:

			// files channel is closed.
			if !more {
				wg.Done()
				return
			}

			// construct the destination path of this particular source `fsrc`
			var fdst string
			if dst.Mode.IsDir() {
				// `dst` is an existing directory/collection, construct the
				// destination file path of this particular file.
				fdst = path.Join(dst.Path, strings.TrimPrefix(fsrc, srcbase))
			} else {
				// destination isn't a directory, then it should be used as the destination file path.
				fdst = dst.Path
			}

			switch {
			case src.Type == TypeIrods && dst.Type == TypeFileSystem:

				psrc, _ := GetPathInfo(ctx, fmt.Sprintf("i:%s", fsrc))
				pdst, _ := GetPathInfo(ctx, fdst)

				if pdst.SameAs(ctx, psrc) {
					log.Debugf("skip transfer: %s == %s\n", fsrc, fdst)
					processed <- SyncError{
						File:  fsrc,
						Error: nil,
					}
					continue
				}

				// get file from irods
				log.Debugf("irods get: %s -> %s\n", fsrc, fdst)
				processed <- SyncError{
					File:  fsrc,
					Error: ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).DownloadFileParallel(fsrc, "", fdst, 0, nil),
				}
			case src.Type == TypeFileSystem && dst.Type == TypeIrods:

				pdst, _ := GetPathInfo(ctx, fmt.Sprintf("i:%s", fdst))
				psrc, _ := GetPathInfo(ctx, fsrc)

				if pdst.SameAs(ctx, psrc) {
					log.Debugf("skip transfer: %s == %s\n", fsrc, fdst)
					processed <- SyncError{
						File:  fsrc,
						Error: nil,
					}
					continue
				}

				// put file to irods
				log.Debugf("irods put: %s -> %s\n", fsrc, fdst)
				processed <- SyncError{
					File:  fsrc,
					Error: ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).UploadFileParallel(fsrc, fdst, "", 0, false, nil),
				}
			default:
				// both source/destination has the same type
				processed <- SyncError{
					File:  fsrc,
					Error: fmt.Errorf("not supported"),
				}
			}
		case <-ctx.Done():
			log.Debugf("sync worker aborted")
			// task has been mark to done.
			wg.Done()
			return
		}
	}
}
